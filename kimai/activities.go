package kimai

import (
	"encoding/json"
	"errors"
)

func (k *KimaiClient) GetActivities() ([]Activity, error) {
	response, err := k.api.get("/activities")

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var activities []Activity

	err = json.NewDecoder(response.Body).Decode(&activities)

	if err != nil {
		return nil, err
	}

	if len(activities) == 0 {
		return []Activity{}, nil
	}

	if activities[0].ID == 0 {
		return nil, errors.New("Malformed activities")
	}

	return activities, nil
}
