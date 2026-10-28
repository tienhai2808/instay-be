package middleware

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/gin-gonic/gin"
)

const (
	CacheKeyPrefix     = "domain_schema:"
	CacheKeyBlacklist  = "domain_blacklist:"
	CacheKeyRateLimit  = "domain_rate:"
	CtxSchema          = "schema"
	CtxDomain          = "domain"
	ValidSchemaTTL     = 24 * time.Hour
	InvalidDomainTTL   = 15 * time.Minute
	BlacklistTTL       = 1 * time.Hour
	RateLimitWindow    = 1 * time.Minute
	RateLimitKeyTTL    = 2 * time.Minute
	MaxDBQueriesPerMin = 100
)

type DomainMiddleware struct {
	logPro     port.LogProvider
	cachePro   port.CacheProvider
	tenantRepo repository.TenantRepository
}

func NewDomainMiddleware(
	logPro port.LogProvider,
	cachePro port.CacheProvider,
	tenantRepo repository.TenantRepository,
) *DomainMiddleware {
	return &DomainMiddleware{
		logPro:     logPro,
		cachePro:   cachePro,
		tenantRepo: tenantRepo,
	}
}

func (m *DomainMiddleware) ResolveTenant(c *gin.Context) {
	if m.shouldSkipTenantResolution(c.Request.Method, c.Request.URL.Path) {
		c.Next()
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	domain := c.Request.Header.Get("Origin")
	if domain == "" {
		domain = c.Request.Host
	}

	domain = m.normalizeDomain(domain)
	if domain == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, response.APIResponse{
			Code:    exceptions.ErrInvalidDomain.Code,
			Message: exceptions.ErrInvalidDomain.Message,
		})
		return
	}

	schema, valid := m.resolveSchema(ctx, domain)
	if !valid {
		c.AbortWithStatusJSON(http.StatusForbidden, response.APIResponse{
			Code:    exceptions.ErrDomainNotAllowed.Code,
			Message: exceptions.ErrDomainNotAllowed.Message,
		})
		return
	}

	c.Set(CtxSchema, schema)
	c.Set(CtxDomain, domain)

	c.Next()
}

func (m *DomainMiddleware) shouldSkipTenantResolution(method, path string) bool {
	if strings.HasSuffix(path, "/ping") || strings.HasSuffix(path, "/metrics") {
		return true
	}

	if method == http.MethodGet {
		if strings.HasSuffix(path, "/health") || strings.HasSuffix(path, "/conversations") {
			return true
		}

		if strings.HasSuffix(path, "/messages") && !strings.Contains(path, "/chats/") {
			return true
		}
	}

	return method == http.MethodPost && strings.HasSuffix(path, "/guests/verify")
}

func (m *DomainMiddleware) AllowOriginWithContext(c *gin.Context, origin string) bool {
	domain := m.normalizeDomain(origin)
	if domain == "" {
		return false
	}

	savedDomain := c.GetString(CtxDomain)
	savedSchema := c.GetString(CtxSchema)
	if savedDomain == "" || savedSchema == "" || savedDomain != domain {
		return false
	}

	return true
}

func (m *DomainMiddleware) resolveSchema(ctx context.Context, domain string) (string, bool) {
	if m.isBlacklisted(ctx, domain) {
		return "", false
	}

	schema, found := m.getFromCache(ctx, domain)
	if found {
		return schema, schema != ""
	}

	rateLimited, err := m.isRateLimited(ctx, domain)
	if err != nil {
		return "", false
	}
	if rateLimited {
		m.addToBlacklist(ctx, domain)
		return "", false
	}

	schema, err = m.queryDatabase(ctx, domain)
	if err != nil {
		return "", false
	}

	if schema != "" {
		m.cacheValidDomain(ctx, domain, schema)
		return schema, true
	}

	m.cacheInvalidDomain(ctx, domain)
	return "", false
}

func (m *DomainMiddleware) normalizeDomain(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))

	if strings.HasPrefix(domain, "http://") || strings.HasPrefix(domain, "https://") {
		u, err := url.Parse(domain)
		if err != nil {
			return ""
		}
		domain = u.Host
	}

	if idx := strings.LastIndex(domain, ":"); idx > 0 {
		if !strings.Contains(domain, "[") {
			domain = domain[:idx]
		}
	}

	if len(domain) == 0 || len(domain) > 253 {
		return ""
	}

	if domain == "localhost" || domain == "127.0.0.1" {
		return domain
	}

	if !strings.Contains(domain, ".") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return ""
	}

	return domain
}

func (m *DomainMiddleware) isBlacklisted(ctx context.Context, domain string) bool {
	key := CacheKeyBlacklist + domain
	exists, err := m.cachePro.Exists(ctx, key)
	if err != nil {
		m.logPro.Error("Get blacklist domain in cache failed", port.String("domain", domain), port.Error(err))
	}

	return exists
}

func (m *DomainMiddleware) addToBlacklist(ctx context.Context, domain string) {
	key := CacheKeyBlacklist + domain
	if err := m.cachePro.SetString(ctx, key, "1", BlacklistTTL); err != nil {
		m.logPro.Error("Save blacklist domain in cache failed", port.String("domain", domain), port.Error(err))
	}
}

func (m *DomainMiddleware) getFromCache(ctx context.Context, domain string) (string, bool) {
	key := CacheKeyPrefix + domain
	schema, err := m.cachePro.GetString(ctx, key)
	if err != nil {
		m.logPro.Error("Get schema for domain in cache failed", port.String("domain", domain), port.Error(err))
		return "", false
	}
	if schema == "" {
		return "", false
	}

	return schema, true
}

func (m *DomainMiddleware) isRateLimited(ctx context.Context, domain string) (bool, error) {
	bucket := time.Now().UTC().Truncate(RateLimitWindow).Format("200601021504")
	key := fmt.Sprintf("%s%s:%s", CacheKeyRateLimit, domain, bucket)

	count, err := m.cachePro.IncrementAndGet(ctx, key, RateLimitKeyTTL)
	if err != nil {
		m.logPro.Error("Increment domain rate limit failed", port.String("domain", domain), port.Error(err))
		return false, err
	}

	return count > MaxDBQueriesPerMin, nil
}

func (m *DomainMiddleware) queryDatabase(ctx context.Context, domain string) (string, error) {
	tenant, err := m.tenantRepo.FindByDomain(ctx, domain)
	if err != nil {
		m.logPro.Error("Find schema for domain failed", port.String("domain", domain), port.Error(err))
		return "", err
	}

	if tenant == nil || tenant.SchemaName == "" {
		return "", nil
	}

	if !isValidSchemaName(tenant.SchemaName) {
		return "", nil
	}

	return tenant.SchemaName, nil
}

func (m *DomainMiddleware) cacheValidDomain(ctx context.Context, domain, schema string) {
	key := CacheKeyPrefix + domain
	if err := m.cachePro.SetString(ctx, key, schema, ValidSchemaTTL); err != nil {
		m.logPro.Error("Cache valid domain failed", port.String("domain", domain), port.String("schema", schema), port.Error(err))
	}
}

func (m *DomainMiddleware) cacheInvalidDomain(ctx context.Context, domain string) {
	key := CacheKeyPrefix + domain
	if err := m.cachePro.SetString(ctx, key, "", InvalidDomainTTL); err != nil {
		m.logPro.Error("Cache invalid domain failed", port.String("domain", domain), port.Error(err))
	}
}

func isValidSchemaName(schema string) bool {
	if len(schema) == 0 || len(schema) > 63 {
		return false
	}

	first := rune(schema[0])
	if !unicode.IsLower(first) && first != '_' {
		return false
	}

	for _, r := range schema {
		if !unicode.IsLower(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}

	return true
}
