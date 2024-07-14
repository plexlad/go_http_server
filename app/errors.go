package main

const (
  RouteRegistrationError int = iota
  RouteNotFoundError
  InternalServerError
)

type BaseError struct {
  code    int
  message string
}

func (e BaseError) Error() string {
  return e.message
}

func NewError(errorCode int, message string) error {
  return error(BaseError{
    errorCode,
    message,
  })
}
