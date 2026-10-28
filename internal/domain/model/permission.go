package model

type PermissionName string

const (
	PermBookingManagement                    PermissionName = "booking:management"
	PermOrderServiceManagement               PermissionName = "order_service:management"
	PermSupportTicketManagement              PermissionName = "support_ticket:management"
	PermGuestManagement                      PermissionName = "guest:management"
	PermEventManagement                      PermissionName = "event:management"
	PermBookingManagementDisplayNameEN                      = "Management booking"
	PermBookingManagementDisplayNameVI                      = "Quản lý booking"
	PermOrderServiceManagementDisplayNameEN                 = "Management order service"
	PermOrderServiceManagementDisplayNameVI                 = "Quản lý đơn dịch vụ"
	PermSupportTicketManagementDisplayNameEN                = "Management support ticket"
	PermSupportTicketManagementDisplayNameVI                = "Quản lý phiếu hỗ trợ"
	PermGuestManagementDisplayNameEN                        = "Management accommodation "
	PermGuestManagementDisplayNameVI                        = "Quản lý lưu trú"
	PermEventManagementDisplayNameEN                        = "Event management"
	PermEventManagementDisplayNameVI                        = "Quản lý sự kiện"
)

type Permission struct {
	ID            int64
	Name          PermissionName
	DisplayNameVI string
	DisplayNameEN string

	Departments []*Department
}

func (p PermissionName) IsValid() bool {
	switch p {
	case PermBookingManagement,
		PermOrderServiceManagement,
		PermSupportTicketManagement,
		PermGuestManagement,
		PermEventManagement:
		return true
	default:
		return false
	}
}
