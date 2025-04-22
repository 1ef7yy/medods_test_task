package errors

import "fmt"

var (
	TokenInvalidErr       = fmt.Errorf("token is not valid")
	GuidIsDifferentErr    = fmt.Errorf("guid in tokens does not match")
	CouldNotFindSecretErr = fmt.Errorf("could not find JWT_SECRET in environment")
)
