package httperror

import "fmt"

var (
	ErrUnauthorized   = New(401, "401", "unauthorized, please login again")
	ErrInvalidRequest = New(400, "400", "invalid request")
	ErrInternalServer = New(500, "500", "something went wrong, please try again later")
)

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

func New(status int, code, msg string, args ...interface{}) Error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return Error{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}
