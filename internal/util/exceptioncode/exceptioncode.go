package exceptioncode

import "errors"

var (
	ErrEmptyResult        = errors.New("empty result")
	ErrTokenInvalid       = errors.New(CodeTokenInvalid)
	ErrTokenExpired       = errors.New(CodeTokenExpired)
	ErrInvalidRequest     = errors.New("invalid request")
	ErrFailedReadEnum     = errors.New("failed to read enum")
	ErrUnableToLock       = errors.New("unable to lock")
	ErrMethodNotAllowed   = errors.New("method not allowed")
	ErrUnauthorizedMethod = errors.New("unauthorized method")

	// spesific mysql error
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrDupeKey             = errors.New("duplicate key value")
)

const (
	CodeAccountBlocked       = "ACCOUNT_BLOCKED"
	CodeAccountDeleted       = "ACCOUNT_DELETED"
	CodeAccountLocked        = "ACCOUNT_LOCKED"
	CodeBadRequest           = "BAD_REQUEST"
	CodeConflict             = "CONFLICT"
	CodeDataAlreadyExist     = "DATA_ALREADY_EXIST"
	CodeDataLocked           = "DATA_LOCKED"
	CodeDataNotFound         = "DATA_NOT_FOUND"
	CodeDateExpired          = "DATE_EXPIRED"
	CodeEmailExist           = "EMAIL_EXIST"
	CodeEmailNotVerified     = "EMAIL_NOT_VERIFIED"
	CodeFileNotProcessedYet  = "FILE_NOT_PROCESSED_YET"
	CodeFileTypeNotSupported = "FILE_TYPE_NOT_SUPPORTED"
	CodeForbidden            = "FORBIDDEN"
	CodeInternalServerError  = "INTERNAL_SERVER_ERROR"
	CodeInvalidAge           = "INVALID_AGE"
	CodeInvalidCredential    = "INVALID_CREDENTIAL"
	CodeInvalidData          = "INVALID_DATA"
	CodeInvalidFormat        = "INVALID_FORMAT"
	CodeInvalidRequest       = "INVALID_REQUEST"
	CodeInvalidValidation    = "INVALID_VALIDATION"
	CodeMethodNotAllowed     = "METHOD_NOT_ALLOWED"
	CodeMissingRequiredData  = "MISSING_REQUIRED_DATA"
	CodeOtpFailed            = "OTP_FAILED"
	CodeOtpInvalid           = "OTP_INVALID"
	CodeQuotaLimitReached    = "QUOTA_LIMIT_REACHED"
	CodeRequestFailed        = "REQUEST_FAILED"
	CodeRequestTooFast       = "REQUEST_TOO_FAST"
	CodeServiceUnavailable   = "SERVICE_UNAVAILABLE"
	CodeTokenExpired         = "TOKEN_EXPIRED"
	CodeTokenInvalid         = "TOKEN_INVALID"
	CodeUnauthorized         = "UNAUTHORIZED"
)

type (
	errorType struct {
		ErrorMessage string
	}
	ErrorNotFound            errorType
	ErrorOTPFailed           errorType
	ErrorForeignKeyViolation errorType
)

type NotFoundError struct {
	Error string
}

func NewNotFoundError(error string) NotFoundError {
	return NotFoundError{Error: error}
}
