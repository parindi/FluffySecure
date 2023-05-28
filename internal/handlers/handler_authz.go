package handlers

import (
	"fmt"
	"net/url"

	"github.com/authelia/authelia/v4/internal/authentication"
	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/session"
	"github.com/authelia/authelia/v4/internal/utils"
)

// Handler is the middlewares.RequestHandler for Authz.
func (authz *Authz) Handler(ctx *middlewares.AutheliaCtx) {
	var (
		object      authorization.Object
		autheliaURL *url.URL
		provider    *session.Session
		err         error
	)

	if object, err = authz.handleGetObject(ctx); err != nil {
		ctx.Logger.WithError(err).Error("Error getting Target URL and Request Method")

		ctx.ReplyStatusCode(authz.config.StatusCodeBadRequest)

		return
	}

	if !utils.IsURISecure(object.URL) {
		ctx.Logger.Errorf("Target URL '%s' has an insecure scheme '%s', only the 'https' and 'wss' schemes are supported so session cookies can be transmitted securely", object.URL.String(), object.URL.Scheme)

		ctx.ReplyStatusCode(authz.config.StatusCodeBadRequest)

		return
	}

	if provider, err = ctx.GetSessionProviderByTargetURL(object.URL); err != nil {
		ctx.Logger.WithError(err).WithField("target_url", object.URL.String()).Error("Target URL does not appear to have a relevant session cookies configuration")

		ctx.ReplyStatusCode(authz.config.StatusCodeBadRequest)

		return
	}

	if autheliaURL, err = authz.getAutheliaURL(ctx, provider); err != nil {
		ctx.Logger.WithError(err).WithField("target_url", object.URL.String()).Error("Error occurred trying to determine the external Authelia URL for Target URL")

		ctx.ReplyStatusCode(authz.config.StatusCodeBadRequest)

		return
	}

	var (
		authn    Authn
		strategy AuthnStrategy
	)

	if authn, strategy, err = authz.authn(ctx, provider); err != nil {
		authn.Object = object

		ctx.Logger.WithError(err).Error("Error occurred while attempting to authenticate a request")

		switch strategy {
		case nil:
			ctx.ReplyUnauthorized()
		default:
			strategy.HandleUnauthorized(ctx, &authn, authz.getRedirectionURL(&object, autheliaURL))
		}

		return
	}

	authn.Object = object
	authn.Method = friendlyMethod(authn.Object.Method)

	ruleHasSubject, required := ctx.Providers.Authorizer.GetRequiredLevel(
		authorization.Subject{
			Username: authn.Details.Username,
			Groups:   authn.Details.Groups,
			IP:       ctx.RemoteIP(),
		},
		object,
	)

	switch isAuthzResult(authn.Level, required, ruleHasSubject) {
	case AuthzResultForbidden:
		ctx.Logger.Infof("Access to '%s' is forbidden to user '%s'", object.URL.String(), authn.Username)
		ctx.ReplyForbidden()
	case AuthzResultUnauthorized:
		var handler HandlerAuthzUnauthorized

		if strategy != nil {
			handler = strategy.HandleUnauthorized
		} else {
			handler = authz.handleUnauthorized
		}

		handler(ctx, &authn, authz.getRedirectionURL(&object, autheliaURL))
	case AuthzResultAuthorized:
		authz.handleAuthorized(ctx, &authn)
	}
}

func (authz *Authz) getAutheliaURL(ctx *middlewares.AutheliaCtx, provider *session.Session) (autheliaURL *url.URL, err error) {
	if autheliaURL, err = authz.handleGetAutheliaURL(ctx); err != nil {
		return nil, err
	}

	switch {
	case authz.implementation == AuthzImplLegacy:
		return autheliaURL, nil
	case autheliaURL != nil:
		switch {
		case utils.HasURIDomainSuffix(autheliaURL, provider.Config.Domain):
			return autheliaURL, nil
		default:
			return nil, fmt.Errorf("authelia url '%s' is not valid for detected domain '%s' as the url does not have the domain as a suffix", autheliaURL.String(), provider.Config.Domain)
		}
	}

	if provider.Config.AutheliaURL != nil {
		return provider.Config.AutheliaURL, nil
	}

	return nil, fmt.Errorf("authelia url lookup failed")
}

func (authz *Authz) getRedirectionURL(object *authorization.Object, autheliaURL *url.URL) (redirectionURL *url.URL) {
	if autheliaURL == nil {
		return nil
	}

	redirectionURL, _ = url.ParseRequestURI(autheliaURL.String())

	if redirectionURL.Path == "" {
		redirectionURL.Path = "/"
	}

	qry := redirectionURL.Query()

	qry.Set(queryArgRD, object.URL.String())

	if object.Method != "" {
		qry.Set(queryArgRM, object.Method)
	}

	redirectionURL.RawQuery = qry.Encode()

	return redirectionURL
}

func (authz *Authz) authn(ctx *middlewares.AutheliaCtx, provider *session.Session) (authn Authn, strategy AuthnStrategy, err error) {
	for _, strategy = range authz.strategies {
		if authn, err = strategy.Get(ctx, provider); err != nil {
			if strategy.CanHandleUnauthorized() {
				return Authn{Type: authn.Type, Level: authentication.NotAuthenticated, Username: anonymous}, strategy, err
			}

			return Authn{Type: authn.Type, Level: authentication.NotAuthenticated, Username: anonymous}, nil, err
		}

		if authn.Level != authentication.NotAuthenticated {
			break
		}
	}

	if strategy.CanHandleUnauthorized() {
		return authn, strategy, err
	}

	return authn, nil, nil
}
