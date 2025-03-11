package errs

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Constant Errors
var (
	ErrNotEnoughBalance = errors.New("not enough balance")
	ErrNoUser           = errors.New("the user not found")
	ErrNoItem           = errors.New("the item does not exist")
	ErrInvalidReq       = errors.New("invalid request body")
)

func ValidationError(errs error) string {
	s := "Folloing Input fields need to be changed:\n"
	for _, err := range errs.(validator.ValidationErrors) {

		i := fmt.Sprintf("Field : %s does not satisfy the requirement : [%s]\n", err.Field(), err.Tag())
		s = s + i

	}
	return s
}
