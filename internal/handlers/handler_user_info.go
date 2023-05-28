package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/model"
	"github.com/authelia/authelia/v4/internal/session"
	"github.com/authelia/authelia/v4/internal/utils"
)

// UserInfoPOST handles setting up info for users if necessary when they login.
func UserInfoPOST(ctx *middlewares.AutheliaCtx) {
	var (
		userSession session.UserSession
		userInfo    model.UserInfo
		err         error
	)

	if userSession, err = ctx.GetSession(); err != nil {
		ctx.Logger.WithError(err).Error("Error occurred retrieving user session")

		ctx.ReplyForbidden()

		return
	}

	if _, err = ctx.Providers.StorageProvider.LoadPreferred2FAMethod(ctx, userSession.Username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if err = ctx.Providers.StorageProvider.SavePreferred2FAMethod(ctx, userSession.Username, ""); err != nil {
				ctx.Error(fmt.Errorf("unable to load user information: %v", err), messageOperationFailed)
			}
		} else {
			ctx.Error(fmt.Errorf("unable to load user information: %v", err), messageOperationFailed)
		}
	}

	if userInfo, err = ctx.Providers.StorageProvider.LoadUserInfo(ctx, userSession.Username); err != nil {
		ctx.Error(fmt.Errorf("unable to load user information: %v", err), messageOperationFailed)
		return
	}

	var (
		changed bool
	)

	if changed = userInfo.SetDefaultPreferred2FAMethod(ctx.AvailableSecondFactorMethods(), ctx.Configuration.Default2FAMethod); changed {
		if err = ctx.Providers.StorageProvider.SavePreferred2FAMethod(ctx, userSession.Username, userInfo.Method); err != nil {
			ctx.Error(fmt.Errorf("unable to save user two factor method: %v", err), messageOperationFailed)
			return
		}
	}

	userInfo.DisplayName = userSession.DisplayName

	err = ctx.SetJSONBody(userInfo)
	if err != nil {
		ctx.Logger.Errorf("Unable to set user info response in body: %s", err)
	}
}

// UserInfoGET get the info related to the user identified by the session.
func UserInfoGET(ctx *middlewares.AutheliaCtx) {
	var (
		userSession session.UserSession
		err         error
	)

	if userSession, err = ctx.GetSession(); err != nil {
		ctx.Logger.WithError(err).Error("Error occurred retrieving user session")

		ctx.ReplyForbidden()

		return
	}

	userInfo, err := ctx.Providers.StorageProvider.LoadUserInfo(ctx, userSession.Username)
	if err != nil {
		ctx.Error(fmt.Errorf("unable to load user information: %v", err), messageOperationFailed)
		return
	}

	userInfo.DisplayName = userSession.DisplayName

	err = ctx.SetJSONBody(userInfo)
	if err != nil {
		ctx.Logger.Errorf("Unable to set user info response in body: %s", err)
	}
}

// MethodPreferencePOST update the user preferences regarding 2FA method.
func MethodPreferencePOST(ctx *middlewares.AutheliaCtx) {
	var (
		bodyJSON bodyPreferred2FAMethod

		userSession session.UserSession
		err         error
	)

	if userSession, err = ctx.GetSession(); err != nil {
		ctx.Logger.WithError(err).Error("Error occurred retrieving user session")

		ctx.Error(err, messageOperationFailed)

		return
	}

	if err = ctx.ParseBody(&bodyJSON); err != nil {
		ctx.Error(err, messageOperationFailed)
		return
	}

	if !utils.IsStringInSlice(bodyJSON.Method, ctx.AvailableSecondFactorMethods()) {
		ctx.Error(fmt.Errorf("unknown or unavailable method '%s', it should be one of %s", bodyJSON.Method, strings.Join(ctx.AvailableSecondFactorMethods(), ", ")), messageOperationFailed)
		return
	}

	ctx.Logger.Debugf("Save new preferred 2FA method of user %s to %s", userSession.Username, bodyJSON.Method)
	err = ctx.Providers.StorageProvider.SavePreferred2FAMethod(ctx, userSession.Username, bodyJSON.Method)

	if err != nil {
		ctx.Error(fmt.Errorf("unable to save new preferred 2FA method: %s", err), messageOperationFailed)
		return
	}

	ctx.ReplyOK()
}
