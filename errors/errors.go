// package errors makes it easier to associate HTTP response codes with Go errors.
package errors

// Error encapsulates an error message with an HTTP status code.
type Error struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func New(code int, message string) *Error {
	e := new(Error)
	e.StatusCode = code
	e.Message = message
	return e
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Code() int {
	return e.StatusCode
}
