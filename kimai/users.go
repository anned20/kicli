package kimai

import (
	"encoding/json"
	"errors"
)

func (k *KimaiClient) GetMe() (*User, error) {
	response, err := k.api.get("/users/me")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var user User

	err = json.NewDecoder(response.Body).Decode(&user)

	if err != nil {
		return nil, err
	}

	if user.ID == 0 {
		return nil, errors.New("Malformed user")
	}

	return &user, nil
}
