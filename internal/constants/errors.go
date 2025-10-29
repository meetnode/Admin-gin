package constants

var (
	// General errors
	ErrSomethingWentWrong = "something went wrong"
	ErrInternalServer     = "internal server error"
	ErrInvalidRequest     = "invalid request"
	ErrUnauthorized       = "unauthorized access"
	ErrForbidden          = "forbidden access"
	ErrNotFound           = "record not found"
	ErrEmailSendFailed    = "failed to send email"
	ErrDatabase           = "database error"
	ErrInvalidEmail       = "email is not valid"

	// User-related errors
	ErrUserNotFound        = "user not found"
	ErrInvalidUserId       = "Invalid user ID"
	ErrInvalidPermissionId = "Invalid permission ID"
	ErrInvalidPassword     = "invalid password"
	ErrEmailExists         = "email already exists"
	ErrTokenExpired        = "token expired or invalid"

	// File-related errors
	ErrFileTooLarge     = "file size too large"
	ErrFileUploadFailed = "file upload failed"
)
