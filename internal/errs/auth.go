package errs

import "fmt"

// todo universal error handling
var WrongPassword = fmt.Errorf("wrong password")
var AccessTokenExpired = fmt.Errorf("access token is expired")
var InvalidRefreshToken = fmt.Errorf("invalid refresh token")
var InvalidAccessToken = fmt.Errorf("invalid access token")
var RefreshTokenExpired = fmt.Errorf("refresh token is expired")
var RefreshTokenNotFound = fmt.Errorf("refresh token not found")
