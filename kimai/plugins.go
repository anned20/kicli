package kimai

import (
	"encoding/json"
	"errors"
)

func (k *KimaiClient) GetPlugins() ([]Plugin, error) {
	response, err := k.api.get("/plugins")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var plugins []Plugin

	err = json.NewDecoder(response.Body).Decode(&plugins)

	if err != nil {
		return nil, err
	}

	if len(plugins) == 0 {
		return []Plugin{}, nil
	}

	if plugins[0].Name == "" {
		return nil, errors.New("Malformed customers")
	}

	return plugins, nil
}
