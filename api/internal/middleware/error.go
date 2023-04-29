package middleware

import (
	"net/http"

	"github.com/parwin-pp/todo-application/internal/httperror"
	"github.com/uptrace/bunrouter"
)

func NewErrorHandler(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		err := next(w, req)
		if err == nil {
			return nil
		}

		switch err := err.(type) {
		case httperror.Error: // already a HTTPError
			w.WriteHeader(err.Status)
			_ = bunrouter.JSON(w, err)
		default:
			httpErr := httperror.New(500, "500", "Internal Server Error")
			w.WriteHeader(httpErr.Status)
			_ = bunrouter.JSON(w, httpErr)
		}

		return err
	}
}
