package kimai

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type object struct {
	RawStart string
	RawEnd   string

	Start *time.Time
	End   *time.Time
}

func (t *object) GetRawStart() string {
	return t.RawStart
}

func (t *object) GetRawEnd() string {
	return t.RawEnd
}

func (t *object) SetStart(start *time.Time) {
	t.Start = start
}

func (t *object) SetEnd(end *time.Time) {
	t.End = end
}

func TestNormalizeTimestamps(t *testing.T) {
	tests := map[string]struct {
		obj            *object
		createExpected func() *object
		errString      string
	}{
		"empty": {
			obj: &object{},
			createExpected: func() *object {
				return &object{}
			},
		},
		"start only": {
			obj: &object{
				RawStart: dateFormat,
			},
			createExpected: func() *object {
				dt, _ := time.Parse(dateFormat, dateFormat)

				return &object{
					Start: &dt,
				}
			},
		},
		"end only": {
			obj: &object{
				RawEnd: dateFormat,
			},
			createExpected: func() *object {
				dt, _ := time.Parse(dateFormat, dateFormat)

				return &object{
					End: &dt,
				}
			},
		},
		"start and end": {
			obj: &object{
				RawStart: dateFormat,
				RawEnd:   dateFormat,
			},
			createExpected: func() *object {
				dt, _ := time.Parse(dateFormat, dateFormat)

				return &object{
					Start: &dt,
					End:   &dt,
				}
			},
		},
		"start and end with invalid date": {
			obj: &object{
				RawStart: "invalid",
				RawEnd:   dateFormat,
			},
			errString: "cannot parse",
		},
		"start and end with invalid date 2": {
			obj: &object{
				RawStart: dateFormat,
				RawEnd:   "invalid",
			},
			errString: "cannot parse",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := normalizeTimestamps(test.obj)

			if test.errString != "" {
				assert.ErrorContains(t, err, test.errString)

				return
			} else {
				assert.NoError(t, err)
			}

			object, ok := result.(*object)

			assert.True(t, ok)

			expected := test.createExpected()

			assert.NoError(t, err)
			assert.Equal(t, expected.Start, object.Start)
			assert.Equal(t, expected.End, object.End)
		})
	}
}
