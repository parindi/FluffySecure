package storage

const (
	queryFmtSelectMigrations = `
		SELECT id, applied, version_before, version_after, application_version
		FROM %s;`

	queryFmtSelectLatestMigration = `
		SELECT id, applied, version_before, version_after, application_version
		FROM %s
		ORDER BY id DESC
		LIMIT 1;`

	queryFmtInsertMigration = `
		INSERT INTO %s (applied, version_before, version_after, application_version)
		VALUES (?, ?, ?, ?);`
)

const (
	queryMySQLSelectExistingTables = `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE' AND table_schema = database();`

	queryPostgreSelectExistingTables = `
		SELECT table_name
		FROM information_schema.tables
		WHERE table_type = 'BASE TABLE' AND table_schema = $1;`

	querySQLiteSelectExistingTables = `
		SELECT name
		FROM sqlite_master
		WHERE type = 'table';`
)

const (
	queryFmtSelectUserInfo = `
		SELECT second_factor_method, (SELECT EXISTS (SELECT id FROM %s WHERE username = ?)) AS has_totp, (SELECT EXISTS (SELECT id FROM %s WHERE username = ?)) AS has_webauthn, (SELECT EXISTS (SELECT id FROM %s WHERE username = ?)) AS has_duo
		FROM %s
		WHERE username = ?;`

	queryFmtSelectPreferred2FAMethod = `
		SELECT second_factor_method
		FROM %s
		WHERE username = ?;`

	queryFmtUpsertPreferred2FAMethod = `
		REPLACE INTO %s (username, second_factor_method)
		VALUES (?, ?);`

	queryFmtUpsertPreferred2FAMethodPostgreSQL = `
		INSERT INTO %s (username, second_factor_method)
		VALUES ($1, $2)
			ON CONFLICT (username)
			DO UPDATE SET second_factor_method = $2;`
)

const (
	queryFmtSelectIdentityVerification = `
		SELECT id, jti, iat, issued_ip, exp, username, action, consumed, consumed_ip
		FROM %s
		WHERE jti = ?;`

	queryFmtInsertIdentityVerification = `
		INSERT INTO %s (jti, iat, issued_ip, exp, username, action)
		VALUES (?, ?, ?, ?, ?, ?);`

	queryFmtConsumeIdentityVerification = `
		UPDATE %s
		SET consumed = CURRENT_TIMESTAMP, consumed_ip = ?
		WHERE jti = ?;`
)

const (
	queryFmtSelectTOTPConfiguration = `
		SELECT id, username, issuer, algorithm, digits, period, secret
		FROM %s
		WHERE username = ?;`

	queryFmtSelectTOTPConfigurations = `
		SELECT id, username, issuer, algorithm, digits, period, secret
		FROM %s
		LIMIT ?
		OFFSET ?;`

	queryFmtSelectTOTPConfigurationsEncryptedData = `
		SELECT id, secret
		FROM %s;`

	//nolint:gosec // These are not hardcoded credentials it's a query to obtain credentials.
	queryFmtUpdateTOTPConfigurationSecret = `
		UPDATE %s
		SET secret = ?
		WHERE id = ?;`

	queryFmtUpsertTOTPConfiguration = `
		REPLACE INTO %s (created_at, last_used_at, username, issuer, algorithm, digits, period, secret)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtUpsertTOTPConfigurationPostgreSQL = `
		INSERT INTO %s (created_at, last_used_at, username, issuer, algorithm, digits, period, secret)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (username)
			DO UPDATE SET created_at = $1, last_used_at = $2, issuer = $4, algorithm = $5, digits = $6, period = $7, secret = $8;`

	queryFmtUpdateTOTPConfigRecordSignIn = `
		UPDATE %s
		SET last_used_at = ?
		WHERE id = ?;`

	queryFmtUpdateTOTPConfigRecordSignInByUsername = `
		UPDATE %s
		SET last_used_at = ?
		WHERE username = ?;`

	queryFmtDeleteTOTPConfiguration = `
		DELETE FROM %s
		WHERE username = ?;`
)

const (
	queryFmtSelectWebAuthnDevices = `
		SELECT id, created_at, last_used_at, rpid, username, description, kid, public_key, attestation_type, transport, aaguid, sign_count, clone_warning
		FROM %s
		LIMIT ?
		OFFSET ?;`

	queryFmtSelectWebAuthnDevicesEncryptedData = `
		SELECT id, public_key
		FROM %s;`

	queryFmtSelectWebAuthnDevicesByUsername = `
		SELECT id, created_at, last_used_at, rpid, username, description, kid, public_key, attestation_type, transport, aaguid, sign_count, clone_warning
		FROM %s
		WHERE username = ?;`

	queryFmtUpdateWebAuthnDevicePublicKey = `
		UPDATE %s
		SET public_key = ?
		WHERE id = ?;`

	queryFmtUpdateWebAuthnDeviceRecordSignIn = `
		UPDATE %s
		SET
			rpid = ?, last_used_at = ?, sign_count = ?,
			clone_warning = CASE clone_warning WHEN TRUE THEN TRUE ELSE ? END
		WHERE id = ?;`

	queryFmtUpdateWebAuthnDeviceRecordSignInByUsername = `
		UPDATE %s
		SET
			rpid = ?, last_used_at = ?, sign_count = ?,
			clone_warning = CASE clone_warning WHEN TRUE THEN TRUE ELSE ? END
		WHERE username = ? AND kid = ?;`

	queryFmtUpsertWebAuthnDevice = `
		REPLACE INTO %s (created_at, last_used_at, rpid, username, description, kid, public_key, attestation_type, transport, aaguid, sign_count, clone_warning)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtUpsertWebAuthnDevicePostgreSQL = `
		INSERT INTO %s (created_at, last_used_at, rpid, username, description, kid, public_key, attestation_type, transport, aaguid, sign_count, clone_warning)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			ON CONFLICT (username, description)
			DO UPDATE SET created_at = $1, last_used_at = $2, rpid = $3, kid = $6, public_key = $7, attestation_type = $8, transport = $9, aaguid = $10, sign_count = $11, clone_warning = $12;`

	queryFmtDeleteWebAuthnDevice = `
		DELETE FROM %s
		WHERE kid = ?;`

	queryFmtDeleteWebAuthnDeviceByUsername = `
		DELETE FROM %s
		WHERE username = ?;`

	queryFmtDeleteWebAuthnDeviceByUsernameAndDescription = `
		DELETE FROM %s
		WHERE username = ? AND description = ?;`
)

const (
	queryFmtUpsertDuoDevice = `
		REPLACE INTO %s (username, device, method)
		VALUES (?, ?, ?);`

	queryFmtUpsertDuoDevicePostgreSQL = `
		INSERT INTO %s (username, device, method)
		VALUES ($1, $2, $3)
			ON CONFLICT (username)
			DO UPDATE SET device = $2, method = $3;`

	queryFmtDeleteDuoDevice = `
		DELETE
		FROM %s
		WHERE username = ?;`

	queryFmtSelectDuoDevice = `
		SELECT id, username, device, method
		FROM %s
		WHERE username = ?
		ORDER BY id;`
)

const (
	queryFmtInsertAuthenticationLogEntry = `
		INSERT INTO %s (time, successful, banned, username, auth_type, remote_ip, request_uri, request_method)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtSelect1FAAuthenticationLogEntryByUsername = `
		SELECT time, successful, username
		FROM %s
		WHERE time > ? AND username = ? AND auth_type = '1FA' AND banned = FALSE
		ORDER BY time DESC
		LIMIT ?
		OFFSET ?;`
)

const (
	queryFmtSelectEncryptionValue = `
		SELECT (value)
        FROM %s
        WHERE name = ?`

	queryFmtUpsertEncryptionValue = `
		REPLACE INTO %s (name, value)
		VALUES (?, ?);`

	queryFmtUpsertEncryptionValuePostgreSQL = `
		INSERT INTO %s (name, value)
		VALUES ($1, $2)
			ON CONFLICT (name)
			DO UPDATE SET value = $2;`
)

const (
	queryFmtSelectOAuth2ConsentPreConfigurations = `
		SELECT id, client_id, subject, created_at, expires_at, revoked, scopes, audience
		FROM %s
		WHERE client_id = ? AND subject = ? AND
			  revoked = FALSE AND (expires_at IS NULL OR expires_at >= CURRENT_TIMESTAMP);`

	queryFmtInsertOAuth2ConsentPreConfiguration = `
		INSERT INTO %s (client_id, subject, created_at, expires_at, revoked, scopes, audience)
		VALUES(?, ?, ?, ?, ?, ?, ?);`

	queryFmtInsertOAuth2ConsentPreConfigurationPostgreSQL = `
		INSERT INTO %s (client_id, subject, created_at, expires_at, revoked, scopes, audience)
		VALUES($1, $2, $3, $4, $5, $6, $7)
		RETURNING id;`

	queryFmtSelectOAuth2ConsentSessionByChallengeID = `
		SELECT id, challenge_id, client_id, subject, authorized, granted, requested_at, responded_at,
		form_data, requested_scopes, granted_scopes, requested_audience, granted_audience, preconfiguration
		FROM %s
		WHERE challenge_id = ?;`

	queryFmtInsertOAuth2ConsentSession = `
		INSERT INTO %s (challenge_id, client_id, subject, authorized, granted, requested_at, responded_at,
		form_data, requested_scopes, granted_scopes, requested_audience, granted_audience, preconfiguration)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtUpdateOAuth2ConsentSessionSubject = `
		UPDATE %s
		SET subject = ?
		WHERE id = ?;`

	queryFmtUpdateOAuth2ConsentSessionSessionData = `
		UPDATE %s
		SET session_data = ?
		WHERE id = ?;`

	queryFmtUpdateOAuth2ConsentSessionResponse = `
		UPDATE %s
		SET authorized = ?, responded_at = CURRENT_TIMESTAMP, granted_scopes = ?, granted_audience = ?, preconfiguration = ?
		WHERE id = ? AND responded_at IS NULL;`

	queryFmtUpdateOAuth2ConsentSessionGranted = `
		UPDATE %s
		SET granted = TRUE
		WHERE id = ? AND responded_at IS NOT NULL;`

	queryFmtSelectOAuth2Session = `
		SELECT id, challenge_id, request_id, client_id, signature, subject, requested_at,
		requested_scopes, granted_scopes, requested_audience, granted_audience,
		active, revoked, form_data, session_data
		FROM %s
		WHERE signature = ? AND revoked = FALSE;`

	queryFmtSelectOAuth2SessionEncryptedData = `
		SELECT id, session_data
		FROM %s;`

	queryFmtInsertOAuth2Session = `
		INSERT INTO %s (challenge_id, request_id, client_id, signature, subject, requested_at,
		requested_scopes, granted_scopes, requested_audience, granted_audience,
		active, revoked, form_data, session_data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtRevokeOAuth2Session = `
		UPDATE %s
		SET revoked = TRUE
		WHERE signature = ?;`

	queryFmtRevokeOAuth2SessionByRequestID = `
		UPDATE %s
		SET revoked = TRUE
		WHERE request_id = ?;`

	queryFmtDeactivateOAuth2Session = `
		UPDATE %s
		SET active = FALSE
		WHERE signature = ?;`

	queryFmtDeactivateOAuth2SessionByRequestID = `
		UPDATE %s
		SET active = FALSE
		WHERE request_id = ?;`

	queryFmtSelectOAuth2PARContext = `
		SELECT id, signature, request_id, client_id, requested_at, scopes, audience,
		handled_response_types, response_mode, response_mode_default, revoked,
		form_data, session_data
		FROM %s
		WHERE signature = ? AND revoked = FALSE;`

	queryFmtInsertOAuth2PARContext = `
		INSERT INTO %s (signature, request_id, client_id, requested_at, scopes, audience,
		handled_response_types, response_mode, response_mode_default, revoked,
		form_data, session_data)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	queryFmtSelectOAuth2BlacklistedJTI = `
		SELECT id, signature, expires_at
		FROM %s
		WHERE signature = ?;`

	queryFmtUpsertOAuth2BlacklistedJTI = `
		REPLACE INTO %s (signature, expires_at)
		VALUES(?, ?);`

	queryFmtUpsertOAuth2BlacklistedJTIPostgreSQL = `
		INSERT INTO %s (signature, expires_at)
		VALUES ($1, $2)
			ON CONFLICT (signature)
			DO UPDATE SET expires_at = $2;`
)

const (
	queryFmtInsertUserOpaqueIdentifier = `
		INSERT INTO %s (service, sector_id, username, identifier)
		VALUES(?, ?, ?, ?);`

	queryFmtSelectUserOpaqueIdentifier = `
		SELECT id, service, sector_id, username, identifier
		FROM %s
		WHERE identifier = ?;`

	queryFmtSelectUserOpaqueIdentifierBySignature = `
		SELECT id, service, sector_id, username, identifier
		FROM %s
		WHERE service = ? AND sector_id = ? AND username = ?;`

	queryFmtSelectUserOpaqueIdentifiers = `
		SELECT id, service, sector_id, username, identifier
		FROM %s;`
)
