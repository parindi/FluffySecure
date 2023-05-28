package handlers

import (
	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/session"
)

// StateGET is the handler serving the user state.
func StateGET(ctx *middlewares.AutheliaCtx) {
	var (
		userSession session.UserSession
		err         error
	)

	if userSession, err = ctx.GetSession(); err != nil {
		ctx.Logger.WithError(err).Error("Error occurred retrieving user session")

		ctx.ReplyForbidden()

		return
	}

	stateResponse := StateResponse{
		Username:              userSession.Username,
		AuthenticationLevel:   userSession.AuthenticationLevel,
		DefaultRedirectionURL: ctx.Configuration.DefaultRedirectionURL,
	}

	if err = ctx.SetJSONBody(stateResponse); err != nil {
		ctx.Logger.Errorf("Unable to set state response in body: %s", err)
	}
}
