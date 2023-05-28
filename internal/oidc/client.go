package oidc

import (
	"github.com/go-crypt/crypt/algorithm"
	"github.com/ory/fosite"
	"github.com/ory/x/errorsx"
	"gopkg.in/square/go-jose.v2"

	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/model"
)

// NewClient creates a new Client.
func NewClient(config schema.OpenIDConnectClient) (client Client) {
	base := &BaseClient{
		ID:               config.ID,
		Description:      config.Description,
		Secret:           config.Secret,
		SectorIdentifier: config.SectorIdentifier.String(),
		Public:           config.Public,

		EnforcePKCE:                config.EnforcePKCE || config.PKCEChallengeMethod != "",
		EnforcePKCEChallengeMethod: config.PKCEChallengeMethod != "",
		PKCEChallengeMethod:        config.PKCEChallengeMethod,

		Audience:      config.Audience,
		Scopes:        config.Scopes,
		RedirectURIs:  config.RedirectURIs,
		GrantTypes:    config.GrantTypes,
		ResponseTypes: config.ResponseTypes,
		ResponseModes: []fosite.ResponseModeType{},

		EnforcePAR: config.EnforcePAR,

		IDTokenSigningAlg:    config.IDTokenSigningAlg,
		IDTokenSigningKeyID:  config.IDTokenSigningKeyID,
		UserinfoSigningAlg:   config.UserinfoSigningAlg,
		UserinfoSigningKeyID: config.UserinfoSigningKeyID,

		Policy: authorization.NewLevel(config.Policy),

		Consent: NewClientConsent(config.ConsentMode, config.ConsentPreConfiguredDuration),
	}

	for _, mode := range config.ResponseModes {
		base.ResponseModes = append(base.ResponseModes, fosite.ResponseModeType(mode))
	}

	if config.TokenEndpointAuthMethod != "" || config.TokenEndpointAuthSigningAlg != "" ||
		len(config.PublicKeys.Values) != 0 || config.PublicKeys.URI != nil || config.RequestObjectSigningAlg != "" {
		full := &FullClient{
			BaseClient:                        base,
			TokenEndpointAuthMethod:           config.TokenEndpointAuthMethod,
			TokenEndpointAuthSigningAlgorithm: config.TokenEndpointAuthSigningAlg,
			RequestObjectSigningAlgorithm:     config.RequestObjectSigningAlg,

			JSONWebKeys: NewPublicJSONWebKeySetFromSchemaJWK(config.PublicKeys.Values),
		}

		if config.PublicKeys.URI != nil {
			full.JSONWebKeysURI = config.PublicKeys.URI.String()
		}

		client = full
	} else {
		client = base
	}

	return client
}

// GetID returns the ID.
func (c *BaseClient) GetID() string {
	return c.ID
}

// GetDescription returns the Description.
func (c *BaseClient) GetDescription() string {
	if c.Description == "" {
		c.Description = c.GetID()
	}

	return c.Description
}

// GetSecret returns the Secret.
func (c *BaseClient) GetSecret() algorithm.Digest {
	return c.Secret
}

// GetSectorIdentifier returns the SectorIdentifier for this client.
func (c *BaseClient) GetSectorIdentifier() string {
	return c.SectorIdentifier
}

// GetHashedSecret returns the Secret.
func (c *BaseClient) GetHashedSecret() (secret []byte) {
	if c.Secret == nil {
		return []byte(nil)
	}

	return []byte(c.Secret.Encode())
}

// GetRedirectURIs returns the RedirectURIs.
func (c *BaseClient) GetRedirectURIs() (redirectURIs []string) {
	return c.RedirectURIs
}

// GetGrantTypes returns the GrantTypes.
func (c *BaseClient) GetGrantTypes() fosite.Arguments {
	if len(c.GrantTypes) == 0 {
		return fosite.Arguments{"authorization_code"}
	}

	return c.GrantTypes
}

// GetResponseTypes returns the ResponseTypes.
func (c *BaseClient) GetResponseTypes() fosite.Arguments {
	if len(c.ResponseTypes) == 0 {
		return fosite.Arguments{"code"}
	}

	return c.ResponseTypes
}

// GetScopes returns the Scopes.
func (c *BaseClient) GetScopes() fosite.Arguments {
	return c.Scopes
}

// GetAudience returns the Audience.
func (c *BaseClient) GetAudience() fosite.Arguments {
	return c.Audience
}

// GetResponseModes returns the valid response modes for this client.
//
// Implements the fosite.ResponseModeClient.
func (c *BaseClient) GetResponseModes() []fosite.ResponseModeType {
	return c.ResponseModes
}

// GetIDTokenSigningAlg returns the IDTokenSigningAlg.
func (c *BaseClient) GetIDTokenSigningAlg() (alg string) {
	if c.IDTokenSigningAlg == "" {
		c.IDTokenSigningAlg = SigningAlgRSAUsingSHA256
	}

	return c.IDTokenSigningAlg
}

// GetIDTokenSigningKeyID returns the IDTokenSigningKeyID.
func (c *BaseClient) GetIDTokenSigningKeyID() (alg string) {
	return c.IDTokenSigningKeyID
}

// GetUserinfoSigningAlg returns the UserinfoSigningAlg.
func (c *BaseClient) GetUserinfoSigningAlg() string {
	if c.UserinfoSigningAlg == "" {
		c.UserinfoSigningAlg = SigningAlgNone
	}

	return c.UserinfoSigningAlg
}

// GetUserinfoSigningKeyID returns the UserinfoSigningKeyID.
func (c *BaseClient) GetUserinfoSigningKeyID() (kid string) {
	return c.UserinfoSigningKeyID
}

// GetPAREnforcement returns EnforcePAR.
func (c *BaseClient) GetPAREnforcement() bool {
	return c.EnforcePAR
}

// GetPKCEEnforcement returns EnforcePKCE.
func (c *BaseClient) GetPKCEEnforcement() bool {
	return c.EnforcePKCE
}

// GetPKCEChallengeMethodEnforcement returns EnforcePKCEChallengeMethod.
func (c *BaseClient) GetPKCEChallengeMethodEnforcement() bool {
	return c.EnforcePKCEChallengeMethod
}

// GetPKCEChallengeMethod returns PKCEChallengeMethod.
func (c *BaseClient) GetPKCEChallengeMethod() string {
	return c.PKCEChallengeMethod
}

// GetAuthorizationPolicy returns Policy.
func (c *BaseClient) GetAuthorizationPolicy() authorization.Level {
	return c.Policy
}

// GetConsentPolicy returns Consent.
func (c *BaseClient) GetConsentPolicy() ClientConsent {
	return c.Consent
}

// GetConsentResponseBody returns the proper consent response body for this session.OIDCWorkflowSession.
func (c *BaseClient) GetConsentResponseBody(consent *model.OAuth2ConsentSession) ConsentGetResponseBody {
	body := ConsentGetResponseBody{
		ClientID:          c.ID,
		ClientDescription: c.Description,
		PreConfiguration:  c.Consent.Mode == ClientConsentModePreConfigured,
	}

	if consent != nil {
		body.Scopes = consent.RequestedScopes
		body.Audience = consent.RequestedAudience
	}

	return body
}

// IsPublic returns the value of the Public property.
func (c *BaseClient) IsPublic() bool {
	return c.Public
}

// IsAuthenticationLevelSufficient returns if the provided authentication.Level is sufficient for the client of the AutheliaClient.
func (c *BaseClient) IsAuthenticationLevelSufficient(level authentication.Level) bool {
	if level == authentication.NotAuthenticated {
		return false
	}

	return authorization.IsAuthLevelSufficient(level, c.Policy)
}

// ValidatePKCEPolicy is a helper function to validate PKCE policy constraints on a per-client basis.
func (c *BaseClient) ValidatePKCEPolicy(r fosite.Requester) (err error) {
	form := r.GetRequestForm()

	if c.EnforcePKCE {
		if form.Get(FormParameterCodeChallenge) == "" {
			return errorsx.WithStack(fosite.ErrInvalidRequest.
				WithHint("Clients must include a code_challenge when performing the authorize code flow, but it is missing.").
				WithDebug("The server is configured in a way that enforces PKCE for this client."))
		}

		if c.EnforcePKCEChallengeMethod {
			if method := form.Get(FormParameterCodeChallengeMethod); method != c.PKCEChallengeMethod {
				return errorsx.WithStack(fosite.ErrInvalidRequest.
					WithHintf("Client must use code_challenge_method=%s, %s is not allowed.", c.PKCEChallengeMethod, method).
					WithDebugf("The server is configured in a way that enforces PKCE %s as challenge method for this client.", c.PKCEChallengeMethod))
			}
		}
	}

	return nil
}

// ValidatePARPolicy is a helper function to validate additional policy constraints on a per-client basis.
func (c *BaseClient) ValidatePARPolicy(r fosite.Requester, prefix string) (err error) {
	if c.EnforcePAR {
		if !IsPushedAuthorizedRequest(r, prefix) {
			switch requestURI := r.GetRequestForm().Get(FormParameterRequestURI); requestURI {
			case "":
				return errorsx.WithStack(ErrPAREnforcedClientMissingPAR.WithDebug("The request_uri parameter was empty."))
			default:
				return errorsx.WithStack(ErrPAREnforcedClientMissingPAR.WithDebugf("The request_uri parameter '%s' is malformed.", requestURI))
			}
		}
	}

	return nil
}

// ValidateResponseModePolicy is an additional check to the response mode parameter to ensure if it's omitted that the
// default response mode for the fosite.AuthorizeRequester is permitted.
func (c *BaseClient) ValidateResponseModePolicy(r fosite.AuthorizeRequester) (err error) {
	if r.GetResponseMode() != fosite.ResponseModeDefault {
		return nil
	}

	m := r.GetDefaultResponseMode()

	modes := c.GetResponseModes()

	if len(modes) == 0 {
		return nil
	}

	for _, mode := range modes {
		if m == mode {
			return nil
		}
	}

	return errorsx.WithStack(fosite.ErrUnsupportedResponseMode.WithHintf(`The request omitted the response_mode making the default response_mode "%s" based on the other authorization request parameters but registered OAuth 2.0 client doesn't support this response_mode`, m))
}

// GetRequestURIs is an array of request_uri values that are pre-registered by the RP for use at the OP. Servers MAY
// cache the contents of the files referenced by these URIs and not retrieve them at the time they are used in a request.
// OPs can require that request_uri values used be pre-registered with the require_request_uri_registration
// discovery parameter.
func (c *FullClient) GetRequestURIs() []string {
	return c.RequestURIs
}

// GetJSONWebKeys returns the JSON Web Key Set containing the public key used by the client to authenticate.
func (c *FullClient) GetJSONWebKeys() *jose.JSONWebKeySet {
	return c.JSONWebKeys
}

// GetJSONWebKeysURI returns the URL for lookup of JSON Web Key Set containing the
// public key used by the client to authenticate.
func (c *FullClient) GetJSONWebKeysURI() string {
	return c.JSONWebKeysURI
}

// GetRequestObjectSigningAlgorithm returns the JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request
// Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm.
func (c *FullClient) GetRequestObjectSigningAlgorithm() string {
	return c.RequestObjectSigningAlgorithm
}

// GetTokenEndpointAuthMethod returns the requested Client Authentication Method for the Token Endpoint. The options are
// client_secret_post, client_secret_basic, client_secret_jwt, private_key_jwt, and none.
func (c *FullClient) GetTokenEndpointAuthMethod() string {
	if c.TokenEndpointAuthMethod == "" {
		if c.Public {
			c.TokenEndpointAuthMethod = ClientAuthMethodNone
		} else {
			c.TokenEndpointAuthMethod = ClientAuthMethodClientSecretBasic
		}
	}

	return c.TokenEndpointAuthMethod
}

// GetTokenEndpointAuthSigningAlgorithm returns the JWS [JWS] alg algorithm [JWA] that MUST be used for signing the JWT
// [JWT] used to authenticate the Client at the Token Endpoint for the private_key_jwt and client_secret_jwt
// authentication methods.
func (c *FullClient) GetTokenEndpointAuthSigningAlgorithm() string {
	if c.TokenEndpointAuthSigningAlgorithm == "" {
		c.TokenEndpointAuthSigningAlgorithm = SigningAlgRSAUsingSHA256
	}

	return c.TokenEndpointAuthSigningAlgorithm
}
