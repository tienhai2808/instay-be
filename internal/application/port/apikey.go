package port

type APIKeyProvider interface {
	Generate(schema string) (string, error)

	Parse(key string) (string, error)

	Extract(headerValue string) (string, error)

	HeaderName() string

	HeaderPrefix() string
}
