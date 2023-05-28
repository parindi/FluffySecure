package suites

import (
	"testing"

	"github.com/go-rod/rod"
)

func (rs *RodSession) doRegisterThenLogout(t *testing.T, page *rod.Page, username, password string) string {
	secret := rs.doLoginAndRegisterTOTP(t, page, username, password, false)
	rs.doLogout(t, page)

	return secret
}
