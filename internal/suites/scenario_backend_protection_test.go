package suites

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

// WARNING: This scenario is intended to be used with TLS enabled in the authelia backend.

type BackendProtectionScenario struct {
	suite.Suite

	client *http.Client
}

func NewBackendProtectionScenario() *BackendProtectionScenario {
	return &BackendProtectionScenario{}
}

func (s *BackendProtectionScenario) SetupSuite() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // Needs to be enabled in suites. Not used in production.
	}

	s.client = &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func (s *BackendProtectionScenario) AssertRequestStatusCode(method, url string, expectedStatusCode int) {
	s.Run(url, func() {
		req, err := http.NewRequest(method, url, nil)
		s.Assert().NoError(err)

		res, err := s.client.Do(req)

		s.Assert().NoError(err)
		s.Assert().Equal(expectedStatusCode, res.StatusCode)
	})
}

func (s *BackendProtectionScenario) TestProtectionOfBackendEndpoints() {
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/totp", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/webauthn/assertion", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/webauthn/attestation", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/user/info/2fa_method", AutheliaBaseURL), fasthttp.StatusForbidden)

	s.AssertRequestStatusCode(fasthttp.MethodGet, fmt.Sprintf("%s/api/user/info", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodGet, fmt.Sprintf("%s/api/configuration", AutheliaBaseURL), fasthttp.StatusForbidden)

	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/totp/identity/start", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/totp/identity/finish", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/webauthn/identity/start", AutheliaBaseURL), fasthttp.StatusForbidden)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/secondfactor/webauthn/identity/finish", AutheliaBaseURL), fasthttp.StatusForbidden)
}

func (s *BackendProtectionScenario) TestInvalidEndpointsReturn404() {
	s.AssertRequestStatusCode(fasthttp.MethodGet, fmt.Sprintf("%s/api/not_existing", AutheliaBaseURL), fasthttp.StatusNotFound)
	s.AssertRequestStatusCode(fasthttp.MethodHead, fmt.Sprintf("%s/api/not_existing", AutheliaBaseURL), fasthttp.StatusNotFound)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/not_existing", AutheliaBaseURL), fasthttp.StatusNotFound)

	s.AssertRequestStatusCode(fasthttp.MethodGet, fmt.Sprintf("%s/api/not_existing/second", AutheliaBaseURL), fasthttp.StatusNotFound)
	s.AssertRequestStatusCode(fasthttp.MethodHead, fmt.Sprintf("%s/api/not_existing/second", AutheliaBaseURL), fasthttp.StatusNotFound)
	s.AssertRequestStatusCode(fasthttp.MethodPost, fmt.Sprintf("%s/api/not_existing/second", AutheliaBaseURL), fasthttp.StatusNotFound)
}

func (s *BackendProtectionScenario) TestInvalidEndpointsReturn405() {
	s.AssertRequestStatusCode(fasthttp.MethodPut, fmt.Sprintf("%s/api/configuration", AutheliaBaseURL), fasthttp.StatusMethodNotAllowed)
}

func TestRunBackendProtection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping suite test in short mode")
	}

	suite.Run(t, NewBackendProtectionScenario())
}
