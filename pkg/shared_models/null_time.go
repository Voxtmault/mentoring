package shared_models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type CustomNullTime struct {
	sql.NullTime
}

const timeFormat = time.DateTime

// MarshalJSON ensures that if the time is null, it returns an empty string.
func (c CustomNullTime) MarshalJSON() ([]byte, error) {
	if !c.Valid {
		return json.Marshal("") // Return an empty string if the time is null
	}
	return json.Marshal(c.Time.Format(timeFormat)) // Format the time as a string
}

// UnmarshalJSON parses the time from JSON, allowing empty strings to be treated as null.
func (c *CustomNullTime) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		c.Valid = false
		return nil
	}
	parsedTime, err := time.Parse(timeFormat, str)
	if err != nil {
		return err
	}
	c.Time = parsedTime
	c.Valid = true
	return nil
}
