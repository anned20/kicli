package kimai

import "time"

// dateFormat is the format used for timestamps in the Kimai API
const dateFormat = "2006-01-02T15:04:05-0700"

type KimaiClient struct {
	api *API
}

func NewKimaiClient(baseURL string, username string, token string) *KimaiClient {
	return &KimaiClient{
		api: NewAPI(baseURL, username, token),
	}
}

// Convert the timestamps to time.Time objects.
func normalizeTimestamps(object startAndEnd) (startAndEnd, error) {
	// Convert the timestamps to time.Time
	if object.GetRawStart() != "" {
		start, err := time.Parse(dateFormat, object.GetRawStart())

		if err != nil {
			return nil, err
		}

		object.SetStart(&start)
	}

	if object.GetRawEnd() != "" {
		end, err := time.Parse(dateFormat, object.GetRawEnd())

		if err != nil {
			return nil, err
		}

		object.SetEnd(&end)
	}

	return object, nil
}
