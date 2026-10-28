package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToCustomerResponse(customer *model.Customer) *response.CustomerResponse {
	if customer == nil {
		return nil
	}

	return &response.CustomerResponse{
		ID:        customer.ID,
		BirthYear: customer.BirthYear,
		Gender:    customer.Gender,
		Phone:     customer.Phone,
		CreatedAt: customer.CreatedAt,
	}
}

func ToCustomerQuery(query request.CustomerPaginationQuery) model.CustomerQuery {
	return model.CustomerQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Sort:   query.Sort,
		Order:  query.Order,
		Search: query.Search,
	}
}

func ToCustomersResponse(customers []*model.Customer) []*response.CustomerResponse {
	if len(customers) == 0 {
		return make([]*response.CustomerResponse, 0)
	}

	customersRes := make([]*response.CustomerResponse, 0, len(customers))
	for _, customer := range customers {
		customersRes = append(customersRes, ToCustomerResponse(customer))
	}

	return customersRes
}
