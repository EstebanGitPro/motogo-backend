package person

import "errors"

var (
	ErrUserCannotSave                  = errors.New("error user can not save")
	ErrGetUsers                        = errors.New("error get users")
	ErrDuplicateUser                   = errors.New("user already exists")
	ErrSavingUser                      = errors.New("error saving user")
	ErrUserCannotGet                   = errors.New("error can not get user")
	ErrUserCannotFound                 = errors.New("error can no found user")
	ErrGettingUserByEmail              = errors.New("error getting user by the email")
	ErrNotFoundUserByEmail             = errors.New("error not found user by email")
	ErrUserCannotLogin                 = errors.New("error user can not login")
	ErrValidationUser                  = errors.New("error validation user")
	ErrInvalidJson                     = errors.New("error invalid json")
	ErrGerateToken                     = errors.New("error generate token")
	ErrUserCannotUpdateEmailVerified   = errors.New("error user can not update email verified")
	ErrorEmailNotVerified              = errors.New("email not verified yet")
	ErrTokenExpired                    = errors.New("verification token has expired")
	ErrTokenAlreadyUsed                = errors.New("verification token has already been used")
	ErrTokenNotFound                   = errors.New("verification token not found")
	ErrCleanupExpiredTokens            = errors.New("error cleaning up expired tokens")
	ErrGetVerificationToken            = errors.New("error getting verification token")
	ErrVerificationTokenNotFound       = errors.New("verification token not found")
	ErrUserCannotSaveVerificationToken = errors.New("error user can not save verification token")
	ErrUserCannotUpdate                = errors.New("error user can not update")
)
