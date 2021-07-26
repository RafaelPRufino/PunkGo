package analytics

import (
	"fmt"
	"time"
)

// ErrGeneric Error
type ErrGeneric struct {
	When time.Time
	What string
}

func (e ErrGeneric) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

type ErrAuthorization struct {
	When time.Time
	What string
}

func (e ErrAuthorization) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

type ErrResponseAuthorization struct {
	When time.Time
	What string
}

func (e ErrResponseAuthorization) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

type ErrMalformedRequest struct {
	When time.Time
	What string
}

func (e ErrMalformedRequest) Error() string {
	return fmt.Sprintf("%v: %v", e.When, e.What)
}

func NewError(errorName string) error {
	errs := map[string]error{
		"MalformedRequest": ErrMalformedRequest{
			time.Now().UTC(),
			"Malformed Request",
		},
		"Unauthenticated": ErrResponseAuthorization{
			time.Now().UTC(),
			"UnAuthenticated",
		},
		"Unauthorized": ErrAuthorization{
			time.Now().UTC(),
			"Unauthorized",
		},
		"Default": ErrGeneric{
			time.Now().UTC(),
			errorName,
		},
	}

	err := errs[errorName]
	if err == nil {
		err = errs["Default"]
	}
	return err
}
