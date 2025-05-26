package erorrs

import "errors"

var ErrUsernameAlreadyExists = errors.New("username already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrTokenNotFound = errors.New("token not found")
var ErrExerciseAlreadyExists = errors.New("exercise already exists")

const ErrorKey = "error"
