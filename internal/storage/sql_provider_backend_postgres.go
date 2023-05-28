package storage

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/utils"
)

// PostgreSQLProvider is a PostgreSQL provider.
type PostgreSQLProvider struct {
	SQLProvider
}

// NewPostgreSQLProvider a PostgreSQL provider.
func NewPostgreSQLProvider(config *schema.Configuration, caCertPool *x509.CertPool) (provider *PostgreSQLProvider) {
	provider = &PostgreSQLProvider{
		SQLProvider: NewSQLProvider(config, providerPostgres, "pgx", dsnPostgreSQL(config.Storage.PostgreSQL, caCertPool)),
	}

	// All providers have differing SELECT existing table statements.
	provider.sqlSelectExistingTables = queryPostgreSelectExistingTables

	// Specific alterations to this provider.
	// PostgreSQL doesn't have a UPSERT statement but has an ON CONFLICT operation instead.
	provider.sqlUpsertWebAuthnDevice = fmt.Sprintf(queryFmtUpsertWebAuthnDevicePostgreSQL, tableWebAuthnDevices)
	provider.sqlUpsertDuoDevice = fmt.Sprintf(queryFmtUpsertDuoDevicePostgreSQL, tableDuoDevices)
	provider.sqlUpsertTOTPConfig = fmt.Sprintf(queryFmtUpsertTOTPConfigurationPostgreSQL, tableTOTPConfigurations)
	provider.sqlUpsertPreferred2FAMethod = fmt.Sprintf(queryFmtUpsertPreferred2FAMethodPostgreSQL, tableUserPreferences)
	provider.sqlUpsertEncryptionValue = fmt.Sprintf(queryFmtUpsertEncryptionValuePostgreSQL, tableEncryption)
	provider.sqlUpsertOAuth2BlacklistedJTI = fmt.Sprintf(queryFmtUpsertOAuth2BlacklistedJTIPostgreSQL, tableOAuth2BlacklistedJTI)
	provider.sqlInsertOAuth2ConsentPreConfiguration = fmt.Sprintf(queryFmtInsertOAuth2ConsentPreConfigurationPostgreSQL, tableOAuth2ConsentPreConfiguration)

	// PostgreSQL requires rebinding of any query that contains a '?' placeholder to use the '$#' notation placeholders.
	provider.sqlFmtRenameTable = provider.db.Rebind(provider.sqlFmtRenameTable)

	provider.sqlSelectPreferred2FAMethod = provider.db.Rebind(provider.sqlSelectPreferred2FAMethod)
	provider.sqlSelectUserInfo = provider.db.Rebind(provider.sqlSelectUserInfo)

	provider.sqlInsertUserOpaqueIdentifier = provider.db.Rebind(provider.sqlInsertUserOpaqueIdentifier)
	provider.sqlSelectUserOpaqueIdentifier = provider.db.Rebind(provider.sqlSelectUserOpaqueIdentifier)
	provider.sqlSelectUserOpaqueIdentifierBySignature = provider.db.Rebind(provider.sqlSelectUserOpaqueIdentifierBySignature)

	provider.sqlSelectIdentityVerification = provider.db.Rebind(provider.sqlSelectIdentityVerification)
	provider.sqlInsertIdentityVerification = provider.db.Rebind(provider.sqlInsertIdentityVerification)
	provider.sqlConsumeIdentityVerification = provider.db.Rebind(provider.sqlConsumeIdentityVerification)

	provider.sqlSelectTOTPConfig = provider.db.Rebind(provider.sqlSelectTOTPConfig)
	provider.sqlUpdateTOTPConfigRecordSignIn = provider.db.Rebind(provider.sqlUpdateTOTPConfigRecordSignIn)
	provider.sqlUpdateTOTPConfigRecordSignInByUsername = provider.db.Rebind(provider.sqlUpdateTOTPConfigRecordSignInByUsername)
	provider.sqlDeleteTOTPConfig = provider.db.Rebind(provider.sqlDeleteTOTPConfig)
	provider.sqlSelectTOTPConfigs = provider.db.Rebind(provider.sqlSelectTOTPConfigs)

	provider.sqlSelectWebAuthnDevices = provider.db.Rebind(provider.sqlSelectWebAuthnDevices)
	provider.sqlSelectWebAuthnDevicesByUsername = provider.db.Rebind(provider.sqlSelectWebAuthnDevicesByUsername)
	provider.sqlUpdateWebAuthnDeviceRecordSignIn = provider.db.Rebind(provider.sqlUpdateWebAuthnDeviceRecordSignIn)
	provider.sqlUpdateWebAuthnDeviceRecordSignInByUsername = provider.db.Rebind(provider.sqlUpdateWebAuthnDeviceRecordSignInByUsername)
	provider.sqlDeleteWebAuthnDevice = provider.db.Rebind(provider.sqlDeleteWebAuthnDevice)
	provider.sqlDeleteWebAuthnDeviceByUsername = provider.db.Rebind(provider.sqlDeleteWebAuthnDeviceByUsername)
	provider.sqlDeleteWebAuthnDeviceByUsernameAndDescription = provider.db.Rebind(provider.sqlDeleteWebAuthnDeviceByUsernameAndDescription)

	provider.sqlSelectDuoDevice = provider.db.Rebind(provider.sqlSelectDuoDevice)
	provider.sqlDeleteDuoDevice = provider.db.Rebind(provider.sqlDeleteDuoDevice)

	provider.sqlInsertAuthenticationAttempt = provider.db.Rebind(provider.sqlInsertAuthenticationAttempt)
	provider.sqlSelectAuthenticationAttemptsByUsername = provider.db.Rebind(provider.sqlSelectAuthenticationAttemptsByUsername)

	provider.sqlInsertMigration = provider.db.Rebind(provider.sqlInsertMigration)
	provider.sqlSelectMigrations = provider.db.Rebind(provider.sqlSelectMigrations)
	provider.sqlSelectLatestMigration = provider.db.Rebind(provider.sqlSelectLatestMigration)

	provider.sqlSelectEncryptionValue = provider.db.Rebind(provider.sqlSelectEncryptionValue)

	provider.sqlSelectOAuth2ConsentPreConfigurations = provider.db.Rebind(provider.sqlSelectOAuth2ConsentPreConfigurations)

	provider.sqlInsertOAuth2ConsentSession = provider.db.Rebind(provider.sqlInsertOAuth2ConsentSession)
	provider.sqlUpdateOAuth2ConsentSessionSubject = provider.db.Rebind(provider.sqlUpdateOAuth2ConsentSessionSubject)
	provider.sqlUpdateOAuth2ConsentSessionResponse = provider.db.Rebind(provider.sqlUpdateOAuth2ConsentSessionResponse)
	provider.sqlUpdateOAuth2ConsentSessionGranted = provider.db.Rebind(provider.sqlUpdateOAuth2ConsentSessionGranted)
	provider.sqlSelectOAuth2ConsentSessionByChallengeID = provider.db.Rebind(provider.sqlSelectOAuth2ConsentSessionByChallengeID)

	provider.sqlInsertOAuth2AccessTokenSession = provider.db.Rebind(provider.sqlInsertOAuth2AccessTokenSession)
	provider.sqlRevokeOAuth2AccessTokenSession = provider.db.Rebind(provider.sqlRevokeOAuth2AccessTokenSession)
	provider.sqlRevokeOAuth2AccessTokenSessionByRequestID = provider.db.Rebind(provider.sqlRevokeOAuth2AccessTokenSessionByRequestID)
	provider.sqlDeactivateOAuth2AccessTokenSession = provider.db.Rebind(provider.sqlDeactivateOAuth2AccessTokenSession)
	provider.sqlDeactivateOAuth2AccessTokenSessionByRequestID = provider.db.Rebind(provider.sqlDeactivateOAuth2AccessTokenSessionByRequestID)
	provider.sqlSelectOAuth2AccessTokenSession = provider.db.Rebind(provider.sqlSelectOAuth2AccessTokenSession)

	provider.sqlInsertOAuth2AuthorizeCodeSession = provider.db.Rebind(provider.sqlInsertOAuth2AuthorizeCodeSession)
	provider.sqlRevokeOAuth2AuthorizeCodeSession = provider.db.Rebind(provider.sqlRevokeOAuth2AuthorizeCodeSession)
	provider.sqlRevokeOAuth2AuthorizeCodeSessionByRequestID = provider.db.Rebind(provider.sqlRevokeOAuth2AuthorizeCodeSessionByRequestID)
	provider.sqlDeactivateOAuth2AuthorizeCodeSession = provider.db.Rebind(provider.sqlDeactivateOAuth2AuthorizeCodeSession)
	provider.sqlDeactivateOAuth2AuthorizeCodeSessionByRequestID = provider.db.Rebind(provider.sqlDeactivateOAuth2AuthorizeCodeSessionByRequestID)
	provider.sqlSelectOAuth2AuthorizeCodeSession = provider.db.Rebind(provider.sqlSelectOAuth2AuthorizeCodeSession)

	provider.sqlInsertOAuth2OpenIDConnectSession = provider.db.Rebind(provider.sqlInsertOAuth2OpenIDConnectSession)
	provider.sqlRevokeOAuth2OpenIDConnectSession = provider.db.Rebind(provider.sqlRevokeOAuth2OpenIDConnectSession)
	provider.sqlRevokeOAuth2OpenIDConnectSessionByRequestID = provider.db.Rebind(provider.sqlRevokeOAuth2OpenIDConnectSessionByRequestID)
	provider.sqlDeactivateOAuth2OpenIDConnectSession = provider.db.Rebind(provider.sqlDeactivateOAuth2OpenIDConnectSession)
	provider.sqlDeactivateOAuth2OpenIDConnectSessionByRequestID = provider.db.Rebind(provider.sqlDeactivateOAuth2OpenIDConnectSessionByRequestID)
	provider.sqlSelectOAuth2OpenIDConnectSession = provider.db.Rebind(provider.sqlSelectOAuth2OpenIDConnectSession)

	provider.sqlInsertOAuth2PARContext = provider.db.Rebind(provider.sqlInsertOAuth2PARContext)
	provider.sqlRevokeOAuth2PARContext = provider.db.Rebind(provider.sqlRevokeOAuth2PARContext)
	provider.sqlSelectOAuth2PARContext = provider.db.Rebind(provider.sqlSelectOAuth2PARContext)

	provider.sqlInsertOAuth2PKCERequestSession = provider.db.Rebind(provider.sqlInsertOAuth2PKCERequestSession)
	provider.sqlRevokeOAuth2PKCERequestSession = provider.db.Rebind(provider.sqlRevokeOAuth2PKCERequestSession)
	provider.sqlRevokeOAuth2PKCERequestSessionByRequestID = provider.db.Rebind(provider.sqlRevokeOAuth2PKCERequestSessionByRequestID)
	provider.sqlDeactivateOAuth2PKCERequestSession = provider.db.Rebind(provider.sqlDeactivateOAuth2PKCERequestSession)
	provider.sqlDeactivateOAuth2PKCERequestSessionByRequestID = provider.db.Rebind(provider.sqlDeactivateOAuth2PKCERequestSessionByRequestID)
	provider.sqlSelectOAuth2PKCERequestSession = provider.db.Rebind(provider.sqlSelectOAuth2PKCERequestSession)

	provider.sqlInsertOAuth2RefreshTokenSession = provider.db.Rebind(provider.sqlInsertOAuth2RefreshTokenSession)
	provider.sqlRevokeOAuth2RefreshTokenSession = provider.db.Rebind(provider.sqlRevokeOAuth2RefreshTokenSession)
	provider.sqlRevokeOAuth2RefreshTokenSessionByRequestID = provider.db.Rebind(provider.sqlRevokeOAuth2RefreshTokenSessionByRequestID)
	provider.sqlDeactivateOAuth2RefreshTokenSession = provider.db.Rebind(provider.sqlDeactivateOAuth2RefreshTokenSession)
	provider.sqlDeactivateOAuth2RefreshTokenSessionByRequestID = provider.db.Rebind(provider.sqlDeactivateOAuth2RefreshTokenSessionByRequestID)
	provider.sqlSelectOAuth2RefreshTokenSession = provider.db.Rebind(provider.sqlSelectOAuth2RefreshTokenSession)

	provider.sqlSelectOAuth2BlacklistedJTI = provider.db.Rebind(provider.sqlSelectOAuth2BlacklistedJTI)

	provider.schema = config.Storage.PostgreSQL.Schema

	return provider
}

func dsnPostgreSQL(config *schema.PostgreSQLStorageConfiguration, globalCACertPool *x509.CertPool) (dsn string) {
	dsnConfig, _ := pgx.ParseConfig("")

	dsnConfig.Host = config.Address.SocketHostname()
	dsnConfig.Port = uint16(config.Address.Port())
	dsnConfig.Database = config.Database
	dsnConfig.User = config.Username
	dsnConfig.Password = config.Password
	dsnConfig.TLSConfig = loadPostgreSQLTLSConfig(config, globalCACertPool)
	dsnConfig.ConnectTimeout = config.Timeout
	dsnConfig.RuntimeParams = map[string]string{
		"search_path": config.Schema,
	}

	if dsnConfig.Port == 0 && config.Address.IsUnixDomainSocket() {
		dsnConfig.Port = 5432
	}

	return stdlib.RegisterConnConfig(dsnConfig)
}

func loadPostgreSQLTLSConfig(config *schema.PostgreSQLStorageConfiguration, globalCACertPool *x509.CertPool) (tlsConfig *tls.Config) {
	if config.TLS == nil && config.SSL == nil {
		return nil
	}

	if config.TLS != nil {
		return utils.NewTLSConfig(config.TLS, globalCACertPool)
	}

	ca, certs := loadPostgreSQLLegacyTLSConfig(config)

	switch config.SSL.Mode {
	case "disable":
		return nil
	default:
		var caCertPool *x509.CertPool

		switch ca {
		case nil:
			caCertPool = globalCACertPool
		default:
			caCertPool = globalCACertPool.Clone()
			caCertPool.AddCert(ca)
		}

		tlsConfig = &tls.Config{
			Certificates:       certs,
			RootCAs:            caCertPool,
			InsecureSkipVerify: true, //nolint:gosec
		}

		switch {
		case config.SSL.Mode == "require" && config.SSL.RootCertificate != "" || config.SSL.Mode == "verify-ca":
			tlsConfig.VerifyPeerCertificate = newPostgreSQLVerifyCAFunc(tlsConfig)
		case config.SSL.Mode == "verify-full":
			tlsConfig.InsecureSkipVerify = false
			tlsConfig.ServerName = config.Address.Hostname()
		}
	}

	return tlsConfig
}

func loadPostgreSQLLegacyTLSConfig(config *schema.PostgreSQLStorageConfiguration) (ca *x509.Certificate, certs []tls.Certificate) {
	var (
		err error
	)

	if config.SSL.RootCertificate != "" {
		var caPEMBlock []byte

		if caPEMBlock, err = os.ReadFile(config.SSL.RootCertificate); err != nil {
			return nil, nil
		}

		if ca, err = x509.ParseCertificate(caPEMBlock); err != nil {
			return nil, nil
		}
	}

	if config.SSL.Certificate != "" && config.SSL.Key != "" {
		var (
			keyPEMBlock  []byte
			certPEMBlock []byte
		)

		if keyPEMBlock, err = os.ReadFile(config.SSL.Key); err != nil {
			return nil, nil
		}

		if certPEMBlock, err = os.ReadFile(config.SSL.Certificate); err != nil {
			return nil, nil
		}

		var cert tls.Certificate

		if cert, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock); err != nil {
			return nil, nil
		}

		certs = []tls.Certificate{cert}
	}

	return ca, certs
}

func newPostgreSQLVerifyCAFunc(config *tls.Config) func(certificates [][]byte, _ [][]*x509.Certificate) (err error) {
	return func(certificates [][]byte, _ [][]*x509.Certificate) (err error) {
		certs := make([]*x509.Certificate, len(certificates))

		var cert *x509.Certificate

		for i, asn1Data := range certificates {
			if cert, err = x509.ParseCertificate(asn1Data); err != nil {
				return errors.New("failed to parse certificate from server: " + err.Error())
			}

			certs[i] = cert
		}

		// Leave DNSName empty to skip hostname verification.
		opts := x509.VerifyOptions{
			Roots:         config.RootCAs,
			Intermediates: x509.NewCertPool(),
		}

		// Skip the first cert because it's the leaf. All others
		// are intermediates.
		for _, cert = range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		_, err = certs[0].Verify(opts)

		return err
	}
}
