package hubspot

import (
	"encoding/json"
	"fmt"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
)

const (
	dateTimeMSFormat string = "2006-01-02T15:04:05.999Z"
	dateTimeMSNull   string = "1970-01-01T00:00:00Z"
)

type DateTimeMSString time.Time

func (d *DateTimeMSString) UnmarshalJSON(b []byte) error {
	var returnError = func() error {
		errortools.CaptureError(fmt.Sprintf("Cannot parse '%s' to DateTimeMSString", string(b)))
		return nil
	}

	var s string

	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("DateTimeMSString", string(b))
		return returnError()
	}

	if s == "" || s == dateTimeMSNull {
		d = nil
		return nil
	}

	_t, err := time.Parse(dateTimeMSFormat, s)
	if err != nil {
		return returnError()
	}

	*d = DateTimeMSString(_t)
	return nil
}

func (d *DateTimeMSString) MarshalJSON() ([]byte, error) {
	if d == nil {
		return nil, nil
	}

	return json.Marshal(time.Time(*d).Format(dateTimeMSFormat))
}

func (d *DateTimeMSString) ValuePtr() *time.Time {
	if d == nil {
		return nil
	}

	_d := time.Time(*d)
	return &_d
}

func (d DateTimeMSString) Value() time.Time {
	return time.Time(d)
}
