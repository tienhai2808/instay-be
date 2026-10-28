package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToServiceTypeResponse(svcType *model.ServiceType) *response.ServiceTypeResponse {
	if svcType == nil {
		return nil
	}

	return &response.ServiceTypeResponse{
		ID:           svcType.ID,
		Name:         svcType.Name,
		Icon:         svcType.Icon,
		CreatedAt:    svcType.CreatedAt,
		UpdatedAt:    svcType.UpdatedAt,
		CreatedBy:    ToBasicUserResponse(svcType.CreatedBy),
		UpdatedBy:    ToBasicUserResponse(svcType.UpdatedBy),
		Department:   ToBasicDepartmentResponse(svcType.Department),
		ServiceCount: svcType.ServicesCount,
	}
}

func ToBasicServiceImageResponse(img *model.ServiceImage) *response.BasicServiceImageResponse {
	if img == nil {
		return nil
	}

	return &response.BasicServiceImageResponse{
		Key: img.Key,
	}
}

func ToExtendedServiceResponse(svc *model.Service) *response.ExtendedServiceResponse {
	if svc == nil {
		return nil
	}

	return &response.ExtendedServiceResponse{
		ID:              svc.ID,
		Name:            svc.Name,
		Price:           svc.Price,
		DiscountType:    svc.DiscountType,
		DiscountValue:   svc.DiscountValue,
		DiscountStartAt: svc.DiscountStartAt,
		DiscountEndAt:   svc.DiscountEndAt,
		Description:     svc.Description,
		ServiceType:     ToBasicServiceTypeResponse(svc.ServiceType),
		Branches:        ToBasicBranchesResponse(svc.Branches),
		Images:          ToServiceImagesResponse(svc.ServiceImages),
	}
}

func ToServiceResponse(svc *model.Service) *response.ServiceResponse {
	if svc == nil {
		return nil
	}

	var thumbnail *model.ServiceImage
	if len(svc.ServiceImages) > 0 {
		thumbnail = svc.ServiceImages[0]
	}

	return &response.ServiceResponse{
		ID:              svc.ID,
		Price:           svc.Price,
		DiscountType:    svc.DiscountType,
		DiscountValue:   svc.DiscountValue,
		DiscountStartAt: svc.DiscountStartAt,
		DiscountEndAt:   svc.DiscountEndAt,
		Name:            svc.Name,
		Slug:            svc.Slug,
		Thumbnail:       ToBasicServiceImageResponse(thumbnail),
		Branches:        ToBasicBranchesResponse(svc.Branches),
	}
}

func ToServicesResponse(svcs []*model.Service) []*response.ServiceResponse {
	if len(svcs) == 0 {
		return make([]*response.ServiceResponse, 0)
	}

	svcsRes := make([]*response.ServiceResponse, 0, len(svcs))
	for _, svc := range svcs {
		svcsRes = append(svcsRes, ToServiceResponse(svc))
	}

	return svcsRes
}

func ToOrderServiceResponse(os *model.OrderService) *response.OrderServiceResponse {
	if os == nil {
		return nil
	}

	return &response.OrderServiceResponse{
		ID:           os.ID,
		BookingRoom:  ToBasicBookingRoomResponse(os.BookingRoom),
		Service:      ToServiceResponse(os.Service),
		Quantity:     os.Quantity,
		UnitPrice:    os.UnitPrice,
		TotalPrice:   os.TotalPrice,
		Status:       os.Status,
		GuestNote:    os.GuestNote,
		StaffNote:    os.StaffNote,
		CancelReason: os.CancelReason,
		EstimatedAt:  os.EstimatedAt,
		CreatedAt:    os.CreatedAt,
		UpdatedAt:    os.UpdatedAt,
		CreatedBy:    ToBasicUserResponse(os.CreatedBy),
		UpdatedBy:    ToBasicUserResponse(os.UpdatedBy),
	}
}

func ToServiceTypeWithServicesQuery(query request.ServiceTypeWithActiveServicesQuery) model.ServiceTypeWithServicesQuery {
	return model.ServiceTypeWithServicesQuery{
		Search:        query.Search,
		ServiceTypeID: query.ServiceTypeID,
		BranchID:      query.BranchID,
	}
}

func ToBasicOrderServiceResponse(orderSvc *model.OrderService) *response.BasicOrderServiceResponse {
	if orderSvc == nil {
		return nil
	}

	return &response.BasicOrderServiceResponse{
		ID:          orderSvc.ID,
		Service:     ToBasicServiceResponse(orderSvc.Service),
		BookingRoom: ToBasicBookingRoomResponse(orderSvc.BookingRoom),
		UnitPrice:   orderSvc.UnitPrice,
		Quantity:    orderSvc.Quantity,
		TotalPrice:  orderSvc.TotalPrice,
		Status:      orderSvc.Status,
		EstimatedAt: orderSvc.EstimatedAt,
		GuestNote:   orderSvc.GuestNote,
	}
}

func ToOrderServiceQuery(query request.OrderServicePaginationQuery) model.OrderServiceQuery {
	return model.OrderServiceQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Sort:     query.Sort,
		Order:    query.Order,
		Status:   query.Status,
		FromDate: query.FromDate,
		ToDate:   query.ToDate,
	}
}

func ToBasicOrderServicesResponse(ordeSvcs []*model.OrderService) []*response.BasicOrderServiceResponse {
	if len(ordeSvcs) == 0 {
		return make([]*response.BasicOrderServiceResponse, 0)
	}

	orderSvcsRes := make([]*response.BasicOrderServiceResponse, 0, len(ordeSvcs))
	for _, orderSvc := range ordeSvcs {
		orderSvcsRes = append(orderSvcsRes, ToBasicOrderServiceResponse(orderSvc))
	}

	return orderSvcsRes
}

func ToServiceTypesResponse(svcTypes []*model.ServiceType) []*response.ServiceTypeResponse {
	if len(svcTypes) == 0 {
		return make([]*response.ServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*response.ServiceTypeResponse, 0, len(svcTypes))
	for _, svcType := range svcTypes {
		serviceTypesRes = append(serviceTypesRes, ToServiceTypeResponse(svcType))
	}

	return serviceTypesRes
}

func ToBasicServiceTypeResponse(svcType *model.ServiceType) *response.BasicServiceTypeResponse {
	if svcType == nil {
		return nil
	}

	return &response.BasicServiceTypeResponse{
		ID:   svcType.ID,
		Name: svcType.Name,
		Icon: svcType.Icon,
		Slug: svcType.Slug,
	}
}

func ToBasicServiceTypesResponse(svcTypes []*model.ServiceType) []*response.BasicServiceTypeResponse {
	if len(svcTypes) == 0 {
		return make([]*response.BasicServiceTypeResponse, 0)
	}

	serviceTypesRes := make([]*response.BasicServiceTypeResponse, 0, len(svcTypes))
	for _, svcType := range svcTypes {
		serviceTypesRes = append(serviceTypesRes, ToBasicServiceTypeResponse(svcType))
	}

	return serviceTypesRes
}

func ToServiceImageResponse(svcImg *model.ServiceImage) *response.ServiceImageResponse {
	if svcImg == nil {
		return nil
	}

	return &response.ServiceImageResponse{
		Key:         svcImg.Key,
		IsThumbnail: svcImg.IsThumbnail,
		SortOrder:   svcImg.SortOrder,
	}
}

func ToServiceImagesResponse(svcImgs []*model.ServiceImage) []*response.ServiceImageResponse {
	if len(svcImgs) == 0 {
		return make([]*response.ServiceImageResponse, 0)
	}

	svcImgsRes := make([]*response.ServiceImageResponse, 0, len(svcImgs))
	for _, svcImg := range svcImgs {
		svcImgsRes = append(svcImgsRes, ToServiceImageResponse(svcImg))
	}

	return svcImgsRes
}

func ToServiceDetailsResponse(svc *model.Service) *response.ServiceDetailsResponse {
	if svc == nil {
		return nil
	}

	return &response.ServiceDetailsResponse{
		ID:              svc.ID,
		Name:            svc.Name,
		Price:           svc.Price,
		DiscountType:    svc.DiscountType,
		DiscountValue:   svc.DiscountValue,
		DiscountStartAt: svc.DiscountStartAt,
		DiscountEndAt:   svc.DiscountEndAt,
		IsActive:        svc.IsActive,
		Description:     svc.Description,
		ServiceType:     ToBasicServiceTypeResponse(svc.ServiceType),
		Branches:        ToBasicBranchesResponse(svc.Branches),
		Images:          ToServiceImagesResponse(svc.ServiceImages),
		CreatedAt:       svc.CreatedAt,
		UpdatedAt:       svc.UpdatedAt,
		CreatedBy:       ToBasicUserResponse(svc.CreatedBy),
		UpdatedBy:       ToBasicUserResponse(svc.UpdatedBy),
	}
}

func ToBasicServiceResponse(svc *model.Service) *response.BasicServiceResponse {
	if svc == nil {
		return nil
	}

	return &response.BasicServiceResponse{
		ID:   svc.ID,
		Name: svc.Name,
		Slug: svc.Slug,
	}
}

func ToSimpleServiceTypeResponse(st *model.ServiceType) *response.SimpleServiceTypeResponse {
	if st == nil {
		return nil
	}

	return &response.SimpleServiceTypeResponse{
		ID:       st.ID,
		Name:     st.Name,
		Icon:     st.Icon,
		Services: ToServicesResponse(st.Services),
	}
}

func ToSimpleServiceTypesResponse(sts []*model.ServiceType) []*response.SimpleServiceTypeResponse {
	if len(sts) == 0 {
		return make([]*response.SimpleServiceTypeResponse, 0)
	}

	stsRes := make([]*response.SimpleServiceTypeResponse, 0, len(sts))
	for _, st := range sts {
		stsRes = append(stsRes, ToSimpleServiceTypeResponse(st))
	}

	return stsRes
}

func ToSimpleServiceResponse(svc *model.Service) *response.SimpleServiceResponse {
	if svc == nil {
		return nil
	}

	return &response.SimpleServiceResponse{
		ID:          svc.ID,
		Name:        svc.Name,
		ServiceType: ToBasicServiceTypeResponse(svc.ServiceType),
		Branches:    ToBasicBranchesResponse(svc.Branches),
		IsActive:    svc.IsActive,
		Price:       svc.Price,
		CreatedAt:   svc.CreatedAt,
	}
}

func ToSimpleServicesResponse(svcs []*model.Service) []*response.SimpleServiceResponse {
	if len(svcs) == 0 {
		return make([]*response.SimpleServiceResponse, 0)
	}

	svcsRes := make([]*response.SimpleServiceResponse, 0, len(svcs))
	for _, svc := range svcs {
		svcsRes = append(svcsRes, ToSimpleServiceResponse(svc))
	}

	return svcsRes
}

func ToSimpleOrderServiceResponse(orderSvc *model.OrderService) *response.SimpleOrderServiceResponse {
	if orderSvc == nil {
		return nil
	}

	return &response.SimpleOrderServiceResponse{
		ID:           orderSvc.ID,
		Service:      ToServiceResponse(orderSvc.Service),
		Quantity:     orderSvc.Quantity,
		UnitPrice:    orderSvc.UnitPrice,
		TotalPrice:   orderSvc.TotalPrice,
		Status:       orderSvc.Status,
		GuestNote:    orderSvc.GuestNote,
		StaffNote:    orderSvc.StaffNote,
		CancelReason: orderSvc.CancelReason,
		EstimatedAt:  orderSvc.EstimatedAt,
		CreatedAt:    orderSvc.CreatedAt,
	}
}

func ToSimpleOrderServicesResponse(orderServices []*model.OrderService) []*response.SimpleOrderServiceResponse {
	if len(orderServices) == 0 {
		return make([]*response.SimpleOrderServiceResponse, 0)
	}

	orderServicesRes := make([]*response.SimpleOrderServiceResponse, 0, len(orderServices))
	for _, orderSvc := range orderServices {
		orderServicesRes = append(orderServicesRes, ToSimpleOrderServiceResponse(orderSvc))
	}

	return orderServicesRes
}

func ToServiceTypeQuery(query request.ServiceTypePaginationQuery) model.ServiceTypeQuery {
	return model.ServiceTypeQuery{
		Page:         query.Page,
		Limit:        query.Limit,
		Search:       query.Search,
		DepartmentID: query.DepartmentID,
	}
}

func ToServiceQuery(query request.ServicePaginationQuery) model.ServiceQuery {
	return model.ServiceQuery{
		Page:          query.Page,
		Limit:         query.Limit,
		Search:        query.Search,
		ServiceTypeID: query.ServiceTypeID,
		BranchID:      query.BranchID,
		Sort:          query.Sort,
		Order:         query.Order,
		IsActive:      query.IsActive,
		HasDiscount:   query.HasDiscount,
	}
}
