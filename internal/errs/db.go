package errs

import (
	"fmt"
)

var UserAlreadyExists = fmt.Errorf("user already exists")
var UserNotExists = fmt.Errorf("user does not exist")
