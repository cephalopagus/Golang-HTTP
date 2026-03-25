package http

import (
	"encoding/json"
	"errors"
	"time"
)

type TaskDTO struct {
	Title       string
	Description string
}

type ErrorDTO struct {
	Message string
	Time    time.Time
}
type CompleteDTO struct {
	Complete bool
}

func (e ErrorDTO) ErrorToString() string {
	b, err := json.MarshalIndent(e, "", "	")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (t TaskDTO) Validate() error {
	if t.Title == "" {
		return errors.New("empty title")
	}
	if t.Description == "" {
		return errors.New("empty description")
	}
	return nil
}
