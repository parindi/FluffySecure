package validator

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
)

type AccessControl struct {
	suite.Suite
	config    *schema.Configuration
	validator *schema.StructValidator
}

func (suite *AccessControl) SetupTest() {
	suite.validator = schema.NewStructValidator()
	suite.config = &schema.Configuration{
		AccessControl: schema.AccessControlConfiguration{
			DefaultPolicy: policyDeny,

			Networks: schema.DefaultACLNetwork,
			Rules:    schema.DefaultACLRule,
		},
	}
}

func (suite *AccessControl) TestShouldValidateCompleteConfiguration() {
	ValidateAccessControl(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Assert().Len(suite.validator.Errors(), 0)
}

func (suite *AccessControl) TestShouldValidateEitherDomainsOrDomainsRegex() {
	domainsRegex := regexp.MustCompile(`^abc.example.com$`)

	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: []string{"abc.example.com"},
			Policy:  "bypass",
		},
		{
			DomainsRegex: []regexp.Regexp{*domainsRegex},
			Policy:       "bypass",
		},
		{
			Policy: "bypass",
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	assert.EqualError(suite.T(), suite.validator.Errors()[0], "access control: rule #3: option 'domain' or 'domain_regex' must be present but are both absent")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidDefaultPolicy() {
	suite.config.AccessControl.DefaultPolicy = testInvalid

	ValidateAccessControl(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: option 'default_policy' must be one of 'bypass', 'one_factor', 'two_factor', or 'deny' but it's configured as 'invalid'")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidNetworkGroupNetwork() {
	suite.config.AccessControl.Networks = []schema.ACLNetwork{
		{
			Name:     "internal",
			Networks: []string{"abc.def.ghi.jkl"},
		},
	}

	ValidateAccessControl(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: networks: network group 'internal' is invalid: the network 'abc.def.ghi.jkl' is not a valid IP or CIDR notation")
}

func (suite *AccessControl) TestShouldRaiseWarningOnBadDomain() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: []string{"*example.com"},
			Policy:  "one_factor",
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 1)
	suite.Require().Len(suite.validator.Errors(), 0)

	suite.Assert().EqualError(suite.validator.Warnings()[0], "access control: rule #1: domain #1: domain '*example.com' is ineffective and should probably be '*.example.com' instead")
}

func (suite *AccessControl) TestShouldRaiseErrorWithNoRulesDefined() {
	suite.config.AccessControl.Rules = []schema.ACLRule{}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: 'default_policy' option 'deny' is invalid: when no rules are specified it must be 'two_factor' or 'one_factor'")
}

func (suite *AccessControl) TestShouldRaiseWarningWithNoRulesDefined() {
	suite.config.AccessControl.Rules = []schema.ACLRule{}

	suite.config.AccessControl.DefaultPolicy = policyTwoFactor

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Errors(), 0)
	suite.Require().Len(suite.validator.Warnings(), 1)

	suite.Assert().EqualError(suite.validator.Warnings()[0], "access control: no rules have been specified so the 'default_policy' of 'two_factor' is going to be applied to all requests")
}

func (suite *AccessControl) TestShouldRaiseErrorsWithEmptyRules() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{},
		{
			Policy: "wrong",
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 4)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1: option 'domain' or 'domain_regex' must be present but are both absent")
	suite.Assert().EqualError(suite.validator.Errors()[1], "access control: rule #1: option 'policy' must be present but it's absent")
	suite.Assert().EqualError(suite.validator.Errors()[2], "access control: rule #2: option 'domain' or 'domain_regex' must be present but are both absent")
	suite.Assert().EqualError(suite.validator.Errors()[3], "access control: rule #2: option 'policy' must be one of 'bypass', 'one_factor', 'two_factor', or 'deny' but it's configured as 'wrong'")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidPolicy() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: []string{"public.example.com"},
			Policy:  testInvalid,
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): option 'policy' must be one of 'bypass', 'one_factor', 'two_factor', or 'deny' but it's configured as 'invalid'")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidNetwork() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains:  []string{"public.example.com"},
			Policy:   "bypass",
			Networks: []string{"abc.def.ghi.jkl/32"},
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): the network 'abc.def.ghi.jkl/32' is not a valid Group Name, IP, or CIDR notation")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidMethod() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: []string{"public.example.com"},
			Policy:  "bypass",
			Methods: []string{fasthttp.MethodGet, "HOP"},
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): option 'methods' must only have the values 'GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE', 'TRACE', 'CONNECT', 'OPTIONS', 'COPY', 'LOCK', 'MKCOL', 'MOVE', 'PROPFIND', 'PROPPATCH', or 'UNLOCK' but the values 'HOP' are present")
}

func (suite *AccessControl) TestShouldRaiseErrorDuplicateMethod() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: []string{"public.example.com"},
			Policy:  "bypass",
			Methods: []string{fasthttp.MethodGet, fasthttp.MethodGet},
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): option 'methods' must have unique values but the values 'GET' are duplicated")
}

func (suite *AccessControl) TestShouldRaiseErrorInvalidSubject() {
	domains := []string{"public.example.com"}
	subjects := [][]string{{testInvalid}}
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains:  domains,
			Policy:   "bypass",
			Subjects: subjects,
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Require().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 2)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): 'subject' option 'invalid' is invalid: must start with 'user:' or 'group:'")
	suite.Assert().EqualError(suite.validator.Errors()[1], fmt.Sprintf(errAccessControlRuleBypassPolicyInvalidWithSubjects, ruleDescriptor(1, suite.config.AccessControl.Rules[0])))
}

func (suite *AccessControl) TestShouldRaiseErrorBypassWithSubjectDomainRegexGroup() {
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			DomainsRegex: MustCompileRegexps([]string{`^(?P<User>\w+)\.example\.com$`}),
			Policy:       "bypass",
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Require().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 1)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1: 'policy' option 'bypass' is not supported when 'domain_regex' option contains the user or group named matches. For more information see: https://www.authelia.com/c/acl-match-concept-2")
}

func (suite *AccessControl) TestShouldSetQueryDefaults() {
	domains := []string{"public.example.com"}
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "", Key: "example"},
				},
				{
					{Operator: "", Key: "example", Value: "test"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "pattern", Key: "a", Value: "^(x|y|z)$"},
				},
			},
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Assert().Len(suite.validator.Errors(), 0)

	suite.Assert().Equal("present", suite.config.AccessControl.Rules[0].Query[0][0].Operator)
	suite.Assert().Equal("equal", suite.config.AccessControl.Rules[0].Query[1][0].Operator)

	suite.Require().Len(suite.config.AccessControl.Rules, 2)
	suite.Require().Len(suite.config.AccessControl.Rules[1].Query, 1)
	suite.Require().Len(suite.config.AccessControl.Rules[1].Query[0], 1)

	t := &regexp.Regexp{}

	suite.Assert().IsType(t, suite.config.AccessControl.Rules[1].Query[0][0].Value)
}

func (suite *AccessControl) TestShouldErrorOnInvalidRulesQuery() {
	domains := []string{"public.example.com"}
	suite.config.AccessControl.Rules = []schema.ACLRule{
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "equal", Key: "example"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "present"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "present", Key: "a"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "absent", Key: "a"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "not", Key: "a", Value: "a"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "pattern", Key: "a", Value: "(bad pattern"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "present", Key: "a", Value: "not good"},
				},
			},
		},
		{
			Domains: domains,
			Policy:  "bypass",
			Query: [][]schema.ACLQueryRule{
				{
					{Operator: "present", Key: "a", Value: 5},
				},
			},
		},
	}

	ValidateRules(suite.config, suite.validator)

	suite.Assert().Len(suite.validator.Warnings(), 0)
	suite.Require().Len(suite.validator.Errors(), 7)

	suite.Assert().EqualError(suite.validator.Errors()[0], "access control: rule #1 (domain 'public.example.com'): query: option 'value' must be present when the option 'operator' is 'equal' but it's absent")
	suite.Assert().EqualError(suite.validator.Errors()[1], "access control: rule #2 (domain 'public.example.com'): query: option 'key' is required but it's absent")
	suite.Assert().EqualError(suite.validator.Errors()[2], "access control: rule #5 (domain 'public.example.com'): query: option 'key' is required but it's absent")
	suite.Assert().EqualError(suite.validator.Errors()[3], "access control: rule #6 (domain 'public.example.com'): query: option 'operator' must be one of 'present', 'absent', 'equal', 'not equal', 'pattern', or 'not pattern' but it's configured as 'not'")
	suite.Assert().EqualError(suite.validator.Errors()[4], "access control: rule #7 (domain 'public.example.com'): query: option 'value' is invalid: error parsing regexp: missing closing ): `(bad pattern`")
	suite.Assert().EqualError(suite.validator.Errors()[5], "access control: rule #8 (domain 'public.example.com'): query: option 'value' must not be present when the option 'operator' is 'present' but it's present")
	suite.Assert().EqualError(suite.validator.Errors()[6], "access control: rule #9 (domain 'public.example.com'): query: option 'value' is invalid: expected type was string but got int")
}

func TestAccessControl(t *testing.T) {
	suite.Run(t, new(AccessControl))
}

func TestShouldReturnCorrectResultsForValidNetworkGroups(t *testing.T) {
	config := schema.AccessControlConfiguration{
		Networks: schema.DefaultACLNetwork,
	}

	validNetwork := IsNetworkGroupValid(config, "internal")
	invalidNetwork := IsNetworkGroupValid(config, loopback)

	assert.True(t, validNetwork)
	assert.False(t, invalidNetwork)
}

func MustCompileRegexps(exps []string) (regexps []regexp.Regexp) {
	regexps = make([]regexp.Regexp, len(exps))

	for i, exp := range exps {
		regexps[i] = *regexp.MustCompile(exp)
	}

	return regexps
}
