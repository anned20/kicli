package kimai

import "time"

// dateTimeFormat is the format used for timestamps in the Kimai API
const dateTimeFormat = "2006-01-02T15:04:05-0700"
const fallbackDatetimeformat = "2006-01-02 15:04:05"
const fallbackDateformat = "2006-01-02"

type KimaiClient struct {
	api *API
}

func NewKimaiClient(baseURL string, token string) *KimaiClient {
	return &KimaiClient{
		api: NewAPI(baseURL, token),
	}
}

// Convert the timestamps to time.Time objects.
func normalizeTimestamps(object startAndEnd) (startAndEnd, error) {
	// Convert the timestamps to time.Time
	if object.GetRawStart() != "" {
		start, err := normalizeTimestamp(object.GetRawStart())

		if err != nil {
			return nil, err
		}

		object.SetStart(start)
	}

	if object.GetRawEnd() != "" {
		end, err := normalizeTimestamp(object.GetRawEnd())

		if err != nil {
			return nil, err
		}

		object.SetEnd(end)
	}

	return object, nil
}

func normalizeTimestamp(timestamp string) (*time.Time, error) {
	if timestamp == "" {
		return nil, nil
	}

	t, err := time.Parse(dateTimeFormat, timestamp)

	if err != nil {
		t, err = time.Parse(fallbackDatetimeformat, timestamp)

		if err != nil {
			t, err = time.Parse(fallbackDateformat, timestamp)

			if err != nil {
				return nil, err
			}
		}
	}

	return &t, nil
}
