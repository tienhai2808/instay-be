package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToCustomerEntity(customer *model.Customer) *entity.CustomerEntity {
	if customer == nil {
		return nil
	}

	return &entity.CustomerEntity{
		ID:            customer.ID,
		CreatedAt:     customer.CreatedAt,
		BirthYear:     customer.BirthYear,
		Phone:         customer.Phone,
		Gender:        string(customer.Gender),
		BookingRoomID: customer.BookingRoomID,
	}
}

func ToCustomer(customer *entity.CustomerEntity) *model.Customer {
	if customer == nil {
		return nil
	}

	return &model.Customer{
		ID:            customer.ID,
		CreatedAt:     customer.CreatedAt,
		BirthYear:     customer.BirthYear,
		Phone:         customer.Phone,
		Gender:        model.CustomerGender(customer.Gender),
		BookingRoomID: customer.BookingRoomID,
		BookingRoom:   ToBookingRoom(customer.BookingRoom),
	}
}

func ToCustomers(customers []*entity.CustomerEntity) []*model.Customer {
	if len(customers) == 0 {
		return make([]*model.Customer, 0)
	}

	customersRes := make([]*model.Customer, 0, len(customers))
	for _, customer := range customers {
		customersRes = append(customersRes, ToCustomer(customer))
	}

	return customersRes
}
