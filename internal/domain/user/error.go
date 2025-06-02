package user

type LoginError struct {
	Reason string
}

func (e *LoginError) Error() string {
	return e.Reason
}

var (
	ErrUserNotFound    = &LoginError{"user not found"}
	ErrInvalidPassword = &LoginError{"invalid password"}
)
