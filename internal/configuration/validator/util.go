package validator

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rsa"
	"fmt"
	"strings"

	"golang.org/x/net/publicsuffix"
	"gopkg.in/square/go-jose.v2"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/oidc"
	"github.com/authelia/authelia/v4/internal/utils"
)

func isCookieDomainAPublicSuffix(domain string) (valid bool) {
	var suffix string

	suffix, _ = publicsuffix.PublicSuffix(domain)

	return len(strings.TrimLeft(domain, ".")) == len(suffix)
}

func strJoinOr(items []string) string {
	return strJoinComma("or", items)
}

func strJoinAnd(items []string) string {
	return strJoinComma("and", items)
}

func strJoinComma(word string, items []string) string {
	if word == "" {
		return buildJoinedString(",", "", "'", items)
	}

	return buildJoinedString(",", word, "'", items)
}

func buildJoinedString(sep, sepFinal, quote string, items []string) string {
	n := len(items)

	if n == 0 {
		return ""
	}

	b := &strings.Builder{}

	for i := 0; i < n; i++ {
		if quote != "" {
			b.WriteString(quote)
		}

		b.WriteString(items[i])

		if quote != "" {
			b.WriteString(quote)
		}

		if i == (n - 1) {
			continue
		}

		if sep != "" {
			if sepFinal == "" || n != 2 {
				b.WriteString(sep)
			}

			b.WriteString(" ")
		}

		if sepFinal != "" && i == (n-2) {
			b.WriteString(strings.Trim(sepFinal, " "))
			b.WriteString(" ")
		}
	}

	return b.String()
}

func validateList(values, valid []string, chkDuplicate bool) (invalid, duplicates []string) { //nolint:unparam
	chkValid := len(valid) != 0

	for i, value := range values {
		if chkValid {
			if !utils.IsStringInSlice(value, valid) {
				invalid = append(invalid, value)

				// Skip checking duplicates for invalid values.
				continue
			}
		}

		if chkDuplicate {
			for j, valueAlt := range values {
				if i == j {
					continue
				}

				if value != valueAlt {
					continue
				}

				if utils.IsStringInSlice(value, duplicates) {
					continue
				}

				duplicates = append(duplicates, value)
			}
		}
	}

	return
}

type JWKProperties struct {
	Use       string
	Algorithm string
	Bits      int
	Curve     elliptic.Curve
}

func schemaJWKGetProperties(jwk schema.JWK) (properties *JWKProperties, err error) {
	switch key := jwk.Key.(type) {
	case nil:
		return nil, nil
	case ed25519.PrivateKey, ed25519.PublicKey:
		return &JWKProperties{}, nil
	case *rsa.PrivateKey:
		if key.PublicKey.N == nil {
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgRSAUsingSHA256, 0, nil}, nil
		}

		return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgRSAUsingSHA256, key.Size(), nil}, nil
	case *rsa.PublicKey:
		if key.N == nil {
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgRSAUsingSHA256, 0, nil}, nil
		}

		return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgRSAUsingSHA256, key.Size(), nil}, nil
	case *ecdsa.PublicKey:
		switch key.Curve {
		case elliptic.P256():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP256AndSHA256, -1, key.Curve}, nil
		case elliptic.P384():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP384AndSHA384, -1, key.Curve}, nil
		case elliptic.P521():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP521AndSHA512, -1, key.Curve}, nil
		default:
			return &JWKProperties{oidc.KeyUseSignature, "", -1, key.Curve}, nil
		}
	case *ecdsa.PrivateKey:
		switch key.Curve {
		case elliptic.P256():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP256AndSHA256, -1, key.Curve}, nil
		case elliptic.P384():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP384AndSHA384, -1, key.Curve}, nil
		case elliptic.P521():
			return &JWKProperties{oidc.KeyUseSignature, oidc.SigningAlgECDSAUsingP521AndSHA512, -1, key.Curve}, nil
		default:
			return &JWKProperties{oidc.KeyUseSignature, "", -1, key.Curve}, nil
		}
	default:
		return nil, fmt.Errorf("the key type '%T' is unknown or not valid for the configuration", key)
	}
}

func jwkCalculateThumbprint(key schema.CryptographicKey) (thumbprintStr string, err error) {
	j := jose.JSONWebKey{}

	switch k := key.(type) {
	case schema.CryptographicPrivateKey:
		j.Key = k.Public()
	case *rsa.PublicKey, *ecdsa.PublicKey, ed25519.PublicKey:
		j.Key = k
	default:
		return "", nil
	}

	var thumbprint []byte

	if thumbprint, err = j.Thumbprint(crypto.SHA256); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", thumbprint)[:6], nil
}

func getResponseObjectAlgFromKID(config *schema.OpenIDConnect, kid, alg string) string {
	for _, jwk := range config.IssuerPrivateKeys {
		if kid == jwk.KeyID {
			return jwk.Algorithm
		}
	}

	return alg
}
