package exceptions

import "github.com/InstaySystem/is_v2-be/internal/domain/constants"

var (
	ErrLoginFailed = NewException(constants.CodeLoginFailed, "Incorrect username or password")

	ErrInvalidToken = NewException(constants.CodeInvalidToken, "Invalid or expired token")

	ErrBadRequest = NewException(constants.CodeBadRequest, "Invalid data")

	ErrUnauth = NewException(constants.CodeUnauth, "Unauthorized")

	ErrForbidden = NewException(constants.CodeForbidden, "Forbidden")

	ErrInvalidPassword = NewException(constants.CodeInvalidPassword, "Incorrect password")

	ErrInvalidUser = NewException(constants.CodeNoRefreshToken, "Please login again")

	ErrUserNotFound = NewException(constants.CodeUserNotFound, "User not found")

	ErrEmailDoesNotExist = NewException(constants.CodeEmailDoesNotExist, "Email does not exist")

	ErrTooManyAttempts = NewException(constants.CodeTooManyAttempts, "Too many attempts")

	ErrInvalidOTP = NewException(constants.CodeInvalidOTP, "Invalid or expired OTP")

	ErrEmailAlreadyExists = NewException(constants.CodeEmailAlreadyExists, "Email already exists")

	ErrNameAlreadyExists = NewException(constants.CodeNameAlreadyExists, "Name already exists")

	ErrPhoneAlreadyExists = NewException(constants.CodePhoneAlreadyExists, "Phone already exists")

	ErrDepartmentNotFound = NewException(constants.CodeDepartmentNotFound, "Department not found")

	ErrTenantNotFound = NewException(constants.CodeTenantNotFound, "Tenant not found")

	ErrInvalidID = NewException(constants.CodeInvalidID, "Invalid id")

	ErrProtectedRecord = NewException(constants.CodeProtectedRecord, "Protected record")

	ErrHasUserNotFound = NewException(constants.CodeHasUserNotFound, "Has user not found")

	ErrNeedAdmin = NewException(constants.CodeNeedAdmin, "Need 1 active administrator")

	ErrUsernameAlreadyExists = NewException(constants.CodeUsernameAlreadyExists, "Username already exists")

	ErrInternalServerError = NewException(constants.CodeInternalError, "Internal server error")

	ErrInvalidDomain = NewException(constants.CodeInvalidDomain, "Invalid domain")

	ErrDomainNotAllowed = NewException(constants.CodeDomainNotAllowed, "Domain not allowed")

	ErrDomainAlreadyExists = NewException(constants.CodeBadRequest, "Domain already exists")

	ErrSchemaNameAlreadyExists = NewException(constants.CodeBadRequest, "Schema name already exists")

	ErrServiceTypeAlreadyExists = NewException(constants.CodeServiceTypeAlreadyExists, "Service type already exists")

	ErrServiceTypeNotFound = NewException(constants.CodeServiceTypeNotFound, "Service type not found")

	ErrServiceAlreadyExists = NewException(constants.CodeServiceAlreadyExists, "Service already exists")

	ErrServiceNotFound = NewException(constants.CodeServiceNotFound, "Service not found")

	ErrSupportCategoryAlreadyExists = NewException(constants.CodeSupportCategoryAlreadyExists, "Support category already exists")

	ErrSupportCategoryNotFound = NewException(constants.CodeSupportCategoryNotFound, "Support category not found")

	ErrSourceAlreadyExists = NewException(constants.CodeSourceAlreadyExists, "Source already exists")

	ErrSourceNotFound = NewException(constants.CodeSourceNotFound, "Source not found")

	ErrRoomTypeAlreadyExists = NewException(constants.CodeRoomTypeAlreadyExists, "Room type already exists")

	ErrBookingAlreadyExists = NewException(constants.CodeBookingAlreadyExists, "Booking already exists")

	ErrRoomTypeNotFound = NewException(constants.CodeRoomTypeNotFound, "Room type not found")

	ErrBookingNotFound = NewException(constants.CodeBookingNotFound, "Booking not found")

	ErrFloorAlreadyExists = NewException(constants.CodeFloorAlreadyExists, "Floor already exists")

	ErrRoomAlreadyExists = NewException(constants.CodeRoomAlreadyExists, "Room already exists")

	ErrFloorNotFound = NewException(constants.CodeFloorNotFound, "Floor not found")

	ErrHasPermissionNotFound = NewException(constants.CodeHasPermissionNotFound, "Has permission not found")

	ErrBookingExpired = NewException(constants.CodeBookingExpired, "Booking expired")

	ErrMaxRoomReached = NewException(constants.CodeMaxRoomReached, "Max room reached")

	ErrRoomNotFound = NewException(constants.CodeRoomNotFound, "Room not found")

	ErrRoomCurrentlyOccupied = NewException(constants.CodeRoomCurrentlyOccupied, "Room currently occupied")

	ErrOrderServiceNotFound = NewException(constants.CodeOrderServiceNotFound, "Order service not found")

	ErrInvalidStatus = NewException(constants.CodeInvalidStatus, "Invalid status")

	ErrRefreshTokenNotFound = NewException(constants.CodeRefreshTokenNotFound, "Refresh token not found")

	ErrBookingRoomNotFound = NewException(constants.CodeBookingRoomNotFound, "Booking room not found")

	ErrSupportTicketNotFound = NewException(constants.CodeSupportTicketNotFound, "Support ticket not found")

	ErrDeviceTokenNotFound = NewException(constants.CodeDeviceTokenNotFound, "Device token not found")

	ErrReviewAlreadyExists = NewException(constants.CodeReviewAlreadyExists, "Guest have already left review")

	ErrReviewNotFound = NewException(constants.CodeReviewNotFound, "Review not found")

	ErrGuestNotFound = NewException(constants.CodeGuestNotFound, "Guest not found")

	ErrBookingSelectedRoom = NewException(constants.CodeBookingSelectedRoom, "Booking selected room")

	ErrBookingCantUpdate = NewException(constants.CodeBookingCantUpdate, "Booking can't update")

	ErrBookingRoomsNotMatched = NewException(constants.CodeBookingRoomsNotMatched, "Booking rooms not matched")

	ErrGuestInBookingRoomAlreadyExists = NewException(constants.CodeGuestInBookingRoomAlreadyExists, "Guest provided accommodation information")

	ErrMessageNotFound = NewException(constants.CodeMessageNotFound, "Message not found")

	ErrChatNotFound = NewException(constants.CodeChatNotFound, "Chat not found")

	ErrChatAlreadyExists = NewException(constants.CodeChatAlreadyExists, "Chat already exists")

	ErrPricingTagAlreadyExists = NewException(constants.CodePricingTagAlreadyExists, "Pricing tag already exists")

	ErrPricingTagNotFound = NewException(constants.CodePricingTagNotFound, "Pricing tag not found")

	ErrPricingCalendarDateConflict = NewException(constants.CodePricingCalendarDateConflict, "Date is assigned to multiple pricing tags")

	ErrRoomTypePriceAlreadyExists = NewException(constants.CodeRoomTypePriceAlreadyExists, "Room type price already exists")

	ErrRoomTypePriceNotFound = NewException(constants.CodeRoomTypePriceNotFound, "Room type price not found")

	ErrEventAlreadyExists = NewException(constants.CodeEventAlreadyExists, "Event already exists")

	ErrEventNotFound = NewException(constants.CodeEventNotFound, "Event not found")

	ErrBranchAlreadyExists = NewException(constants.CodeBranchAlreadyExists, "Branch already exists")

	ErrBranchNotFound = NewException(constants.CodeBranchNotFound, "Branch not found")

	ErrRoomTypeNotInBranch = NewException(constants.CodeRoomTypeNotInBranch, "Room type not in branch")

	ErrRoomNotInBookingBranch = NewException(constants.CodeRoomNotInBookingBranch, "Room not in booking branch")

	ErrServiceNotInBranch = NewException(constants.CodeServiceNotInBranch, "Service not in branch")
)

type Exception struct {
	Code    int
	Message string
	Data    any
}

func NewException(code int, msg string) *Exception {
	return &Exception{
		Code:    code,
		Message: msg,
	}
}

func (e *Exception) Error() string {
	return e.Message
}

func (e *Exception) WithData(data any) *Exception {
	return &Exception{
		Code:    e.Code,
		Message: e.Message,
		Data:    data,
	}
}

func (e *Exception) Is(target error) bool {
	t, ok := target.(*Exception)
	if !ok {
		return false
	}

	return e.Code == t.Code
}
