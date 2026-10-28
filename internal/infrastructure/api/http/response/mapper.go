package response

import (
	"net/http"

	"github.com/InstaySystem/is_v2-be/internal/domain/constants"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/gin-gonic/gin"
)

func StatusFromExceptionCode(code int) int {
	switch code {
	case constants.CodeBadRequest,
		constants.CodeLoginFailed,
		constants.CodeInvalidToken,
		constants.CodeInvalidPassword,
		constants.CodeEmailDoesNotExist,
		constants.CodeInvalidOTP,
		constants.CodeNeedAdmin,
		constants.CodeInvalidID,
		constants.CodeInvalidDomain,
		constants.CodeBookingExpired,
		constants.CodeMaxRoomReached:
		return http.StatusBadRequest
	case constants.CodeUnauth:
		return http.StatusUnauthorized
	case constants.CodeNoRefreshToken,
		constants.CodeForbidden,
		constants.CodeDomainNotAllowed:
		return http.StatusForbidden
	case constants.CodeTooManyAttempts:
		return http.StatusTooManyRequests
	case constants.CodeUserNotFound,
		constants.CodeDepartmentNotFound,
		constants.CodeTenantNotFound,
		constants.CodeServiceTypeNotFound,
		constants.CodeServiceNotFound,
		constants.CodeSupportCategoryNotFound,
		constants.CodeSourceNotFound,
		constants.CodeRoomTypeNotFound,
		constants.CodeBookingNotFound,
		constants.CodeFloorNotFound,
		constants.CodeHasPermissionNotFound,
		constants.CodeRoomNotFound,
		constants.CodeOrderServiceNotFound,
		constants.CodeRefreshTokenNotFound,
		constants.CodeBookingRoomNotFound,
		constants.CodeSupportTicketNotFound,
		constants.CodeDeviceTokenNotFound,
		constants.CodeReviewNotFound,
		constants.CodeGuestNotFound,
		constants.CodeMessageNotFound,
		constants.CodeChatNotFound,
		constants.CodePricingTagNotFound,
		constants.CodeRoomTypePriceNotFound,
		constants.CodeEventNotFound,
		constants.CodeBranchNotFound:
		return http.StatusNotFound
	case constants.CodeEmailAlreadyExists,
		constants.CodeNameAlreadyExists,
		constants.CodePhoneAlreadyExists,
		constants.CodeProtectedRecord,
		constants.CodeHasUserNotFound,
		constants.CodeUsernameAlreadyExists,
		constants.CodeServiceTypeAlreadyExists,
		constants.CodeServiceAlreadyExists,
		constants.CodeSupportCategoryAlreadyExists,
		constants.CodeSourceAlreadyExists,
		constants.CodeRoomTypeAlreadyExists,
		constants.CodeBookingAlreadyExists,
		constants.CodeFloorAlreadyExists,
		constants.CodeRoomAlreadyExists,
		constants.CodeRoomCurrentlyOccupied,
		constants.CodeInvalidStatus,
		constants.CodeReviewAlreadyExists,
		constants.CodeBookingSelectedRoom,
		constants.CodeBookingCantUpdate,
		constants.CodeBookingRoomsNotMatched,
		constants.CodeGuestInBookingRoomAlreadyExists,
		constants.CodeChatAlreadyExists,
		constants.CodePricingTagAlreadyExists,
		constants.CodePricingCalendarDateConflict,
		constants.CodeEventAlreadyExists,
		constants.CodeRoomTypePriceAlreadyExists,
		constants.CodeBranchAlreadyExists,
		constants.CodeRoomTypeNotInBranch,
		constants.CodeRoomNotInBookingBranch,
		constants.CodeServiceNotInBranch:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func ISEResponse(c *gin.Context) {
	APIResponse(
		c,
		http.StatusInternalServerError,
		exceptions.ErrInternalServerError.Code,
		exceptions.ErrInternalServerError.Message,
		nil,
	)
}
