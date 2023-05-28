package validator

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
)

func newDefaultConfig() schema.Configuration {
	config := schema.Configuration{}
	config.Server.Address = &schema.AddressTCP{Address: schema.NewAddressFromNetworkValues("tcp", loopback, 9090)}
	config.Log.Level = "info"
	config.Log.Format = "text"
	config.JWTSecret = testJWTSecret
	config.AuthenticationBackend.File = &schema.FileAuthenticationBackend{
		Path: "/a/path",
	}
	config.AccessControl = schema.AccessControlConfiguration{
		DefaultPolicy: "two_factor",
	}
	config.Session = schema.SessionConfiguration{
		Secret: "secret",
		Cookies: []schema.SessionCookieConfiguration{
			{
				SessionCookieCommonConfiguration: schema.SessionCookieCommonConfiguration{
					Name:   "authelia_session",
					Domain: exampleDotCom,
				},
			},
		},
	}
	config.Storage.EncryptionKey = testEncryptionKey
	config.Storage.Local = &schema.LocalStorageConfiguration{
		Path: "abc",
	}
	config.Notifier = schema.NotifierConfiguration{
		FileSystem: &schema.FileSystemNotifierConfiguration{
			Filename: "/tmp/file",
		},
	}

	return config
}

func TestShouldEnsureNotifierConfigIsProvided(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()

	ValidateConfiguration(&config, validator)
	require.Len(t, validator.Errors(), 0)

	config = newDefaultConfig()

	config.Notifier.SMTP = nil
	config.Notifier.FileSystem = nil

	ValidateConfiguration(&config, validator)
	require.Len(t, validator.Errors(), 1)
	assert.EqualError(t, validator.Errors()[0], "notifier: you must ensure either the 'smtp' or 'filesystem' notifier is configured")
}

func TestShouldAddDefaultAccessControl(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()

	config.AccessControl.DefaultPolicy = ""
	config.AccessControl.Rules = []schema.ACLRule{
		{
			Policy: "bypass",
			Domains: []string{
				"public.example.com",
			},
		},
	}

	ValidateConfiguration(&config, validator)
	require.Len(t, validator.Errors(), 0)
	assert.NotNil(t, config.AccessControl)
	assert.Equal(t, "deny", config.AccessControl.DefaultPolicy)
}

func TestShouldRaiseErrorWithUndefinedJWTSecretKey(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()
	config.JWTSecret = ""

	ValidateConfiguration(&config, validator)
	require.Len(t, validator.Errors(), 1)
	require.Len(t, validator.Warnings(), 1)

	assert.EqualError(t, validator.Errors()[0], "option 'jwt_secret' is required")
	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func TestShouldRaiseErrorWithBadDefaultRedirectionURL(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()
	config.DefaultRedirectionURL = "bad_default_redirection_url"

	ValidateConfiguration(&config, validator)
	require.Len(t, validator.Errors(), 1)
	require.Len(t, validator.Warnings(), 1)

	assert.EqualError(t, validator.Errors()[0], "option 'default_redirection_url' is invalid: could not parse 'bad_default_redirection_url' as a URL")
	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func TestShouldNotOverrideCertificatesDirectoryAndShouldPassWhenBlank(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()

	ValidateConfiguration(&config, validator)

	assert.Len(t, validator.Errors(), 0)
	require.Len(t, validator.Warnings(), 1)

	require.Equal(t, "", config.CertificatesDirectory)

	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func TestShouldRaiseErrorOnInvalidCertificatesDirectory(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()
	config.CertificatesDirectory = "not-a-real-file.go"

	ValidateConfiguration(&config, validator)

	require.Len(t, validator.Errors(), 1)
	require.Len(t, validator.Warnings(), 1)

	if runtime.GOOS == "windows" {
		assert.EqualError(t, validator.Errors()[0], "the location 'certificates_directory' could not be inspected: CreateFile not-a-real-file.go: The system cannot find the file specified.")
	} else {
		assert.EqualError(t, validator.Errors()[0], "the location 'certificates_directory' could not be inspected: stat not-a-real-file.go: no such file or directory")
	}

	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")

	config = newDefaultConfig()

	validator = schema.NewStructValidator()
	config.CertificatesDirectory = "const.go"

	ValidateConfiguration(&config, validator)

	require.Len(t, validator.Errors(), 1)
	require.Len(t, validator.Warnings(), 1)

	assert.EqualError(t, validator.Errors()[0], "the location 'certificates_directory' refers to 'const.go' is not a directory")
	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func TestShouldNotRaiseErrorOnValidCertificatesDirectory(t *testing.T) {
	validator := schema.NewStructValidator()
	config := newDefaultConfig()
	config.CertificatesDirectory = "../../suites/common/pki"

	ValidateConfiguration(&config, validator)

	assert.Len(t, validator.Errors(), 0)
	require.Len(t, validator.Warnings(), 1)

	assert.EqualError(t, validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func TestValidateDefault2FAMethod(t *testing.T) {
	testCases := []struct {
		desc         string
		have         *schema.Configuration
		expectedErrs []string
	}{
		{
			desc: "ShouldAllowConfiguredMethodTOTP",
			have: &schema.Configuration{
				Default2FAMethod: "totp",
				DuoAPI: schema.DuoAPIConfiguration{
					SecretKey:      "a key",
					IntegrationKey: "another key",
					Hostname:       "none",
				},
			},
		},
		{
			desc: "ShouldAllowConfiguredMethodWebAuthn",
			have: &schema.Configuration{
				Default2FAMethod: "webauthn",
				DuoAPI: schema.DuoAPIConfiguration{
					SecretKey:      "a key",
					IntegrationKey: "another key",
					Hostname:       "none",
				},
			},
		},
		{
			desc: "ShouldAllowConfiguredMethodMobilePush",
			have: &schema.Configuration{
				Default2FAMethod: "mobile_push",
				DuoAPI: schema.DuoAPIConfiguration{
					SecretKey:      "a key",
					IntegrationKey: "another key",
					Hostname:       "none",
				},
			},
		},
		{
			desc: "ShouldNotAllowDisabledMethodTOTP",
			have: &schema.Configuration{
				Default2FAMethod: "totp",
				DuoAPI: schema.DuoAPIConfiguration{
					SecretKey:      "a key",
					IntegrationKey: "another key",
					Hostname:       "none",
				},
				TOTP: schema.TOTPConfiguration{Disable: true},
			},
			expectedErrs: []string{
				"option 'default_2fa_method' must be one of the enabled options 'webauthn' or 'mobile_push' but it's configured as 'totp'",
			},
		},
		{
			desc: "ShouldNotAllowDisabledMethodWebAuthn",
			have: &schema.Configuration{
				Default2FAMethod: "webauthn",
				DuoAPI: schema.DuoAPIConfiguration{
					SecretKey:      "a key",
					IntegrationKey: "another key",
					Hostname:       "none",
				},
				WebAuthn: schema.WebAuthnConfiguration{Disable: true},
			},
			expectedErrs: []string{
				"option 'default_2fa_method' must be one of the enabled options 'totp' or 'mobile_push' but it's configured as 'webauthn'",
			},
		},
		{
			desc: "ShouldNotAllowDisabledMethodMobilePush",
			have: &schema.Configuration{
				Default2FAMethod: "mobile_push",
				DuoAPI:           schema.DuoAPIConfiguration{Disable: true},
			},
			expectedErrs: []string{
				"option 'default_2fa_method' must be one of the enabled options 'totp' or 'webauthn' but it's configured as 'mobile_push'",
			},
		},
		{
			desc: "ShouldNotAllowInvalidMethodDuo",
			have: &schema.Configuration{
				Default2FAMethod: "duo",
			},
			expectedErrs: []string{
				"option 'default_2fa_method' must be one of 'totp', 'webauthn', or 'mobile_push' but it's configured as 'duo'",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			validator := schema.NewStructValidator()

			validateDefault2FAMethod(tc.have, validator)

			assert.Len(t, validator.Warnings(), 0)

			errs := validator.Errors()

			require.Len(t, errs, len(tc.expectedErrs))

			for i, expected := range tc.expectedErrs {
				t.Run(fmt.Sprintf("Err%d", i+1), func(t *testing.T) {
					assert.EqualError(t, errs[i], expected)
				})
			}
		})
	}
}
