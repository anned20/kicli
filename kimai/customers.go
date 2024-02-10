package kimai

import (
	"encoding/json"
	"errors"
)

func (k *KimaiClient) GetCustomers() ([]Customer, error) {
	response, err := k.api.get("/customers")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var customers []Customer

	err = json.NewDecoder(response.Body).Decode(&customers)

	if err != nil {
		return nil, err
	}

	if len(customers) == 0 {
		return []Customer{}, nil
	}

	if customers[0].ID == 0 {
		return nil, errors.New("Malformed customers")
	}

	return customers, nil
}

func (k *KimaiClient) CreateCustomer(customer *Customer) (*Customer, error) {
	createCustomer := &CreateCustomer{
		Name:     customer.Name,
		Country:  customer.Country,
		Currency: customer.Currency,
		Timezone: customer.Timezone,
		Visible:  true,
	}

	response, err := k.api.post("/customers", createCustomer)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var newCustomer Customer

	err = json.NewDecoder(response.Body).Decode(&newCustomer)

	if err != nil {
		return nil, err
	}

	if newCustomer.ID == 0 {
		return nil, errors.New("Malformed customer")
	}

	return &newCustomer, nil
}
