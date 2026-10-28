package worker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/datatypes"
)

const (
	maxReconnectAttempts = 10
	baseReconnectDelay   = 2 * time.Second
	maxReconnectDelay    = 5 * time.Minute
	idleTimeout          = 25 * time.Minute
	healthCheckInterval  = 5 * time.Minute
	maxEmailPartBytes    = 2 << 20
)

var (
	spaceRegex           = regexp.MustCompile(`[ \t\f\v]+`)
	htmlLineBreakRegex   = regexp.MustCompile(`(?is)<\s*/?\s*(?:br|p|div|tr|li|h[1-6]|td|th)\b[^>]*>`)
	htmlScriptStyleRe    = regexp.MustCompile(`(?is)<\s*script\b[^>]*>.*?<\s*/\s*script\s*>|<\s*style\b[^>]*>.*?<\s*/\s*style\s*>`)
	htmlTagRegex         = regexp.MustCompile(`(?s)<[^>]+>`)
	leadingNumberRegex   = regexp.MustCompile(`\d+`)
	roomBookedPrefixRe   = regexp.MustCompile(`(?i)^\s*\d+\s*(?:x|rooms?|room)?\s*`)
	nonNumericPriceRegex = regexp.MustCompile(`[^0-9,.-]+`)
	decimalCommaSuffixRe = regexp.MustCompile(`,\d{1,2}$`)
	decimalDotSuffixRe   = regexp.MustCompile(`\.\d{1,2}$`)
	dateCandidateRegex   = regexp.MustCompile(`(?i)\b(?:\d{4}-\d{2}-\d{2}(?:\s+\d{2}:\d{2}(?::\d{2})?)?|\d{1,2}[/-]\d{1,2}[/-]\d{2,4}|[A-Za-z]+,\s+[A-Za-z]+\s+\d{1,2},\s+\d{4}(?:\s+(?:from|until)\s+\d{1,2}:\d{2})?|[A-Za-z]+\s+\d{1,2},\s+\d{4}|\d{1,2}\s+[A-Za-z]+\s+\d{4})\b`)
	emailValueRegex      = regexp.MustCompile(`(?i)\b[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}\b`)
	phoneValueRegex      = regexp.MustCompile(`(?i)(?:\d{1,4}\s+)?\+?\d[\d\s()+.-]{7,}\d`)
	bookingNumberRegex   = regexp.MustCompile(`(?i)\b[A-Z0-9][A-Z0-9\-]{5,}\b`)
	dateLikeTokenRegex   = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$|^\d{1,2}[-/]\d{1,2}[-/]\d{2,4}$`)
	guestKeywordRegex    = regexp.MustCompile(`(?i)\b(guest|travell?er|customer|contact)\b`)
	guestLabelPrefixRe   = regexp.MustCompile(`(?i)^\s*(?:guest|travell?er|customer|contact(?:\s*person)?)\s*(?:details?)?\s*[:\-]?\s*`)
	guestNoiseTokenRe    = regexp.MustCompile(`(?i)\b(?:email|e-?mail|phone|mobile|tel|telephone|adults?|adult|children?|child|guests?|pax|rooms?|room|nights?)\b`)
	nameTokenRegex       = regexp.MustCompile(`[\p{L}]+(?:[ '\-][\p{L}]+)*`)
	guestWordRegex       = regexp.MustCompile(`[\p{L}]+`)
	lineLabelPrefixRe    = regexp.MustCompile(`^[\p{L}][\p{L}0-9 &().\-_\/]{0,60}$`)
	knownFieldLabelRe    = regexp.MustCompile(`(?i)^(booking\s*number|guest|check[\s-]?in|check[\s-]?out|booked\s*on|rooms?\s*booked|booking\s*source|total\s*net\s*price|total\s*sell\s*price|number\s*of\s*guests?|number\s*of\s*rooms?|promo(?:tion)?\s*name|meal\s*plan|booking\s*conditions?|booking\s*preferences?)\s*:?\s*$`)
	alphaNumRegex        = regexp.MustCompile(`[A-Za-z0-9]`)
)

var bookingFieldAliases = map[string][]string{
	"booking_number":      {`booking\s*(?:number|no\.?|id)`, `reservation\s*(?:number|no\.?|id)`, `confirmation\s*(?:number|no\.?|id)`},
	"guest_full_name":     {`guest\s*(?:full\s*)?name`, `customer\s*name`, `travell?er\s*name`},
	"guest_email":         {`guest\s*e-?mail`, `customer\s*e-?mail`, `e-?mail`},
	"guest_phone":         {`guest\s*(?:phone|mobile|telephone|tel)`, `contact\s*(?:phone|number)`, `phone`},
	"guest_blob":          {`guest(?:\s*details?)?`, `travell?er(?:\s*details?)?`, `customer(?:\s*details?)?`, `contact\s*(?:person|details?)`},
	"check_in":            {`check[\s-]?in`, `arrival(?:\s*date)?`},
	"check_out":           {`check[\s-]?out`, `departure(?:\s*date)?`},
	"booked_on":           {`booked\s*on`, `booking\s*date`, `reservation\s*date`},
	"room_booked":         {`room\s*booked`, `room\s*type`, `rooms?\s*(?:booked|selected)`},
	"source":              {`source`, `channel`, `ota`, `booking\s*channel`},
	"total_net_price":     {`net\s*(?:price|amount|total)`, `commissionable(?:\s*amount)?`},
	"total_sell_price":    {`sell\s*(?:price|amount|total)`, `total\s*(?:price|amount)`, `amount\s*due`},
	"guest_number":        {`guests?`, `pax`, `occupancy`},
	"room_number":         {`rooms?\s*(?:count|qty|quantity|number|booked)`},
	"promotion_name":      {`promotion(?:\s*name)?`, `promo(?:tion)?`, `discount`},
	"meal_plan":           {`meal\s*plan`, `board\s*type`, `meal`},
	"booking_conditions":  {`booking\s*conditions?`, `cancellation\s*policy`, `conditions?`},
	"booking_preferences": {`booking\s*preferences?`, `special\s*requests?`, `preferences?`},
}

type tenantWorker struct {
	logPro           port.LogProvider
	idPro            port.IDProvider
	tenant           *model.Tenant
	tenantRepo       repository.TenantRepository
	bookingRepo      repository.BookingRepository
	sourceRepo       repository.SourceRepository
	roomTypeRepo     repository.RoomTypeRepository
	dailyRateRepo    repository.DailyRoomTypeRateRepository
	userRepo         repository.UserRepository
	deviceTokenRepo  repository.DeviceTokenRepository
	notificationRepo repository.NotificationRepository
	departmentRepo   repository.DepartmentRepository
	mqPro            port.MessageQueueProvider
	aiPro            port.AIProvider
	settings         model.TenantSettings
}

func newTenantWorker(
	logPro port.LogProvider,
	idPro port.IDProvider,
	tenant *model.Tenant,
	tenantRepo repository.TenantRepository,
	bookingRepo repository.BookingRepository,
	sourceRepo repository.SourceRepository,
	roomTypeRepo repository.RoomTypeRepository,
	dailyRateRepo repository.DailyRoomTypeRateRepository,
	userRepo repository.UserRepository,
	deviceTokenRepo repository.DeviceTokenRepository,
	notificationRepo repository.NotificationRepository,
	departmentRepo repository.DepartmentRepository,
	mqPro port.MessageQueueProvider,
	aiPro port.AIProvider,
	settings model.TenantSettings,
) *tenantWorker {
	return &tenantWorker{
		logPro,
		idPro,
		tenant,
		tenantRepo,
		bookingRepo,
		sourceRepo,
		roomTypeRepo,
		dailyRateRepo,
		userRepo,
		deviceTokenRepo,
		notificationRepo,
		departmentRepo,
		mqPro,
		aiPro,
		settings,
	}
}

func (w *tenantWorker) run(ctx context.Context) {
	w.validateBookingRegexSettings()

	attempt := 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err := w.listenEmails(ctx)
		if err == nil {
			return
		}

		attempt++
		delay := calculateBackoff(attempt)
		w.logPro.Error("Imap listener error, reconnecting",
			port.Error(err),
			port.Int("attempt", attempt),
			port.Duration("retry_in", delay))

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return
		}
	}
}

func (w *tenantWorker) listenEmails(ctx context.Context) error {
	cfg := w.settings.IMAP
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	c, err := client.DialTLS(addr, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = c.Logout()
	}()

	if err := c.Login(cfg.User, cfg.Password); err != nil {
		return err
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return err
	}

	lastSeen := mbox.Messages
	w.logPro.Info("INBOX selected", port.String("schema", w.tenant.SchemaName), port.Uint32("total_messages", lastSeen))

	updates := make(chan client.Update, 10)
	c.Updates = updates

	stop := make(chan struct{})
	done := make(chan error, 1)
	healthTicker := time.NewTicker(healthCheckInterval)
	defer healthTicker.Stop()

	startIdle := func() {
		go func() {
			done <- c.Idle(stop, &client.IdleOptions{
				LogoutTimeout: idleTimeout,
				PollInterval:  0,
			})
		}()
	}
	startIdle()

	for {
		select {
		case <-ctx.Done():
			close(stop)
			<-done
			return nil

		case <-healthTicker.C:

		case update := <-updates:
			mboxUpdate, ok := update.(*client.MailboxUpdate)
			if !ok {
				continue
			}
			if mboxUpdate.Mailbox.Messages <= lastSeen {
				continue
			}

			close(stop)
			if err := <-done; err != nil {
				return err
			}

			if err := w.fetchNewEmails(ctx, c, lastSeen+1, mboxUpdate.Mailbox.Messages); err != nil {
				w.logPro.Error("Fetch new emails failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			} else {
				lastSeen = mboxUpdate.Mailbox.Messages
			}

			stop = make(chan struct{})
			done = make(chan error, 1)
			startIdle()

		case err := <-done:
			if err != nil {
				return err
			}
		}
	}
}

func (w *tenantWorker) fetchNewEmails(ctx context.Context, c *client.Client, from, to uint32) error {
	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, imap.FetchFlags, imap.FetchUid, section.FetchItem()}

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	targetBookingSubjects := w.settings.EmailBooking.Subject
	targetCancelSubjects := w.settings.EmailCancel.Subject
	targetCancelBodies := w.settings.EmailCancel.Bodies
	for msg := range messages {
		if msg == nil {
			continue
		}
		if msg.Envelope == nil {
			continue
		}

		subject := msg.Envelope.Subject
		isBookingSubject := matchPatternList(targetBookingSubjects, subject)
		isCancelSubject := matchPatternList(targetCancelSubjects, subject)

		needBody := isBookingSubject || isCancelSubject || len(targetCancelBodies) > 0
		htmlBody := ""
		plainBody := ""
		bodyText := ""
		if needBody {
			var err error
			htmlBody, plainBody, err = readEmailBodies(msg, section)
			if err != nil {
				w.logPro.Error("Read email body failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
				continue
			}

			bodyText = strings.Join(compactNonEmpty([]string{normalizeLines(plainBody), htmlToPlainText(htmlBody)}), "\n")
		}

		if isCancelSubject || len(targetCancelBodies) > 0 {
			if handled := w.processCancelEmail(ctx, subject, bodyText, isCancelSubject, targetCancelBodies); handled {
				continue
			}
		}
		if !isBookingSubject {
			continue
		}
		w.processBookingEmail(ctx, subject, htmlBody, plainBody)
	}

	return <-done
}

func (w *tenantWorker) processBookingEmail(ctx context.Context, subject, htmlBody, plainBody string) {
	if htmlBody == "" && plainBody == "" {
		return
	}

	booking, err := w.parseBooking(ctx, subject, htmlBody, plainBody)
	if err != nil {
		w.logPro.Error("Parse booking failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return
	}

	id, err := w.idPro.NextID()
	if err != nil {
		w.logPro.Error("Generate booking id failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return
	}
	booking.ID = int64(id)
	booking.Status = model.BookingStatusConfirmed

	if err := w.bookingRepo.Create(ctx, w.tenant.SchemaName, booking); err != nil {
		w.logPro.Error("Create booking failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return
	}
	w.logPro.Info("Booking created", port.String("schema", w.tenant.SchemaName), port.Int64("booking_id", booking.ID), port.String("booking_number", booking.BookingNumber))
	w.predictBooking(ctx, booking.ID)

	notificationContentVI := fmt.Sprintf("Bạn có 1 booking mới từ %s", booking.Source.Name)
	notificationContentEN := fmt.Sprintf("You have a new booking from %s", booking.Source.Name)

	w.notifyBookingManagement(ctx, booking.ID, notificationContentVI, notificationContentEN, constants.CodePushNotificationBookingNew)
}

func (w *tenantWorker) predictBooking(ctx context.Context, bookingID int64) {
	go func() {
		predictCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		booking, err := w.bookingRepo.FindByIDWithDetails(predictCtx, w.tenant.SchemaName, bookingID)
		if err != nil {
			w.logPro.Error(
				"find booking by id for prediction failed",
				port.String("schema", w.tenant.SchemaName),
				port.Int64("booking_id", bookingID),
				port.Error(err),
			)
			return
		}
		if booking == nil {
			return
		}

		predictRes, err := w.aiPro.PredictBooking(predictCtx, booking)
		if err != nil {
			w.logPro.Error(
				"predict booking from gRPC failed",
				port.String("schema", w.tenant.SchemaName),
				port.Int64("booking_id", bookingID),
				port.Error(err),
			)
			return
		}

		predict := model.BookingPredict{
			CancellationProbability: predictRes.CancellationProbability,
			WillCancel:              predictRes.WillCancel,
			RiskLevel:               predictRes.RiskLevel,
		}

		predictJSON, err := json.Marshal(predict)
		if err != nil {
			w.logPro.Error("Marshal booking prediction failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			return
		}

		if err := w.bookingRepo.Update(predictCtx, w.tenant.SchemaName, bookingID, map[string]any{
			"predict": datatypes.JSON(predictJSON),
		}); err != nil {
			w.logPro.Error(
				"update booking prediction failed",
				port.String("schema", w.tenant.SchemaName),
				port.Int64("booking_id", bookingID),
				port.Error(err),
			)
			return
		}
	}()
}

func (w *tenantWorker) processCancelEmail(
	ctx context.Context,
	subject string,
	bodyText string,
	isCancelSubject bool,
	cancelBodyPatterns model.PatternList,
) bool {
	subject = cleanExtractedValue(subject)

	isCancelBody := matchPatternList(cancelBodyPatterns, bodyText)
	if !isCancelSubject && !isCancelBody {
		return false
	}
	isCancelSignal := isCancelSubject || isCancelBody

	candidates := w.parseCancelBookingCandidates(subject, bodyText)
	booking, matchedBookingNumber, err := w.findBestCancelBooking(ctx, candidates)
	if err != nil {
		w.logPro.Error("Find booking for cancel email failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return isCancelSignal
	}
	if booking == nil {
		w.logPro.Warn(
			"booking not found for cancel email",
			port.String("schema", w.tenant.SchemaName),
			port.String("subject", subject),
		)
		return isCancelSubject
	}

	if booking.Status != model.BookingStatusConfirmed {
		w.logPro.Info(
			"skip cancel email because booking is not confirmed",
			port.String("schema", w.tenant.SchemaName),
			port.Int64("booking_id", booking.ID),
			port.String("booking_number", booking.BookingNumber),
			port.String("status", string(booking.Status)),
		)
		return true
	}

	updateData := map[string]any{
		"status": model.BookingStatusCancelled,
	}
	if err = w.bookingRepo.Update(ctx, w.tenant.SchemaName, booking.ID, updateData); err != nil {
		w.logPro.Error(
			"update booking status failed",
			port.String("schema", w.tenant.SchemaName),
			port.Int64("booking_id", booking.ID),
			port.String("booking_number", booking.BookingNumber),
			port.Error(err),
		)
		return true
	}

	w.logPro.Info(
		"booking cancelled by email",
		port.String("schema", w.tenant.SchemaName),
		port.Int64("booking_id", booking.ID),
		port.String("booking_number", booking.BookingNumber),
		port.String("matched_booking_number", matchedBookingNumber),
	)

	notificationContentVI := fmt.Sprintf("Booking %s đã bị hủy", booking.BookingNumber)
	notificationContentEN := fmt.Sprintf("Booking %s has been cancelled", booking.BookingNumber)
	w.notifyBookingManagement(ctx, booking.ID, notificationContentVI, notificationContentEN, constants.CodePushNotificationBookingCancelled)

	return true
}

func (w *tenantWorker) parseCancelBookingCandidates(subject, bodyText string) []string {
	candidates := make([]string, 0, 16)
	seen := make(map[string]struct{}, 16)

	appendCandidate := func(values []string) {
		for _, value := range values {
			value = normalizeBookingNumber(value)
			if value == "" {
				continue
			}
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			candidates = append(candidates, value)
		}
	}

	for _, pattern := range w.settings.EmailCancel.Subject {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			continue
		}

		allMatches := re.FindAllStringSubmatch(subject, -1)
		if len(allMatches) == 0 {
			continue
		}

		for _, matches := range allMatches {
			for i := 1; i < len(matches); i++ {
				appendCandidate(extractBookingNumberTokens(matches[i]))
			}
			appendCandidate(extractBookingNumberTokens(matches[0]))
		}
	}

	for _, pattern := range w.settings.EmailCancel.Bodies {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			continue
		}

		allMatches := re.FindAllStringSubmatch(bodyText, -1)
		if len(allMatches) == 0 {
			continue
		}

		for _, matches := range allMatches {
			for i := 1; i < len(matches); i++ {
				appendCandidate(extractBookingNumberTokens(matches[i]))
			}
			appendCandidate(extractBookingNumberTokens(matches[0]))
		}
	}

	appendCandidate(extractBookingNumberTokens(subject))
	appendCandidate(extractBookingNumberTokens(bodyText))

	return candidates
}

func (w *tenantWorker) findBestCancelBooking(ctx context.Context, candidates []string) (*model.Booking, string, error) {
	var fallbackBooking *model.Booking
	fallbackNumber := ""

	for _, bookingNumber := range candidates {
		booking, err := w.bookingRepo.FindByBookingNumber(ctx, w.tenant.SchemaName, bookingNumber)
		if err != nil {
			return nil, "", err
		}
		if booking == nil {
			continue
		}

		if booking.Status == model.BookingStatusConfirmed {
			return booking, bookingNumber, nil
		}

		if fallbackBooking == nil {
			fallbackBooking = booking
			fallbackNumber = bookingNumber
		}
	}

	return fallbackBooking, fallbackNumber, nil
}

func readEmailBodies(msg *imap.Message, section *imap.BodySectionName) (string, string, error) {
	r := msg.GetBody(section)
	if r == nil {
		return "", "", nil
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", "", err
	}

	var htmlBody strings.Builder
	var plainBody strings.Builder
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return htmlBody.String(), plainBody.String(), err
		}
		if h, ok := p.Header.(*mail.InlineHeader); ok {
			ct, _, _ := h.ContentType()
			body, err := io.ReadAll(io.LimitReader(p.Body, maxEmailPartBytes))
			if err != nil {
				continue
			}

			switch strings.ToLower(ct) {
			case "text/html":
				if htmlBody.Len() > 0 {
					htmlBody.WriteString("\n")
				}
				htmlBody.Write(body)
			case "text/plain":
				if plainBody.Len() > 0 {
					plainBody.WriteString("\n")
				}
				plainBody.Write(body)
			}
		}
	}

	return htmlBody.String(), plainBody.String(), nil
}

func (w *tenantWorker) notifyBookingManagement(ctx context.Context, bookingID int64, contentVI, contentEN string, pushCode int) {
	departments, err := w.departmentRepo.FindAllByPermissionNameWithActiveUsers(ctx, w.tenant.SchemaName, model.PermBookingManagement)
	if err != nil {
		w.logPro.Error("Find all department by permission name failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return
	}

	var userIDs []int64
	for _, dept := range departments {
		nid, err := w.idPro.NextID()
		if err != nil {
			w.logPro.Error("Generate notification id failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			return
		}

		notification := model.NewNotification(
			nid,
			&dept.ID,
			nil,
			model.NotificationReceiverTypeStaff,
			model.NotificationContentTypeBooking,
			strconv.FormatInt(bookingID, 10),
			contentVI,
			contentEN,
			false,
		)
		if err = w.notificationRepo.Create(ctx, w.tenant.SchemaName, notification); err != nil {
			w.logPro.Error("Create notification failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			return
		}

		var userDeptIDs []int64
		for _, user := range dept.Users {
			userDeptIDs = append(userDeptIDs, user.ID)
		}

		userIDs = append(userIDs, userDeptIDs...)
	}

	deviceTokens, err := w.deviceTokenRepo.FindAllByUserIDs(ctx, w.tenant.SchemaName, userIDs)
	if err != nil {
		w.logPro.Error("Find all device token by user ids failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return
	}
	if len(deviceTokens) == 0 {
		return
	}

	tokens := make([]string, 0, len(deviceTokens))
	for _, deviceToken := range deviceTokens {
		tokens = append(tokens, deviceToken.Token)
	}

	pushNotificationMsg := data.PushNotificationMessage{
		Schema: w.tenant.SchemaName,
		Tokens: tokens,
		Data: map[string]string{
			"code":         strconv.Itoa(pushCode),
			"content_key":  strconv.FormatInt(bookingID, 10),
			"content_type": string(model.NotificationContentTypeBooking),
		},
	}

	go func(msg data.PushNotificationMessage, pushTokens []string) {
		body, err := json.Marshal(msg)
		if err != nil {
			w.logPro.Error("Json marshal failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		}

		routingKey := utils.GetRoutingKey(w.tenant.SchemaName, port.RoutingKeyPushNotification)
		if err = w.mqPro.PublishMessage(port.ExchangePushNotification, routingKey, body); err != nil {
			w.logPro.Error("Publish push notification message failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		}

		if err = w.deviceTokenRepo.Updates(ctx, w.tenant.SchemaName, pushTokens, map[string]any{"last_used_at": time.Now()}); err != nil {
			w.logPro.Error("Update device token by tokens failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		}
	}(pushNotificationMsg, tokens)
}

func (w *tenantWorker) parseBooking(ctx context.Context, subject, htmlContent, plainContent string) (*model.Booking, error) {
	sel := w.settings.EmailBooking
	sources, lineText := buildSearchSources(subject, plainContent, htmlContent)
	labelValues := buildLabelValueMap(lineText)

	extract := func(field string, patterns model.PatternList) string {
		return extractFieldValue(field, patterns, sources, lineText, labelValues)
	}
	optional := func(v string) *string {
		if v == "" {
			return nil
		}
		return &v
	}

	booking := &model.Booking{}

	booking.BookingNumber = extract("booking_number", sel.BookingNumber)
	booking.BookingNumber = normalizeBookingNumber(booking.BookingNumber)
	if booking.BookingNumber == "" {
		return nil, fmt.Errorf("could not parse booking number")
	}

	guestBlob := extract("guest_blob", sel.Guest)
	guestName, guestEmail, guestPhone := enrichGuestDetails("", "", "", guestBlob, lineText)
	booking.GuestFullName = guestName
	if booking.GuestFullName == "" {
		booking.GuestFullName = "Unknown Guest"
	}
	booking.GuestEmail = optional(guestEmail)
	booking.GuestPhone = optional(guestPhone)
	booking.GuestNumber = extract("guest_number", sel.GuestNumber)
	if booking.GuestNumber == "" {
		booking.GuestNumber = "1"
	}
	booking.PromotionName = optional(extract("promotion_name", sel.PromotionName))
	booking.MealPlan = optional(extract("meal_plan", sel.MealPlan))
	booking.BookingConditions = optional(extract("booking_conditions", sel.BookingConditions))
	booking.BookingPreferences = optional(extract("booking_preferences", sel.BookingPreferences))

	checkInRaw := extract("check_in", sel.CheckIn)
	checkOutRaw := extract("check_out", sel.CheckOut)
	bookedOnRaw := extract("booked_on", sel.BookedOn)

	booking.CheckIn = parseDateString(checkInRaw)
	booking.CheckOut = parseDateString(checkOutRaw)
	booking.BookedOn = parseDateString(bookedOnRaw)

	candidates := extractDateCandidates(lineText)
	if booking.CheckIn.IsZero() && len(candidates) > 0 {
		booking.CheckIn = candidates[0]
	}
	if booking.CheckOut.IsZero() {
		switch {
		case len(candidates) > 1:
			booking.CheckOut = candidates[1]
		case !booking.CheckIn.IsZero():
			booking.CheckOut = booking.CheckIn.Add(24 * time.Hour)
		}
	}
	if booking.BookedOn.IsZero() {
		switch {
		case len(candidates) > 2:
			booking.BookedOn = candidates[2]
		default:
			booking.BookedOn = time.Now().UTC()
		}
	}
	if booking.CheckIn.IsZero() || booking.CheckOut.IsZero() {
		return nil, fmt.Errorf("could not parse check-in/check-out")
	}

	roomBookedRaw := extract("room_booked", sel.RoomsBooked)
	roomNumberRaw := extract("room_number", sel.RoomNumber)
	booking.RoomNumber = parseUint32(roomNumberRaw)
	parsedRoomNumber, roomTypeName := parseRoomBooked(roomBookedRaw)
	if booking.RoomNumber == 0 {
		booking.RoomNumber = parsedRoomNumber
	}
	if booking.RoomNumber == 0 {
		booking.RoomNumber = 1
	}

	roomType, err := w.ensureRoomType(ctx, roomTypeName)
	if err != nil {
		return nil, err
	}
	booking.RoomType = roomType
	booking.RoomTypeID = roomType.ID

	rawNetPrice := extract("total_net_price", sel.TotalNetPrice)
	rawSellPrice := extract("total_sell_price", sel.TotalSellPrice)
	hasNetPrice := false
	hasSellPrice := false

	switch {
	case rawSellPrice != "":
		booking.TotalSellPrice = parsePrice(rawSellPrice)
		hasSellPrice = true
	case rawNetPrice != "":
		booking.TotalSellPrice = parsePrice(rawNetPrice)
		hasSellPrice = true
	}
	if rawSellPrice != "" && rawNetPrice != "" {
		booking.TotalNetPrice = parsePrice(rawNetPrice)
		hasNetPrice = true
	}
	if !hasNetPrice && booking.RoomType != nil {
		totalNetPrice, err := w.calculateBookingTotalNetPrice(ctx, booking)
		if err != nil {
			return nil, err
		}
		booking.TotalNetPrice = totalNetPrice
		hasNetPrice = true
	}
	if !hasSellPrice && hasNetPrice {
		booking.TotalSellPrice = booking.TotalNetPrice
	}

	source, err := w.ensureSource(ctx, extract("source", sel.Source))
	if err != nil {
		return nil, err
	}
	booking.Source = source
	booking.SourceID = source.ID

	sanitizeBookingForDB(booking)

	return booking, nil
}

func buildSearchSources(subject, plainContent, htmlContent string) ([]string, string) {
	subjectText := normalizeLines(subject)
	plainText := normalizeLines(plainContent)
	htmlText := htmlToPlainText(htmlContent)

	sources := make([]string, 0, 3)
	for _, source := range []string{subjectText, plainText, htmlText} {
		if source == "" {
			continue
		}
		sources = append(sources, source)
	}

	lineText := strings.Join(compactNonEmpty([]string{subjectText, plainText, htmlText}), "\n")
	return sources, lineText
}

func compactNonEmpty(in []string) []string {
	res := make([]string, 0, len(in))
	for _, item := range in {
		if strings.TrimSpace(item) == "" {
			continue
		}
		res = append(res, item)
	}
	return res
}

func normalizeLines(raw string) string {
	raw = html.UnescapeString(raw)
	raw = strings.ReplaceAll(raw, "\u00a0", " ")
	raw = strings.ReplaceAll(raw, "\r\n", "\n")
	raw = strings.ReplaceAll(raw, "\r", "\n")
	lines := strings.Split(raw, "\n")
	normalized := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(spaceRegex.ReplaceAllString(line, " "))
		if line == "" {
			continue
		}
		normalized = append(normalized, line)
	}
	return strings.Join(normalized, "\n")
}

func htmlToPlainText(raw string) string {
	if strings.TrimSpace(raw) == "" {
		return ""
	}
	res := htmlScriptStyleRe.ReplaceAllString(raw, " ")
	res = htmlLineBreakRegex.ReplaceAllString(res, "\n")
	res = htmlTagRegex.ReplaceAllString(res, " ")
	return normalizeLines(res)
}

func extractFieldValue(field string, patterns model.PatternList, sources []string, lineText string, labelValues map[string]string) string {
	if v := extractValueFromLabelMap(patterns, labelValues); v != "" {
		return v
	}
	if v := extractValueByPatterns(patterns, sources); v != "" {
		return v
	}
	if v := extractValueFromLabelMap(bookingFieldAliases[field], labelValues); v != "" {
		return v
	}
	return extractValueByAliases(field, lineText)
}

func extractValueFromLabelMap(patterns []string, labelValues map[string]string) string {
	if len(patterns) == 0 || len(labelValues) == 0 {
		return ""
	}
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			continue
		}
		for label, value := range labelValues {
			if !re.MatchString(label) {
				continue
			}
			value = cleanExtractedValue(value)
			if value == "" || looksLikeFieldLabel(value) {
				continue
			}
			return value
		}
	}
	return ""
}

func extractValueByPatterns(patterns model.PatternList, sources []string) string {
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			continue
		}
		for _, source := range sources {
			if re.NumSubexp() == 0 {
				if v := extractValueAfterPattern(source, pattern); v != "" {
					return v
				}
			}

			matches := re.FindStringSubmatch(source)
			if len(matches) > 0 {
				for i := 1; i < len(matches); i++ {
					if v := cleanExtractedValue(matches[i]); v != "" && !looksLikeFieldLabel(v) {
						return v
					}
				}
				if v := cleanExtractedValue(matches[0]); v != "" && !looksLikeFieldLabel(v) {
					return v
				}
			}
		}
	}
	return ""
}

func extractValueAfterPattern(source, pattern string) string {
	blockPattern := `(?is)(?:` + pattern + `)[ \t]*[:\-]?[ \t]*(.+?)(?:\n\s*[A-Za-z][A-Za-z0-9 &().\-/_]{1,60}\s*:\s*|$)`
	if re, err := regexp.Compile(blockPattern); err == nil {
		if m := re.FindStringSubmatch(source); len(m) > 1 {
			v := cleanExtractedValue(m[1])
			if v != "" && !looksLikeFieldLabel(v) {
				return v
			}
		}
	}

	linePattern := `(?im)(?:` + pattern + `)\s*[:\-]?\s*(.+)$`
	if re, err := regexp.Compile(linePattern); err == nil {
		if m := re.FindStringSubmatch(source); len(m) > 1 {
			v := cleanExtractedValue(m[1])
			if v != "" && !looksLikeFieldLabel(v) {
				return v
			}
		}
	}

	return ""
}

func extractValueByAliases(field, lineText string) string {
	aliases := bookingFieldAliases[field]
	for _, alias := range aliases {
		linePattern := `(?im)^\s*(?:[-*]\s*)?(?:` + alias + `)\s*(?:[:\-]\s*|\s{2,})(.+?)\s*$`
		re, err := regexp.Compile(linePattern)
		if err == nil {
			if m := re.FindStringSubmatch(lineText); len(m) > 1 {
				if v := cleanExtractedValue(m[1]); v != "" {
					return v
				}
			}
		}

		inlinePattern := `(?im)(?:` + alias + `)\s*[:\-]\s*(.+?)\s*(?:\n|$)`
		re, err = regexp.Compile(inlinePattern)
		if err == nil {
			if m := re.FindStringSubmatch(lineText); len(m) > 1 {
				if v := cleanExtractedValue(m[1]); v != "" {
					return v
				}
			}
		}
	}
	return ""
}

func cleanExtractedValue(raw string) string {
	raw = html.UnescapeString(raw)
	raw = strings.ReplaceAll(raw, "\u00a0", " ")
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	return strings.Join(strings.Fields(raw), " ")
}

func clipString(raw string, max int) string {
	raw = sanitizeLooseText(raw)
	if raw == "" || max <= 0 {
		return ""
	}
	runes := []rune(raw)
	if len(runes) <= max {
		return raw
	}
	return strings.TrimSpace(string(runes[:max]))
}

func sanitizeLooseText(raw string) string {
	raw = html.UnescapeString(raw)
	raw = strings.ReplaceAll(raw, "\u00a0", " ")
	if strings.Contains(raw, "<") && strings.Contains(raw, ">") {
		raw = htmlLineBreakRegex.ReplaceAllString(raw, "\n")
		raw = htmlTagRegex.ReplaceAllString(raw, " ")
	}
	raw = cleanExtractedValue(raw)
	if strings.HasPrefix(strings.ToLower(raw), "<!doctype") {
		return ""
	}
	return raw
}

func clipStringPtr(raw *string, max int) *string {
	if raw == nil {
		return nil
	}
	value := clipString(*raw, max)
	if value == "" || looksLikeFieldLabel(value) || value == ":" {
		return nil
	}
	return &value
}

func normalizeGuestName(raw string) string {
	raw = sanitizeLooseText(raw)
	if raw == "" || looksLikeFieldLabel(raw) {
		return ""
	}
	words := guestWordRegex.FindAllString(raw, -1)
	if len(words) == 0 {
		return ""
	}
	if len(words) > 5 {
		words = words[:5]
	}
	return clipString(strings.Join(words, " "), 150)
}

func sanitizeBookingForDB(booking *model.Booking) {
	if booking == nil {
		return
	}

	booking.BookingNumber = clipString(normalizeBookingNumber(booking.BookingNumber), 50)
	booking.GuestFullName = normalizeGuestName(booking.GuestFullName)
	if booking.GuestFullName == "" {
		booking.GuestFullName = "Unknown Guest"
	}

	if booking.GuestEmail != nil {
		email := strings.ToLower(clipString(*booking.GuestEmail, 150))
		if email == "" || !emailValueRegex.MatchString(email) {
			booking.GuestEmail = nil
		} else {
			booking.GuestEmail = &email
		}
	}

	if booking.GuestPhone != nil {
		phone := clipString(normalizePhone(*booking.GuestPhone), 20)
		if phone == "" {
			booking.GuestPhone = nil
		} else {
			booking.GuestPhone = &phone
		}
	}

	booking.GuestNumber = clipString(booking.GuestNumber, 50)
	if booking.GuestNumber == "" {
		booking.GuestNumber = "1"
	}

	booking.PromotionName = clipStringPtr(booking.PromotionName, 150)
	booking.MealPlan = clipStringPtr(booking.MealPlan, 150)
	booking.BookingPreferences = clipStringPtr(booking.BookingPreferences, 255)
	booking.BookingConditions = clipStringPtr(booking.BookingConditions, 255)

	if booking.Source != nil {
		booking.Source.Name = clipString(booking.Source.Name, 150)
	}
	if booking.RoomType != nil {
		booking.RoomType.Name = clipString(booking.RoomType.Name, 150)
	}
}

func looksLikeFieldLabel(raw string) bool {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return false
	}
	return knownFieldLabelRe.MatchString(raw)
}

func buildLabelValueMap(lineText string) map[string]string {
	lines := strings.Split(normalizeLines(lineText), "\n")
	res := make(map[string]string, len(lines))

	for i := 0; i < len(lines); i++ {
		line := cleanExtractedValue(lines[i])
		if line == "" {
			continue
		}

		label, value, ok := splitLabelValue(line)
		if !ok {
			continue
		}
		if value == "" && i+1 < len(lines) {
			next := cleanExtractedValue(lines[i+1])
			if next != "" && !looksLikeLineLabel(next) {
				value = next
			}
		}

		if value == "" || looksLikeFieldLabel(value) {
			continue
		}
		res[label] = value
	}

	return res
}

func splitLabelValue(line string) (string, string, bool) {
	idx := strings.Index(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	left := cleanExtractedValue(line[:idx])
	if left == "" || !lineLabelPrefixRe.MatchString(left) {
		return "", "", false
	}
	return normalizeLabelKey(left), cleanExtractedValue(line[idx+1:]), true
}

func normalizeLabelKey(raw string) string {
	raw = strings.ToLower(cleanExtractedValue(raw))
	raw = strings.Trim(raw, ":.-_ ")
	raw = strings.Join(strings.Fields(raw), " ")
	return raw
}

func looksLikeLineLabel(line string) bool {
	line = cleanExtractedValue(line)
	if line == "" {
		return false
	}
	if _, _, ok := splitLabelValue(line); ok {
		return true
	}
	return looksLikeFieldLabel(line)
}

func normalizeBookingNumber(raw string) string {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return ""
	}
	if token := extractBookingNumberToken(raw); token != "" {
		return token
	}
	return raw
}

func extractBookingNumberToken(raw string) string {
	tokens := extractBookingNumberTokens(raw)
	if len(tokens) > 0 {
		return tokens[0]
	}
	return ""
}

func extractBookingNumberTokens(raw string) []string {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return nil
	}

	matches := bookingNumberRegex.FindAllString(strings.ToUpper(raw), -1)
	if len(matches) == 0 {
		return nil
	}

	res := make([]string, 0, len(matches))
	seen := make(map[string]struct{}, len(matches))
	for _, match := range matches {
		if dateLikeTokenRegex.MatchString(match) {
			continue
		}

		hasDigit := false
		for _, r := range match {
			if r >= '0' && r <= '9' {
				hasDigit = true
				break
			}
		}
		if !hasDigit {
			continue
		}

		if _, ok := seen[match]; ok {
			continue
		}
		seen[match] = struct{}{}
		res = append(res, match)
	}

	return res
}

func enrichGuestDetails(name, email, phone, guestBlob, lineText string) (string, string, string) {
	name = cleanExtractedValue(name)
	email = cleanExtractedValue(email)
	phone = cleanExtractedValue(phone)
	guestBlob = cleanExtractedValue(guestBlob)
	lineText = normalizeLines(lineText)

	contexts := buildGuestContexts(guestBlob, lineText)
	if email == "" {
		email = findFirstEmail(contexts)
		if email == "" {
			email = findFirstEmail([]string{lineText})
		}
	}
	if phone == "" {
		phone = findFirstPhone(contexts)
		if phone == "" {
			phone = findFirstPhone([]string{lineText})
		}
	}

	if name == "" || strings.EqualFold(name, "unknown guest") {
		for _, context := range contexts {
			name = parseGuestNameFromContext(context, email, phone)
			if name != "" {
				break
			}
		}
		if name == "" {
			name = parseGuestNameFromContext(lineText, email, phone)
		}
	}

	return name, email, phone
}

func buildGuestContexts(guestBlob, lineText string) []string {
	contexts := make([]string, 0, 8)
	seen := make(map[string]struct{}, 8)
	push := func(value string) {
		value = cleanExtractedValue(value)
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		contexts = append(contexts, value)
	}

	push(guestBlob)
	lines := strings.Split(lineText, "\n")
	for i, line := range lines {
		line = cleanExtractedValue(line)
		if line == "" || !guestKeywordRegex.MatchString(line) {
			continue
		}
		push(line)
		if i+1 < len(lines) && !guestKeywordRegex.MatchString(lines[i+1]) {
			push(line + " " + lines[i+1])
		}
		if i+2 < len(lines) && !guestKeywordRegex.MatchString(lines[i+2]) {
			push(line + " " + lines[i+1] + " " + lines[i+2])
			push(lines[i+1] + " " + lines[i+2])
		}
	}
	return contexts
}

func findFirstEmail(texts []string) string {
	for _, text := range texts {
		email := emailValueRegex.FindString(text)
		if email != "" {
			return cleanExtractedValue(email)
		}
	}
	return ""
}

func findFirstPhone(texts []string) string {
	for _, text := range texts {
		indices := phoneValueRegex.FindAllStringIndex(text, -1)
		for _, idx := range indices {
			if len(idx) != 2 {
				continue
			}
			if !isBoundedByNonAlnum(text, idx[0], idx[1]) {
				continue
			}
			match := text[idx[0]:idx[1]]
			normalized := normalizePhone(match)
			if normalized != "" {
				return normalized
			}
		}
	}
	return ""
}

func isBoundedByNonAlnum(text string, start, end int) bool {
	if start > 0 {
		prev := string([]byte{text[start-1]})
		if alphaNumRegex.MatchString(prev) {
			return false
		}
	}
	if end < len(text) {
		next := string([]byte{text[end]})
		if alphaNumRegex.MatchString(next) {
			return false
		}
	}
	return true
}

func normalizePhone(raw string) string {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return ""
	}
	digitsOnly := make([]rune, 0, len(raw))
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			digitsOnly = append(digitsOnly, r)
		}
	}
	if len(digitsOnly) < 8 {
		return ""
	}
	digits := string(digitsOnly)
	if strings.Contains(raw, "+") {
		return "+" + digits
	}
	return digits
}

func parseGuestNameFromContext(context, email, phone string) string {
	context = cleanExtractedValue(context)
	if context == "" {
		return ""
	}
	context = guestLabelPrefixRe.ReplaceAllString(context, "")
	if email != "" {
		re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(email))
		if err == nil {
			context = re.ReplaceAllString(context, " ")
		}
	}
	if phone != "" {
		re, err := regexp.Compile(regexp.QuoteMeta(phone))
		if err == nil {
			context = re.ReplaceAllString(context, " ")
		}
	}
	context = phoneValueRegex.ReplaceAllString(context, " ")
	context = emailValueRegex.ReplaceAllString(context, " ")
	context = guestNoiseTokenRe.ReplaceAllString(context, " ")
	context = strings.NewReplacer(":", " ", ";", " ", "|", " ", "/", " ", "(", " ", ")", " ").Replace(context)
	context = cleanExtractedValue(context)
	if context == "" {
		return ""
	}

	nameTokens := nameTokenRegex.FindAllString(context, -1)
	if len(nameTokens) == 0 {
		return ""
	}

	if len(nameTokens) > 4 {
		nameTokens = nameTokens[:4]
	}
	name := strings.Join(nameTokens, " ")
	name = cleanExtractedValue(name)
	if len(name) < 2 {
		return ""
	}
	return name
}

func matchPatternList(patterns model.PatternList, value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(patterns) == 0 {
		return false
	}
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		re, err := regexp.Compile("(?i)" + pattern)
		if err != nil {
			continue
		}
		if re.MatchString(value) {
			return true
		}
	}
	return false
}

func parseRoomBooked(raw string) (uint32, string) {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return 0, ""
	}
	roomNumber := parseUint32(raw)
	roomTypeName := raw
	if roomNumber > 0 {
		roomTypeName = strings.TrimSpace(roomBookedPrefixRe.ReplaceAllString(raw, ""))
	}
	return roomNumber, roomTypeName
}

func parseUint32(raw string) uint32 {
	match := leadingNumberRegex.FindString(raw)
	if match == "" {
		return 0
	}
	value, err := strconv.ParseUint(match, 10, 32)
	if err != nil {
		return 0
	}
	return uint32(value)
}

func extractDateCandidates(raw string) []time.Time {
	matches := dateCandidateRegex.FindAllString(raw, -1)
	if len(matches) == 0 {
		return nil
	}
	res := make([]time.Time, 0, len(matches))
	seen := make(map[int64]struct{}, len(matches))
	for _, m := range matches {
		parsed := parseDateString(m)
		if parsed.IsZero() {
			continue
		}
		key := parsed.Unix()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		res = append(res, parsed)
	}
	return res
}

func (w *tenantWorker) ensureRoomType(ctx context.Context, roomTypeName string) (*model.RoomType, error) {
	roomTypeName = cleanExtractedValue(roomTypeName)
	if roomTypeName == "" {
		roomTypeName = "Unknown Room Type"
	}
	if w.roomTypeRepo == nil {
		return nil, fmt.Errorf("room type repository is not initialized")
	}

	rt, err := w.roomTypeRepo.SearchByName(ctx, w.tenant.SchemaName, roomTypeName)
	if err != nil {
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			rt = nil
		} else {
			w.logPro.Error("Search room type by name failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			return nil, err
		}
	}
	if rt != nil {
		return rt, nil
	}

	rid, err := w.idPro.NextID()
	if err != nil {
		w.logPro.Error("Generate room type id failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return nil, err
	}

	rt = model.NewRoomType(rid, cases.Title(language.Und).String(roomTypeName), 0, 0, 0, nil, nil)
	if err := w.roomTypeRepo.Create(ctx, w.tenant.SchemaName, rt); err != nil {
		w.logPro.Error("Create room type failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return nil, err
	}
	return rt, nil
}

func (w *tenantWorker) ensureSource(ctx context.Context, sourceName string) (*model.Source, error) {
	sourceName = cleanExtractedValue(sourceName)
	if sourceName == "" {
		sourceName = "Unknown Source"
	}
	if w.sourceRepo == nil {
		return nil, fmt.Errorf("source repository is not initialized")
	}

	source, err := w.sourceRepo.SearchByName(ctx, w.tenant.SchemaName, sourceName)
	if err != nil {
		if errors.Is(err, exceptions.ErrSourceNotFound) {
			source = nil
		} else {
			w.logPro.Error("Search source by name failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
			return nil, err
		}
	}
	if source != nil {
		return source, nil
	}

	sid, err := w.idPro.NextID()
	if err != nil {
		w.logPro.Error("Generate source id failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return nil, err
	}

	source = model.NewSource(sid, cases.Title(language.Und).String(sourceName), nil, nil)
	if err := w.sourceRepo.Create(ctx, w.tenant.SchemaName, source); err != nil {
		w.logPro.Error("Create source failed", port.String("schema", w.tenant.SchemaName), port.Error(err))
		return nil, err
	}
	return source, nil
}

func (w *tenantWorker) validateBookingRegexSettings() {
	fields := map[string]model.PatternList{
		"email_booking.subject":             w.settings.EmailBooking.Subject,
		"email_cancel.subject":              w.settings.EmailCancel.Subject,
		"email_cancel.bodies":               w.settings.EmailCancel.Bodies,
		"email_booking.booking_number":      w.settings.EmailBooking.BookingNumber,
		"email_booking.guest":               w.settings.EmailBooking.Guest,
		"email_booking.check_in":            w.settings.EmailBooking.CheckIn,
		"email_booking.check_out":           w.settings.EmailBooking.CheckOut,
		"email_booking.booked_on":           w.settings.EmailBooking.BookedOn,
		"email_booking.room_booked":         w.settings.EmailBooking.RoomsBooked,
		"email_booking.source":              w.settings.EmailBooking.Source,
		"email_booking.total_net_price":     w.settings.EmailBooking.TotalNetPrice,
		"email_booking.total_sell_price":    w.settings.EmailBooking.TotalSellPrice,
		"email_booking.guest_number":        w.settings.EmailBooking.GuestNumber,
		"email_booking.room_number":         w.settings.EmailBooking.RoomNumber,
		"email_booking.promotion_name":      w.settings.EmailBooking.PromotionName,
		"email_booking.meal_plan":           w.settings.EmailBooking.MealPlan,
		"email_booking.booking_conditions":  w.settings.EmailBooking.BookingConditions,
		"email_booking.booking_preferences": w.settings.EmailBooking.BookingPreferences,
	}

	for field, patterns := range fields {
		for _, pattern := range patterns {
			pattern = strings.TrimSpace(pattern)
			if pattern == "" {
				continue
			}
			if _, err := regexp.Compile("(?i)" + pattern); err != nil {
				w.logPro.Warn(
					"invalid email booking regex pattern",
					port.String("schema", w.tenant.SchemaName),
					port.String("field", field),
					port.String("pattern", pattern),
					port.Error(err),
				)
			}
		}
	}
}

func parseDateString(raw string) time.Time {
	raw = cleanExtractedValue(raw)
	if raw == "" {
		return time.Time{}
	}

	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("ICT", 7*60*60)
	}
	lowerRaw := strings.ToLower(raw)

	layouts := []struct {
		keyword string
		layout  string
	}{
		{"from", "Monday, January 2, 2006 from 15:04"},
		{"until", "Monday, January 2, 2006 until 15:04"},
		{"", "Monday, January 2, 2006 15:04"},
		{"", "Monday, January 2, 2006"},
		{"", "January 2, 2006 15:04"},
		{"", "January 2, 2006"},
		{"", "Jan 2, 2006 15:04"},
		{"", "Jan 2, 2006"},
		{"", "2 January 2006 15:04"},
		{"", "2 January 2006"},
		{"", "2 Jan 2006 15:04"},
		{"", "2 Jan 2006"},
		{"", "2006-01-02 15:04:05"},
		{"", "2006-01-02 15:04"},
		{"", "2006-01-02"},
		{"", "02/01/2006 15:04"},
		{"", "02/01/2006"},
		{"", "2/1/2006 15:04"},
		{"", "2/1/2006"},
		{"", "02-01-2006 15:04"},
		{"", "02-01-2006"},
	}

	for _, l := range layouts {
		if l.keyword != "" && !strings.Contains(lowerRaw, l.keyword) {
			continue
		}
		if t, err := time.ParseInLocation(l.layout, raw, loc); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

func parsePrice(raw string) float64 {
	processed := nonNumericPriceRegex.ReplaceAllString(raw, "")
	processed = strings.TrimSpace(processed)
	if processed == "" {
		return 0
	}

	lastDot := strings.LastIndex(processed, ".")
	lastComma := strings.LastIndex(processed, ",")

	switch {
	case lastDot >= 0 && lastComma >= 0:
		if lastDot > lastComma {
			processed = strings.ReplaceAll(processed, ",", "")
		} else {
			processed = strings.ReplaceAll(processed, ".", "")
			processed = strings.ReplaceAll(processed, ",", ".")
		}
	case lastComma >= 0:
		if decimalCommaSuffixRe.MatchString(processed) {
			processed = strings.ReplaceAll(processed, ".", "")
			processed = strings.ReplaceAll(processed, ",", ".")
		} else {
			processed = strings.ReplaceAll(processed, ",", "")
		}
	case lastDot >= 0:
		if decimalDotSuffixRe.MatchString(processed) {
			processed = strings.ReplaceAll(processed, ",", "")
		} else {
			processed = strings.ReplaceAll(processed, ".", "")
		}
	}

	price, err := strconv.ParseFloat(processed, 64)
	if err != nil {
		return 0
	}
	return price
}

func (w *tenantWorker) calculateBookingTotalNetPrice(ctx context.Context, booking *model.Booking) (float64, error) {
	if booking == nil || booking.RoomType == nil {
		return 0, nil
	}

	stayDays := utils.CalculateStayDays(booking.CheckIn, booking.CheckOut)
	roomCount := float64(booking.RoomNumber)
	if roomCount < 1 {
		roomCount = 1
	}

	if w.dailyRateRepo == nil {
		return booking.RoomType.Price * float64(stayDays) * roomCount, nil
	}

	checkInDate := utils.NormalizeDateOnly(booking.CheckIn)
	checkOutDate := checkInDate.AddDate(0, 0, int(stayDays))

	dailyTotalPrice, coveredDays, err := w.dailyRateRepo.SumFinalPriceByRoomTypeIDAndDateRange(ctx, w.tenant.SchemaName, booking.RoomTypeID, checkInDate, checkOutDate)
	if err != nil {
		w.logPro.Error(
			"sum daily room type rate by date range failed",
			port.String("schema", w.tenant.SchemaName),
			port.Int64("room_type_id", booking.RoomTypeID),
			port.Time("check_in", checkInDate),
			port.Time("check_out", checkOutDate),
			port.Error(err),
		)
		return 0, err
	}

	if coveredDays < 0 {
		coveredDays = 0
	}
	if coveredDays > stayDays {
		coveredDays = stayDays
	}

	missingDays := stayDays - coveredDays
	pricePerRoom := dailyTotalPrice + (float64(missingDays) * booking.RoomType.Price)

	return pricePerRoom * roomCount, nil
}

func calculateBackoff(attempt int) time.Duration {
	if attempt > maxReconnectAttempts {
		attempt = maxReconnectAttempts
	}
	delay := min(baseReconnectDelay*time.Duration(1<<uint(attempt)), maxReconnectDelay)
	return delay
}
