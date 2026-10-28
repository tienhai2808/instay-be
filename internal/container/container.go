package container

import (
	"fmt"
	"log"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	authUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/auth"
	bookingUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/booking"
	bookingRoomUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/bookingroom"
	branchUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/branch"
	chatUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/chat"
	customerUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/customer"
	dashboardUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/dashboard"
	departmentUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/department"
	deviceTokenUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/devicetoken"
	eventUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/event"
	fileUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/file"
	floorUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/floor"
	guestUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/guest"
	messageUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/message"
	notificationUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/notification"
	orderServiceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/orderservice"
	permissionUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/permission"
	pricingCalendarUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/pricingcalendar"
	pricingTagUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/pricingtag"
	recalculateUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/recalculate"
	reviewUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/review"
	roomUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/room"
	roomTypeUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/roomtype"
	roomTypePriceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/roomtypeprice"
	serviceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/service"
	serviceTypeUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/servicetype"
	sourceUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/source"
	supportCategoryUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/supportcategory"
	supportTicketUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/supportticket"
	tenantUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/tenant"
	userUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/user"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	httpHdl "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/handler"
	httpMid "github.com/InstaySystem/is_v2-be/internal/infrastructure/api/http/middleware"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/config"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/initialization"
	ormImpl "github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/implement"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/ai"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/fcm"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/gcs"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/rabbitmq"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/redis"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/smtp"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/sonyflake"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/provider/zap"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/realtime/ws"
)

type Container struct {
	cfg                      *config.Config
	log                      *initialization.Log
	SQL                      *initialization.SQL
	NoSQL                    *initialization.NoSQL
	cache                    *initialization.Cache
	mq                       *initialization.MessageQueue
	stor                     *initialization.Storage
	idGen                    *initialization.IDGen
	pn                       *initialization.PushNotification
	gRPCClients              *initialization.GRPCClients
	txPro                    port.TxProvider
	tokenPro                 port.TokenProvider
	apiKeyPro                port.APIKeyProvider
	MQPro                    port.MessageQueueProvider
	CachePro                 port.CacheProvider
	MailerPro                port.MailerProvider
	StoragePro               port.StorageProvider
	PNPro                    port.PushNotificationProvider
	KYCPro                   port.KYCProvider
	AIPro                    port.AIProvider
	IDPro                    port.IDProvider
	LogPro                   port.LogProvider
	baseRepo                 ormImpl.BaseRepository
	SchemaRepo               repository.SchemeRepository
	TenantRepo               repository.TenantRepository
	branchRepo               repository.BranchRepository
	UserRepo                 repository.UserRepository
	RefreshTokenRepo         repository.RefreshTokenRepository
	DepartmentRepo           repository.DepartmentRepository
	serviceTypeRepo          repository.ServiceTypeRepository
	serviceRepo              repository.ServiceRepository
	serviceImageRepo         repository.ServiceImageRepository
	supportCategoryRepo      repository.SupportCategoryRepository
	SourceRepo               repository.SourceRepository
	BookingRepo              repository.BookingRepository
	RoomTypeRepo             repository.RoomTypeRepository
	roomRepo                 repository.RoomRepository
	floorRepo                repository.FloorRepository
	PermissionRepo           repository.PermissionRepository
	DeviceTokenRepo          repository.DeviceTokenRepository
	NotificationRepo         repository.NotificationRepository
	userNotificationRepo     repository.UserNotificationRepository
	bookingRoomRepo          repository.BookingRoomRepository
	orderServiceRepo         repository.OrderServiceRepository
	supportTicketRepo        repository.SupportTicketRepository
	pricingTagRepo           repository.PricingTagRepository
	pricingCalendarRepo      repository.PricingCalendarRepository
	roomTypePriceRepo        repository.RoomTypePriceRepository
	dailyRoomTypeRateRepo    repository.DailyRoomTypeRateRepository
	DailyRoomTypeRateRepo    repository.DailyRoomTypeRateRepository
	departmentPermissionRepo repository.DepartmentPermissionRepository
	reviewRepo               repository.ReviewRepository
	guestRepo                repository.GuestRepository
	customerRepo             repository.CustomerRepository
	eventRepo                repository.EventRepository
	chatRepo                 repository.ChatRepository
	messageRepo              repository.MessageRepository
	fileUC                   fileUC.FileUseCase
	authUC                   authUC.AuthUseCase
	branchUC                 branchUC.BranchUseCase
	userUC                   userUC.UserUseCase
	departmentUC             departmentUC.DepartmentUseCase
	serviceTypeUC            serviceTypeUC.ServiceTypeUseCase
	serviceUC                serviceUC.ServiceUseCase
	supportCategoryUC        supportCategoryUC.SupportCategoryUseCase
	sourceUC                 sourceUC.SourceUseCase
	bookingUC                bookingUC.BookingUseCase
	roomTypeUC               roomTypeUC.RoomTypeUseCase
	dashboardUC              dashboardUC.DashboardUseCase
	roomUC                   roomUC.RoomUseCase
	floorUC                  floorUC.FloorUseCase
	permissionUC             permissionUC.PermissionUseCase
	bookingRoomUC            bookingRoomUC.BookingRoomUseCase
	orderServiceUC           orderServiceUC.OrderServiceUseCase
	supportTicketUC          supportTicketUC.SupportTicketUseCase
	deviceTokenUC            deviceTokenUC.DeviceTokenUseCase
	pricingTagUC             pricingTagUC.PricingTagUseCase
	pricingCalendarUC        pricingCalendarUC.PricingCalendarUseCase
	recalculateUC            recalculateUC.RecalculateUseCase
	roomTypePriceUC          roomTypePriceUC.RoomTypePriceUseCase
	reviewUC                 reviewUC.ReviewUseCase
	guestUC                  guestUC.GuestUseCase
	customerUC               customerUC.CustomerUseCase
	eventUC                  eventUC.EventUseCase
	notificationUC           notificationUC.NotificationUseCase
	messageUC                messageUC.MessageUseCase
	tenantUC                 tenantUC.TenantUseCase
	chatUC                   chatUC.ChatUseCase
	FileHTTPHdl              *httpHdl.FileHandler
	AuthHTTPHdl              *httpHdl.AuthHandler
	BranchHTTPHdl            *httpHdl.BranchHandler
	UserHTTPHdl              *httpHdl.UserHandler
	DepartmentHTTPHdl        *httpHdl.DepartmentHandler
	ServiceTypeHTTPHdl       *httpHdl.ServiceTypeHandler
	ServiceHTTPHdl           *httpHdl.ServiceHandler
	SupportCategoryHTTPHdl   *httpHdl.SupportCategoryHandler
	SourceHTTPHdl            *httpHdl.SourceHandler
	BookingHTTPHdl           *httpHdl.BookingHandler
	RoomTypeHTTPHdl          *httpHdl.RoomTypeHandler
	FloorHTTPHdl             *httpHdl.FloorHandler
	RoomHTTPHdl              *httpHdl.RoomHandler
	PermissionHTTPHdl        *httpHdl.PermissionHandler
	PricingTagHTTPHdl        *httpHdl.PricingTagHandler
	PricingCalendarHTTPHdl   *httpHdl.PricingCalenderHandler
	RoomTypePriceHTTPHdl     *httpHdl.RoomTypePriceHandler
	BookingRoomHTTPHdl       *httpHdl.BookingRoomHandler
	OrderServiceHTTPHdl      *httpHdl.OrderServiceHandler
	SupportTicketHTTPHdl     *httpHdl.SupportTicketHandler
	DeviceTokenHTTPHdl       *httpHdl.DeviceTokenHandler
	DashboardHTTPHdl         *httpHdl.DashboardHandler
	ReviewHTTPHdl            *httpHdl.ReviewHandler
	GuestHTTPHdl             *httpHdl.GuestHandler
	CustomerHTTPHdl          *httpHdl.CustomerHandler
	EventHTTPHdl             *httpHdl.EventHandler
	NotificationHTTPHdl      *httpHdl.NotificationHandler
	TenantHTTPHdl            *httpHdl.TenantHandler
	ChatHTTPHdl              *httpHdl.ChatHandler
	CQAHTTPHdl               *httpHdl.CQAHandler
	WSHTTPHdl                *httpHdl.WSHandler
	CtxHTTPMid               *httpMid.ContextMiddleware
	DomainHTTPMid            *httpMid.DomainMiddleware
	AuthHTTPMid              *httpMid.AuthMiddleware
	APIKeyHTTPMid            *httpMid.APIKeyMiddleware
	HubWS                    *ws.HubWS
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		cfg: cfg,
	}
}

func (c *Container) InitServer() (err error) {
	defer func() {
		if err != nil {
			c.Cleanup()
		}
	}()

	if err := c.initCore(); err != nil {
		return err
	}

	c.initLogic()
	c.initAPI()

	if c.HubWS != nil {
		go c.HubWS.Run()
	}

	return nil
}

func (c *Container) InitSeed() (err error) {
	defer func() {
		if err != nil {
			c.Cleanup()
		}
	}()

	c.log, err = initialization.InitLog(c.cfg.Log)
	if err != nil {
		return err
	}

	c.SQL, err = initialization.InitSQL(c.cfg.SQL)
	if err != nil {
		return err
	}

	c.idGen, err = initialization.InitIDGen()
	if err != nil {
		return err
	}

	c.LogPro = zap.NewLogProvider(c.log.Logger())
	c.IDPro = sonyflake.NewIDProvider(c.idGen.Generator())

	c.baseRepo = ormImpl.NewBaseRepository(c.SQL.ORM())
	c.UserRepo = ormImpl.NewUserRepository(c.baseRepo)
	c.SchemaRepo = ormImpl.NewSchemaRepository(c.baseRepo)
	c.PermissionRepo = ormImpl.NewPermissionRepository(c.baseRepo)

	return nil
}

func (c *Container) InitConsumer() (err error) {
	defer func() {
		if err != nil {
			c.Cleanup()
		}
	}()

	c.log, err = initialization.InitLog(c.cfg.Log)
	if err != nil {
		return err
	}

	c.mq, err = initialization.InitMessageQueue(c.cfg.MQ)
	if err != nil {
		return err
	}

	c.pn, err = initialization.InitPushNotification()
	if err != nil {
		return err
	}

	c.stor, err = initialization.InitStorage(c.cfg.Storage)
	if err != nil {
		return err
	}

	c.LogPro = zap.NewLogProvider(c.log.Logger())
	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Connection(), c.log.Logger())
	c.MailerPro = smtp.NewMailerProvider(c.cfg.Mailer)
	c.PNPro = fcm.NewPushNotificationProvider(c.pn.Client())
	c.StoragePro = gcs.NewStorageProvider(c.cfg.Storage, c.stor.Client())

	return nil
}

func (c *Container) InitScheduler() (err error) {
	defer func() {
		if err != nil {
			c.Cleanup()
		}
	}()

	c.log, err = initialization.InitLog(c.cfg.Log)
	if err != nil {
		return err
	}

	c.SQL, err = initialization.InitSQL(c.cfg.SQL)
	if err != nil {
		return err
	}

	c.LogPro = zap.NewLogProvider(c.log.Logger())

	c.baseRepo = ormImpl.NewBaseRepository(c.SQL.ORM())
	c.SchemaRepo = ormImpl.NewSchemaRepository(c.baseRepo)

	return nil
}

func (c *Container) InitWorker() (err error) {
	defer func() {
		if err != nil {
			c.Cleanup()
		}
	}()

	c.mq, err = initialization.InitMessageQueue(c.cfg.MQ)
	if err != nil {
		return err
	}

	c.cache, err = initialization.InitCache(c.cfg.Cache)
	if err != nil {
		return err
	}

	c.log, err = initialization.InitLog(c.cfg.Log)
	if err != nil {
		return err
	}

	c.SQL, err = initialization.InitSQL(c.cfg.SQL)
	if err != nil {
		return err
	}

	c.idGen, err = initialization.InitIDGen()
	if err != nil {
		return err
	}

	aiAddr := fmt.Sprintf("%s:%d", c.cfg.Services.AIHost, c.cfg.Services.AIPort)
	ca := initialization.ClientAddresses{
		AIAddr: aiAddr,
	}

	c.gRPCClients, err = initialization.InitClients(ca)
	if err != nil {
		return err
	}

	c.MQPro = rabbitmq.NewMessageQueueProvider(c.mq.Connection(), c.log.Logger())
	c.CachePro = redis.NewCacheProvider(c.cache.Client())
	c.LogPro = zap.NewLogProvider(c.log.Logger())
	c.IDPro = sonyflake.NewIDProvider(c.idGen.Generator())
	c.AIPro = ai.NewAIProvider(c.gRPCClients.AIClient)

	c.baseRepo = ormImpl.NewBaseRepository(c.SQL.ORM())
	c.TenantRepo = ormImpl.NewTenantRepository(c.baseRepo)
	c.BookingRepo = ormImpl.NewBookingRepository(c.baseRepo)
	c.SourceRepo = ormImpl.NewSourceRepository(c.baseRepo)
	c.RoomTypeRepo = ormImpl.NewRoomTypeRepository(c.baseRepo)
	c.UserRepo = ormImpl.NewUserRepository(c.baseRepo)
	c.DeviceTokenRepo = ormImpl.NewDeviceTokenRepository(c.baseRepo)
	c.NotificationRepo = ormImpl.NewNotificationRepository(c.baseRepo)
	c.DepartmentRepo = ormImpl.NewDepartmentRepository(c.baseRepo)

	return nil
}

func (c *Container) Cleanup() {
	if c.HubWS != nil {
		c.HubWS.Stop()
	}

	if c.SQL != nil {
		c.SQL.Close()
	}
	if c.NoSQL != nil {
		c.NoSQL.Close()
	}
	if c.mq != nil {
		c.mq.Close()
	}
	if c.cache != nil {
		c.cache.Close()
	}
	if c.stor != nil {
		c.stor.Close()
	}
	if c.log != nil {
		c.log.Close()
	}
	if c.gRPCClients != nil {
		c.gRPCClients.Close()
	}

	log.Println("Container cleaned successfully")
}
