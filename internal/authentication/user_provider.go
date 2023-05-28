package authentication

import (
	"github.com/authelia/authelia/v4/internal/model"
)

// UserProvider is the interface for checking user password and
// gathering user details.
type UserProvider interface {
	model.StartupCheck

	CheckUserPassword(username string, password string) (valid bool, err error)
	GetDetails(username string) (details *UserDetails, err error)
	UpdatePassword(username string, newPassword string) (err error)
}
