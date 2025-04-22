package errors

import "fmt"

var (
	TokenInvalidErr         = fmt.Errorf("token is not valid")
	GuidIsDifferentErr      = fmt.Errorf("guid in tokens does not match")
	CouldNotFindSecretErr   = fmt.Errorf("could not find JWT_SECRET in environment")
	CouldNotFindRefreshHash = fmt.Errorf("could not find user with such refresh_hash")
	CouldNotFindGuid        = fmt.Errorf("user with such guid is not in our database")
	HashedRefreshDiffErr    = fmt.Errorf("bcrypt hash did not match")
	UserAlreadyLoggedIn     = fmt.Errorf("user with such guid is already logged in")
)
