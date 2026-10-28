package translator

import (
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	PgErrCodeForeignKeyViolation         = "23503"
	PgErrCodeUniqueViolation             = "23505"
	PgErrCodeForeignKeyViolationRestrict = "23001"
	ConstraintUserDepartmentFK           = "fk_users_department"
	ConstraintUserEmailUnique            = "users_email_key"
	ConstraintUserUsernameUnique         = "users_username_key"
	ConstraintUserPhoneUnique            = "users_phone_key"

	ConstraintDepartmentNameUnique  = "departments_name_key"
	ConstraintDepartmentPhoneUnique = "departments_phone_key"

	ConstraintServiceTypeSlugUnique = "service_types_slug_key"
	ConstraintServiceSlugUnique     = "services_slug_key"
	ConstraintEventSlugUnique       = "events_slug_key"

	ConstraintSupportCategoryNameUnique = "support_categories_name_key"
	ConstraintSourceNameUnique          = "sources_name_key"
	ConstraintRoomTypeNameUnique        = "room_types_name_key"
	ConstraintFloorNameUnique           = "floors_branch_id_name_key"
	ConstraintBranchNameUnique          = "branches_name_key"
	ConstraintRoomNameUnique            = "rooms_name_key"
	ConstraintBookingNumberUnique       = "bookings_booking_number_key"
	ConstraintBookingSelectedRoomUnique = "booking_rooms_booking_id_room_id_key"
	ConstraintBookingRoomCodeUnique     = "booking_rooms_code_key"
	ConstraintChatUnique                = "chats_booking_room_id_department_id"
	ConstraintChatUniqueKey             = "chats_booking_room_id_department_id_key"
	ConstraintReviewUnique              = "reviews_booking_room_id"
	ConstraintReviewUniqueKey           = "reviews_booking_room_id_key"

	ConstraintGuestUnique         = "guests_passport_number_booking_room_id_key"
	ConstraintPricingTagUnique    = "pricing_tags_name_key"
	ConstraintRoomTypePriceUnique = "room_type_prices_room_type_id_pricing_tag_id_key"
	ConstraintTenantDomainUnique  = "tenants_domain_key"
	ConstraintTenantSchemaUnique  = "tenants_schema_name_key"

	ConstraintServiceTypeDepartmentFK  = "fk_service_types_department"
	ConstraintSupportCategoryDeptFK    = "fk_support_categories_department"
	ConstraintServiceServiceTypeFK     = "fk_services_service_type"
	ConstraintBookingBranchFK          = "fk_bookings_branch"
	ConstraintBookingSourceFK          = "fk_bookings_source"
	ConstraintBookingRoomTypeFK        = "fk_bookings_room_type"
	ConstraintRoomFloorFK              = "fk_rooms_floor"
	ConstraintRoomRoomTypeFK           = "fk_rooms_room_type"
	ConstraintGuestBookingRoomFK       = "fk_guests_booking_room"
	ConstraintCustomerBookingRoomFK    = "fk_customers_booking_room"
	ConstraintReviewBookingRoomFK      = "fk_reviews_booking_room"
	ConstraintPricingCalendarTagFK     = "fk_pricing_calendars_pricing_tag"
	ConstraintRoomTypePriceRoomTypeFK  = "fk_room_type_prices_room_type"
	ConstraintRoomTypePricePricingFK   = "fk_room_type_prices_pricing_tag"
	ConstraintBranchRoomTypeBranchFK   = "fk_branch_room_types_branch"
	ConstraintBranchRoomTypeRoomTypeFK = "fk_branch_room_types_room_type"
	ConstraintBranchServiceBranchFK    = "fk_branch_services_branch"
	ConstraintBranchServiceServiceFK   = "fk_branch_services_service"
	ConstraintUserBranchUserFK         = "fk_user_branches_user"
	ConstraintUserBranchBranchFK       = "fk_user_branches_branch"
)

var postgresWriteFKToException = map[string]error{
	ConstraintUserDepartmentFK:         exceptions.ErrDepartmentNotFound,
	ConstraintServiceTypeDepartmentFK:  exceptions.ErrDepartmentNotFound,
	ConstraintSupportCategoryDeptFK:    exceptions.ErrDepartmentNotFound,
	ConstraintServiceServiceTypeFK:     exceptions.ErrServiceTypeNotFound,
	ConstraintBookingBranchFK:          exceptions.ErrBranchNotFound,
	ConstraintBookingSourceFK:          exceptions.ErrSourceNotFound,
	ConstraintBookingRoomTypeFK:        exceptions.ErrRoomTypeNotFound,
	ConstraintRoomFloorFK:              exceptions.ErrFloorNotFound,
	ConstraintRoomRoomTypeFK:           exceptions.ErrRoomTypeNotFound,
	ConstraintGuestBookingRoomFK:       exceptions.ErrBookingRoomNotFound,
	ConstraintCustomerBookingRoomFK:    exceptions.ErrBookingRoomNotFound,
	ConstraintReviewBookingRoomFK:      exceptions.ErrUnauth,
	ConstraintPricingCalendarTagFK:     exceptions.ErrPricingTagNotFound,
	ConstraintRoomTypePriceRoomTypeFK:  exceptions.ErrRoomTypeNotFound,
	ConstraintRoomTypePricePricingFK:   exceptions.ErrPricingTagNotFound,
	ConstraintBranchRoomTypeBranchFK:   exceptions.ErrBranchNotFound,
	ConstraintBranchRoomTypeRoomTypeFK: exceptions.ErrRoomTypeNotFound,
	ConstraintBranchServiceBranchFK:    exceptions.ErrBranchNotFound,
	ConstraintBranchServiceServiceFK:   exceptions.ErrServiceNotFound,
	ConstraintUserBranchUserFK:         exceptions.ErrUserNotFound,
	ConstraintUserBranchBranchFK:       exceptions.ErrBranchNotFound,
}

var postgresUniqueToException = map[string]error{
	ConstraintUserEmailUnique:    exceptions.ErrEmailAlreadyExists,
	ConstraintUserUsernameUnique: exceptions.ErrUsernameAlreadyExists,
	ConstraintUserPhoneUnique:    exceptions.ErrPhoneAlreadyExists,

	ConstraintDepartmentNameUnique:  exceptions.ErrNameAlreadyExists,
	ConstraintDepartmentPhoneUnique: exceptions.ErrPhoneAlreadyExists,

	ConstraintServiceTypeSlugUnique: exceptions.ErrServiceTypeAlreadyExists,
	ConstraintServiceSlugUnique:     exceptions.ErrServiceAlreadyExists,
	ConstraintEventSlugUnique:       exceptions.ErrEventAlreadyExists,

	ConstraintSupportCategoryNameUnique: exceptions.ErrSupportCategoryAlreadyExists,
	ConstraintSourceNameUnique:          exceptions.ErrSourceAlreadyExists,
	ConstraintRoomTypeNameUnique:        exceptions.ErrRoomTypeAlreadyExists,
	ConstraintFloorNameUnique:           exceptions.ErrFloorAlreadyExists,
	ConstraintBranchNameUnique:          exceptions.ErrBranchAlreadyExists,
	ConstraintRoomNameUnique:            exceptions.ErrRoomAlreadyExists,
	ConstraintBookingNumberUnique:       exceptions.ErrBookingAlreadyExists,
	ConstraintBookingSelectedRoomUnique: exceptions.ErrBookingSelectedRoom,
	ConstraintBookingRoomCodeUnique:     exceptions.ErrBookingSelectedRoom,
	ConstraintChatUnique:                exceptions.ErrChatAlreadyExists,
	ConstraintChatUniqueKey:             exceptions.ErrChatAlreadyExists,
	ConstraintReviewUnique:              exceptions.ErrReviewAlreadyExists,
	ConstraintReviewUniqueKey:           exceptions.ErrReviewAlreadyExists,

	ConstraintGuestUnique:         exceptions.ErrGuestInBookingRoomAlreadyExists,
	ConstraintPricingTagUnique:    exceptions.ErrPricingTagAlreadyExists,
	ConstraintRoomTypePriceUnique: exceptions.ErrRoomTypePriceAlreadyExists,
	ConstraintTenantDomainUnique:  exceptions.ErrDomainAlreadyExists,
	ConstraintTenantSchemaUnique:  exceptions.ErrSchemaNameAlreadyExists,
}

func PostgresTranslate(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if mappedErr, ok := postgresUniqueToException[pgErr.ConstraintName]; ok {
			return mappedErr
		}

		switch pgErr.Code {
		case PgErrCodeForeignKeyViolation, PgErrCodeForeignKeyViolationRestrict:
			if mappedErr, ok := postgresWriteFKToException[pgErr.ConstraintName]; ok {
				return mappedErr
			}
			return exceptions.ErrProtectedRecord
		case PgErrCodeUniqueViolation:
			return err
		}
	}

	return err
}

func PostgresTranslateDelete(err error) error {
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if mappedErr, ok := postgresUniqueToException[pgErr.ConstraintName]; ok {
			return mappedErr
		}
		switch pgErr.Code {
		case PgErrCodeForeignKeyViolation, PgErrCodeForeignKeyViolationRestrict:
			return exceptions.ErrProtectedRecord
		}
	}
	return err
}
