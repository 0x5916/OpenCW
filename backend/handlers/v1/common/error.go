package common

import "errors"

var ErrExpiredToken = errors.New("expired")
var ErrInvalidToken = errors.New("invalid")
