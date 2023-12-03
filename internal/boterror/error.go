package boterror

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound                 = NewError("Resident not found", errors.New("not_found"))
	ErrIncorrectCallbackData    = NewError("Incorrect Callback Data", errors.New("incorrect_callback"))
	ErrIncorrectAdminFirstInput = NewError("Must be 4 input values", errors.New("incorrect_input"))
	ErrUniqueViolation          = NewError("Violation must be unique", errors.New("non_unique_value"))
	ErrForeignKeyViolation      = NewError("Foreign Key Violation", errors.New("foreign_key_violation "))
	ErrIsNotAdmin               = NewError("The user is not an admin", errors.New("not_admin"))
	ErrUserNotMember            = NewError("User not member", errors.New("user_not_member"))
	ErrInternalError            = NewError("Internal error", errors.New("internal_error"))
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
	case errors.Is(err, ErrIncorrectAdminFirstInput):
		return "Должно быть введено минимум 4 слова при первом вводе [1]"
	case errors.Is(err, ErrUniqueViolation):
		return "Телеграм пользователя должен быть уникальным [1]"
	case errors.Is(err, ErrForeignKeyViolation):
		return "TODO: Что-то с внешнеим ключом"
	case errors.Is(err, ErrIsNotAdmin):
		return "Доступ ограничен"
	case errors.Is(err, ErrUserNotMember):
		return "Пользователь не является участником группы"
	}

	return "Произошла внутренняя ошибка на сервере"
}
