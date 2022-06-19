package kimai

import (
	"encoding/json"
	"errors"
)

func (k *KimaiClient) GetVersion() (*Version, error) {
	response, err := k.api.get("/version")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var version Version

	err = json.NewDecoder(response.Body).Decode(&version)

	if err != nil {
		return nil, err
	}

	if version.Version == "" {
		return nil, errors.New("Malformed version")
	}

	return &version, nil
}
