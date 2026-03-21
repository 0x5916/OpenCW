package common

import "errors"

var ErrExpiredToken = errors.New("expired")
var ErrInvalidToken = errors.New("invalid")

const (
	ErrorCodeInvalidRequestBody        = "INVALID_REQUEST_BODY"
	ErrorCodeInternalServerError       = "INTERNAL_SERVER_ERROR"
	ErrorCodeDatabaseFailure           = "DATABASE_FAILURE"
	ErrorCodeInvalidCredentials        = "INVALID_CREDENTIALS"
	ErrorCodeConflict                  = "CONFLICT"
	ErrorCodeInvalidToken              = "INVALID_TOKEN"
	ErrorCodeExpiredToken              = "EXPIRED_TOKEN"
	ErrorCodeAuthHeaderRequired        = "AUTH_HEADER_REQUIRED"
	ErrorCodeInvalidAuthHeaderFormat   = "INVALID_AUTH_HEADER_FORMAT"
	ErrorCodeUserNotFound              = "USER_NOT_FOUND"
	ErrorCodeSettingsFetchFailed       = "SETTINGS_FETCH_FAILED"
	ErrorCodeSettingsUpdateFailed      = "SETTINGS_UPDATE_FAILED"
	ErrorCodeProgressQueryFailed       = "PROGRESS_QUERY_FAILED"
	ErrorCodeProgressCreateFailed      = "PROGRESS_CREATE_FAILED"
	ErrorCodePasswordHashFailed        = "PASSWORD_HASH_FAILED"
	ErrorCodeTokenIssueFailed          = "TOKEN_ISSUE_FAILED"
	ErrorCodeEmailAlreadyInUse         = "EMAIL_ALREADY_IN_USE"
	ErrorCodeUsernameAlreadyInUse      = "USERNAME_ALREADY_IN_USE"
	ErrorCodeUsernameEmailAlreadyInUse = "USERNAME_EMAIL_ALREADY_IN_USE"
	ErrorCodeEmailUnchanged            = "EMAIL_UNCHANGED"
	ErrorCodeCallSignAlreadyInUse      = "CALL_SIGN_ALREADY_IN_USE"
)

func NewErrorResponse(code string, message string) ErrorResponse {
	return ErrorResponse{Code: code, Error: message}
}
