package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/ory/fosite/storage"

	"github.com/authelia/authelia/v4/internal/model"
)

// Provider is an interface providing storage capabilities for persisting any kind of data related to Authelia.
type Provider interface {
	model.StartupCheck

	RegulatorProvider

	storage.Transactional

	SavePreferred2FAMethod(ctx context.Context, username string, method string) (err error)
	LoadPreferred2FAMethod(ctx context.Context, username string) (method string, err error)
	LoadUserInfo(ctx context.Context, username string) (info model.UserInfo, err error)

	SaveUserOpaqueIdentifier(ctx context.Context, subject model.UserOpaqueIdentifier) (err error)
	LoadUserOpaqueIdentifier(ctx context.Context, identifier uuid.UUID) (subject *model.UserOpaqueIdentifier, err error)
	LoadUserOpaqueIdentifiers(ctx context.Context) (identifiers []model.UserOpaqueIdentifier, err error)
	LoadUserOpaqueIdentifierBySignature(ctx context.Context, service, sectorID, username string) (subject *model.UserOpaqueIdentifier, err error)

	SaveIdentityVerification(ctx context.Context, verification model.IdentityVerification) (err error)
	ConsumeIdentityVerification(ctx context.Context, jti string, ip model.NullIP) (err error)
	FindIdentityVerification(ctx context.Context, jti string) (found bool, err error)

	SaveTOTPConfiguration(ctx context.Context, config model.TOTPConfiguration) (err error)
	UpdateTOTPConfigurationSignIn(ctx context.Context, id int, lastUsedAt sql.NullTime) (err error)
	DeleteTOTPConfiguration(ctx context.Context, username string) (err error)
	LoadTOTPConfiguration(ctx context.Context, username string) (config *model.TOTPConfiguration, err error)
	LoadTOTPConfigurations(ctx context.Context, limit, page int) (configs []model.TOTPConfiguration, err error)

	SaveWebAuthnDevice(ctx context.Context, device model.WebAuthnDevice) (err error)
	UpdateWebAuthnDeviceSignIn(ctx context.Context, id int, rpid string, lastUsedAt sql.NullTime, signCount uint32, cloneWarning bool) (err error)
	DeleteWebAuthnDevice(ctx context.Context, kid string) (err error)
	DeleteWebAuthnDeviceByUsername(ctx context.Context, username, description string) (err error)
	LoadWebAuthnDevices(ctx context.Context, limit, page int) (devices []model.WebAuthnDevice, err error)
	LoadWebAuthnDevicesByUsername(ctx context.Context, username string) (devices []model.WebAuthnDevice, err error)

	SavePreferredDuoDevice(ctx context.Context, device model.DuoDevice) (err error)
	DeletePreferredDuoDevice(ctx context.Context, username string) (err error)
	LoadPreferredDuoDevice(ctx context.Context, username string) (device *model.DuoDevice, err error)

	SaveOAuth2ConsentPreConfiguration(ctx context.Context, config model.OAuth2ConsentPreConfig) (insertedID int64, err error)
	LoadOAuth2ConsentPreConfigurations(ctx context.Context, clientID string, subject uuid.UUID) (rows *ConsentPreConfigRows, err error)

	SaveOAuth2ConsentSession(ctx context.Context, consent model.OAuth2ConsentSession) (err error)
	SaveOAuth2ConsentSessionSubject(ctx context.Context, consent model.OAuth2ConsentSession) (err error)
	SaveOAuth2ConsentSessionResponse(ctx context.Context, consent model.OAuth2ConsentSession, rejection bool) (err error)
	SaveOAuth2ConsentSessionGranted(ctx context.Context, id int) (err error)
	LoadOAuth2ConsentSessionByChallengeID(ctx context.Context, challengeID uuid.UUID) (consent *model.OAuth2ConsentSession, err error)

	SaveOAuth2Session(ctx context.Context, sessionType OAuth2SessionType, session model.OAuth2Session) (err error)
	RevokeOAuth2Session(ctx context.Context, sessionType OAuth2SessionType, signature string) (err error)
	RevokeOAuth2SessionByRequestID(ctx context.Context, sessionType OAuth2SessionType, requestID string) (err error)
	DeactivateOAuth2Session(ctx context.Context, sessionType OAuth2SessionType, signature string) (err error)
	DeactivateOAuth2SessionByRequestID(ctx context.Context, sessionType OAuth2SessionType, requestID string) (err error)
	LoadOAuth2Session(ctx context.Context, sessionType OAuth2SessionType, signature string) (session *model.OAuth2Session, err error)

	SaveOAuth2PARContext(ctx context.Context, par model.OAuth2PARContext) (err error)
	LoadOAuth2PARContext(ctx context.Context, signature string) (par *model.OAuth2PARContext, err error)
	RevokeOAuth2PARContext(ctx context.Context, signature string) (err error)

	SaveOAuth2BlacklistedJTI(ctx context.Context, blacklistedJTI model.OAuth2BlacklistedJTI) (err error)
	LoadOAuth2BlacklistedJTI(ctx context.Context, signature string) (blacklistedJTI *model.OAuth2BlacklistedJTI, err error)

	SchemaTables(ctx context.Context) (tables []string, err error)
	SchemaVersion(ctx context.Context) (version int, err error)
	SchemaLatestVersion() (version int, err error)

	SchemaMigrate(ctx context.Context, up bool, version int) (err error)
	SchemaMigrationHistory(ctx context.Context) (migrations []model.Migration, err error)
	SchemaMigrationsUp(ctx context.Context, version int) (migrations []model.SchemaMigration, err error)
	SchemaMigrationsDown(ctx context.Context, version int) (migrations []model.SchemaMigration, err error)

	SchemaEncryptionChangeKey(ctx context.Context, key string) (err error)
	SchemaEncryptionCheckKey(ctx context.Context, verbose bool) (result EncryptionValidationResult, err error)

	Close() (err error)
}

// RegulatorProvider is an interface providing storage capabilities for persisting any kind of data related to the regulator.
type RegulatorProvider interface {
	AppendAuthenticationLog(ctx context.Context, attempt model.AuthenticationAttempt) (err error)
	LoadAuthenticationLogs(ctx context.Context, username string, fromDate time.Time, limit, page int) (attempts []model.AuthenticationAttempt, err error)
}
