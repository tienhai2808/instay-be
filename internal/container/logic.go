package container

import (
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
	mongoImpl "github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mongo/implement"
	ormImpl "github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/implement"
)

func (c *Container) initLogic() {
	c.baseRepo = ormImpl.NewBaseRepository(c.SQL.ORM())
	c.SchemaRepo = ormImpl.NewSchemaRepository(c.baseRepo)
	c.TenantRepo = ormImpl.NewTenantRepository(c.baseRepo)
	c.branchRepo = ormImpl.NewBranchRepository(c.baseRepo)
	c.UserRepo = ormImpl.NewUserRepository(c.baseRepo)
	c.RefreshTokenRepo = ormImpl.NewRefreshTokenRepository(c.baseRepo)
	c.DepartmentRepo = ormImpl.NewDepartmentRepository(c.baseRepo)
	c.serviceTypeRepo = ormImpl.NewServiceTypeRepository(c.baseRepo)
	c.serviceRepo = ormImpl.NewServiceRepository(c.baseRepo)
	c.serviceImageRepo = ormImpl.NewServiceImageRepository(c.baseRepo)
	c.supportCategoryRepo = ormImpl.NewSupportCategoryRepository(c.baseRepo)
	c.SourceRepo = ormImpl.NewSourceRepository(c.baseRepo)
	c.RoomTypeRepo = ormImpl.NewRoomTypeRepository(c.baseRepo)
	c.floorRepo = ormImpl.NewFloorRepository(c.baseRepo)
	c.roomRepo = ormImpl.NewRoomRepository(c.baseRepo)
	c.BookingRepo = ormImpl.NewBookingRepository(c.baseRepo)
	c.PermissionRepo = ormImpl.NewPermissionRepository(c.baseRepo)
	c.NotificationRepo = ormImpl.NewNotificationRepository(c.baseRepo)
	c.userNotificationRepo = ormImpl.NewUserNotificationRepository(c.baseRepo)
	c.bookingRoomRepo = ormImpl.NewBookingRoomRepository(c.baseRepo)
	c.orderServiceRepo = ormImpl.NewOrderServiceRepository(c.baseRepo)
	c.supportTicketRepo = ormImpl.NewSupportTicketRepository(c.baseRepo)
	c.DeviceTokenRepo = ormImpl.NewDeviceTokenRepository(c.baseRepo)
	c.pricingTagRepo = ormImpl.NewPricingTagRepository(c.baseRepo)
	c.pricingCalendarRepo = ormImpl.NewPricingCalendarRepository(c.baseRepo)
	c.roomTypePriceRepo = ormImpl.NewRoomTypePriceRepository(c.baseRepo)
	c.dailyRoomTypeRateRepo = ormImpl.NewDailyRoomTypeRateRepository(c.baseRepo)
	c.DailyRoomTypeRateRepo = c.dailyRoomTypeRateRepo
	c.departmentPermissionRepo = ormImpl.NewDepartmentPermissionRepository(c.baseRepo)
	c.reviewRepo = ormImpl.NewReviewRepository(c.baseRepo)
	c.guestRepo = ormImpl.NewGuestRepository(c.baseRepo)
	c.customerRepo = ormImpl.NewCustomerRepository(c.baseRepo)
	c.eventRepo = ormImpl.NewEventRepository(c.baseRepo)
	c.chatRepo = ormImpl.NewChatRepository(c.baseRepo)
	c.messageRepo = mongoImpl.NewMessageRepository(c.NoSQL.DB())

	c.recalculateUC = recalculateUC.NewRecalculateUseCase(c.LogPro, c.RoomTypeRepo, c.roomTypePriceRepo, c.pricingCalendarRepo, c.dailyRoomTypeRateRepo)

	c.fileUC = fileUC.NewFileUseCase(c.LogPro, c.StoragePro)
	c.authUC = authUC.NewAuthUseCase(c.cfg.Token, c.txPro, c.LogPro, c.tokenPro, c.CachePro, c.MQPro, c.StoragePro, c.UserRepo, c.RefreshTokenRepo)
	c.branchUC = branchUC.NewBranchUseCase(c.LogPro, c.IDPro, c.branchRepo)
	c.userUC = userUC.NewUserUseCase(c.txPro, c.MQPro, c.LogPro, c.IDPro, c.CachePro, c.UserRepo, c.DepartmentRepo, c.branchRepo, c.RefreshTokenRepo)
	c.departmentUC = departmentUC.NewDepartmentUseCase(c.LogPro, c.IDPro, c.txPro, c.CachePro, c.DepartmentRepo, c.RefreshTokenRepo, c.PermissionRepo, c.UserRepo, c.departmentPermissionRepo)
	c.serviceTypeUC = serviceTypeUC.NewServiceTypeUseCase(c.LogPro, c.IDPro, c.serviceTypeRepo, c.serviceRepo)
	c.serviceUC = serviceUC.NewServiceUseCase(c.LogPro, c.IDPro, c.txPro, c.MQPro, c.StoragePro, c.serviceRepo, c.serviceImageRepo, c.branchRepo, c.UserRepo)
	c.supportCategoryUC = supportCategoryUC.NewSupportCategoryUseCase(c.LogPro, c.IDPro, c.supportCategoryRepo)
	c.sourceUC = sourceUC.NewSourceUseCase(c.LogPro, c.IDPro, c.SourceRepo)
	c.roomTypeUC = roomTypeUC.NewRoomTypeUseCase(c.LogPro, c.IDPro, c.recalculateUC, c.RoomTypeRepo, c.branchRepo, c.UserRepo, c.dailyRoomTypeRateRepo, c.roomRepo)
	c.floorUC = floorUC.NewFloorUseCase(c.LogPro, c.IDPro, c.floorRepo, c.UserRepo)
	c.roomUC = roomUC.NewRoomUseCase(c.LogPro, c.IDPro, c.roomRepo, c.floorRepo, c.RoomTypeRepo, c.UserRepo)
	c.bookingUC = bookingUC.NewBookingUseCase(c.LogPro, c.IDPro, c.txPro, c.CachePro, c.AIPro, c.BookingRepo, c.bookingRoomRepo, c.branchRepo, c.UserRepo, c.roomRepo, c.RoomTypeRepo, c.dailyRoomTypeRateRepo)
	c.permissionUC = permissionUC.NewPermissionUseCase(c.LogPro, c.PermissionRepo)
	c.bookingRoomUC = bookingRoomUC.NewBookingRoomUseCase(c.LogPro, c.IDPro, c.txPro, c.CachePro, c.tokenPro, c.BookingRepo, c.roomRepo, c.RoomTypeRepo, c.UserRepo, c.bookingRoomRepo, c.guestRepo)
	c.orderServiceUC = orderServiceUC.NewOrderServiceUseCase(c.LogPro, c.IDPro, c.txPro, c.PNPro, c.MQPro, c.serviceRepo, c.bookingRoomRepo, c.orderServiceRepo, c.NotificationRepo, c.UserRepo, c.DeviceTokenRepo)
	c.supportTicketUC = supportTicketUC.NewSupportTicketUseCase(c.LogPro, c.IDPro, c.txPro, c.MQPro, c.supportTicketRepo, c.bookingRoomRepo, c.supportCategoryRepo, c.NotificationRepo, c.UserRepo, c.DeviceTokenRepo)
	c.deviceTokenUC = deviceTokenUC.NewDeviceTokenUseCase(c.LogPro, c.DeviceTokenRepo)
	c.pricingTagUC = pricingTagUC.NewPricingTagUseCase(c.LogPro, c.IDPro, c.pricingTagRepo)
	c.pricingCalendarUC = pricingCalendarUC.NewPricingCalendarUseCase(c.LogPro, c.txPro, c.pricingCalendarRepo, c.RoomTypeRepo, c.recalculateUC)
	c.roomTypePriceUC = roomTypePriceUC.NewRoomTypePriceUseCase(c.LogPro, c.IDPro, c.roomTypePriceRepo, c.recalculateUC)
	c.dashboardUC = dashboardUC.NewDashboardUseCase(c.LogPro, c.BookingRepo, c.reviewRepo, c.guestRepo, c.supportTicketRepo, c.orderServiceRepo, c.messageRepo, c.chatRepo, c.UserRepo)
	c.reviewUC = reviewUC.NewReviewUseCase(c.LogPro, c.IDPro, c.AIPro, c.reviewRepo, c.UserRepo)
	c.guestUC = guestUC.NewGuestUseCase(c.LogPro, c.IDPro, c.KYCPro, c.guestRepo, c.bookingRoomRepo, c.UserRepo)
	c.customerUC = customerUC.NewCustomerUseCase(c.LogPro, c.IDPro, c.customerRepo, c.bookingRoomRepo, c.UserRepo)
	c.eventUC = eventUC.NewEventUseCase(c.LogPro, c.IDPro, c.txPro, c.MQPro, c.eventRepo, c.branchRepo, c.UserRepo, c.NotificationRepo, c.bookingRoomRepo, c.DeviceTokenRepo)
	c.notificationUC = notificationUC.NewNotificationUseCase(c.LogPro, c.NotificationRepo, c.userNotificationRepo, c.UserRepo)
	c.messageUC = messageUC.NewMessageUseCase(c.LogPro, c.IDPro, c.txPro, c.chatRepo, c.bookingRoomRepo, c.messageRepo)
	c.tenantUC = tenantUC.NewTenantUseCase(c.LogPro, c.MQPro, c.TenantRepo)
	c.chatUC = chatUC.NewChatUseCase(c.LogPro, c.chatRepo, c.messageRepo, c.bookingRoomRepo, c.UserRepo, c.serviceRepo, c.supportCategoryRepo)
}
