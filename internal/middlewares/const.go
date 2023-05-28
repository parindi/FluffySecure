package middlewares

import (
	"errors"

	"github.com/valyala/fasthttp"
)

var (
	headerXAutheliaURL = []byte("X-Authelia-URL")

	headerAccept        = []byte(fasthttp.HeaderAccept)
	headerContentLength = []byte(fasthttp.HeaderContentLength)
	headerLocation      = []byte(fasthttp.HeaderLocation)

	headerXForwardedProto = []byte(fasthttp.HeaderXForwardedProto)
	headerXForwardedHost  = []byte(fasthttp.HeaderXForwardedHost)
	headerXForwardedFor   = []byte(fasthttp.HeaderXForwardedFor)
	headerXRequestedWith  = []byte(fasthttp.HeaderXRequestedWith)

	headerXForwardedURI    = []byte("X-Forwarded-URI")
	headerXOriginalURL     = []byte("X-Original-URL")
	headerXOriginalMethod  = []byte("X-Original-Method")
	headerXForwardedMethod = []byte("X-Forwarded-Method")

	headerVary   = []byte(fasthttp.HeaderVary)
	headerAllow  = []byte(fasthttp.HeaderAllow)
	headerOrigin = []byte(fasthttp.HeaderOrigin)

	headerAccessControlAllowCredentials = []byte(fasthttp.HeaderAccessControlAllowCredentials)
	headerAccessControlAllowHeaders     = []byte(fasthttp.HeaderAccessControlAllowHeaders)
	headerAccessControlAllowMethods     = []byte(fasthttp.HeaderAccessControlAllowMethods)
	headerAccessControlAllowOrigin      = []byte(fasthttp.HeaderAccessControlAllowOrigin)
	headerAccessControlMaxAge           = []byte(fasthttp.HeaderAccessControlMaxAge)
	headerAccessControlRequestHeaders   = []byte(fasthttp.HeaderAccessControlRequestHeaders)
	headerAccessControlRequestMethod    = []byte(fasthttp.HeaderAccessControlRequestMethod)

	headerXContentTypeOptions   = []byte(fasthttp.HeaderXContentTypeOptions)
	headerReferrerPolicy        = []byte(fasthttp.HeaderReferrerPolicy)
	headerXFrameOptions         = []byte(fasthttp.HeaderXFrameOptions)
	headerPragma                = []byte(fasthttp.HeaderPragma)
	headerCacheControl          = []byte(fasthttp.HeaderCacheControl)
	headerXXSSProtection        = []byte(fasthttp.HeaderXXSSProtection)
	headerContentSecurityPolicy = []byte(fasthttp.HeaderContentSecurityPolicy)

	headerPermissionsPolicy = []byte("Permissions-Policy")
)

var (
	headerValueFalse           = []byte("false")
	headerValueTrue            = []byte("true")
	headerValueMaxAge          = []byte("100")
	headerValueVary            = []byte("Accept-Encoding, Origin")
	headerValueVaryWildcard    = []byte("Accept-Encoding")
	headerValueOriginWildcard  = []byte("*")
	headerValueZero            = []byte("0")
	headerValueCSPNone         = []byte("default-src 'none'")
	headerValueCSPNoneFormPost = []byte("default-src 'none'; script-src 'sha256-skflBqA90WuHvoczvimLdj49ExKdizFjX2Itd6xKZdU='")

	headerValueNoSniff                 = []byte("nosniff")
	headerValueStrictOriginCrossOrigin = []byte("strict-origin-when-cross-origin")
	headerValueSameOrigin              = []byte("SAMEORIGIN")
	headerValueNoCache                 = []byte("no-cache")
	headerValueNoStore                 = []byte("no-store")
	headerValueXSSModeBlock            = []byte("1; mode=block")
	headerValueCohort                  = []byte("interest-cohort=()")
)

const (
	strProtoHTTPS = "https"
	strProtoHTTP  = "http"
	strSlash      = "/"

	queryArgRedirect    = "rd"
	queryArgAutheliaURL = "authelia_url"
	queryArgToken       = "token"
)

var (
	protoHTTPS = []byte(strProtoHTTPS)
	protoHTTP  = []byte(strProtoHTTP)

	qryArgRedirect    = []byte(queryArgRedirect)
	qryArgAutheliaURL = []byte(queryArgAutheliaURL)

	keyUserValueBaseURL   = []byte("base_url")
	keyUserValueAuthzPath = []byte("authz_path")

	// UserValueKeyFormPost is the User Value key where we indicate the form_post response mode.
	UserValueKeyFormPost = []byte("form_post")

	headerSeparator = []byte(", ")

	contentTypeTextPlain       = []byte("text/plain; charset=utf-8")
	contentTypeTextHTML        = []byte("text/html; charset=utf-8")
	contentTypeApplicationJSON = []byte("application/json; charset=utf-8")
	contentTypeApplicationYAML = []byte("application/yaml; charset=utf-8")
)

const (
	headerValueXRequestedWithXHR = "XMLHttpRequest"
)

var okMessageBytes = []byte("{\"status\":\"OK\"}")

const (
	messageOperationFailed                      = "Operation failed"
	messageIdentityVerificationTokenAlreadyUsed = "The identity verification token has already been used"
	messageIdentityVerificationTokenHasExpired  = "The identity verification token has expired"
)

var protoHostSeparator = []byte("://")

var errPasswordPolicyNoMet = errors.New("the supplied password does not met the security policy")
