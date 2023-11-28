package boterror

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = NewError("Resident not found", errors.New("not_found"))
)

type BotError struct {
	Msg string `json:"message"`
	Err error  `json:"-"`
}

func (a *BotError) Error() string {
	return fmt.Sprintf("%s", a.Msg)
}

func NewError(msg string, err error) *BotError {
	return &BotError{
		Msg: msg,
		Err: err,
	}
}

func ParseErrToText(err error) string {
	switch {
	case errors.Is(err, ErrNotFound):
		return "Резидент не был найден"
	}

	return "Произошла внутренняя ошибка на сервере"
}
