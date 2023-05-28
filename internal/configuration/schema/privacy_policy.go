package schema

import (
	"net/url"
)

// PrivacyPolicy is the privacy policy configuration.
type PrivacyPolicy struct {
	Enabled               bool     `koanf:"enabled"`
	RequireUserAcceptance bool     `koanf:"require_user_acceptance"`
	PolicyURL             *url.URL `koanf:"policy_url"`
}
