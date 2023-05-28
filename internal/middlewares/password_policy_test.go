package middlewares

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
)

func TestNewPasswordPolicyProvider(t *testing.T) {
	testCases := []struct {
		desc     string
		have     schema.PasswordPolicyConfiguration
		expected PasswordPolicyProvider
	}{
		{
			desc:     "ShouldReturnUnconfiguredProvider",
			have:     schema.PasswordPolicyConfiguration{},
			expected: &StandardPasswordPolicyProvider{},
		},
		{
			desc:     "ShouldReturnProviderWhenZxcvbn",
			have:     schema.PasswordPolicyConfiguration{ZXCVBN: schema.PasswordPolicyZXCVBNParams{Enabled: true, MinScore: 10}},
			expected: &ZXCVBNPasswordPolicyProvider{minScore: 10},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWithMin",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8}},
			expected: &StandardPasswordPolicyProvider{min: 8},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWitHMinMax",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, MaxLength: 100}},
			expected: &StandardPasswordPolicyProvider{min: 8, max: 100},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWithMinLowercase",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, RequireLowercase: true}},
			expected: &StandardPasswordPolicyProvider{min: 8, patterns: []regexp.Regexp{*regexp.MustCompile(`[a-z]+`)}},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWithMinLowercaseUppercase",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, RequireLowercase: true, RequireUppercase: true}},
			expected: &StandardPasswordPolicyProvider{min: 8, patterns: []regexp.Regexp{*regexp.MustCompile(`[a-z]+`), *regexp.MustCompile(`[A-Z]+`)}},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWithMinLowercaseUppercaseNumber",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, RequireLowercase: true, RequireUppercase: true, RequireNumber: true}},
			expected: &StandardPasswordPolicyProvider{min: 8, patterns: []regexp.Regexp{*regexp.MustCompile(`[a-z]+`), *regexp.MustCompile(`[A-Z]+`), *regexp.MustCompile(`[0-9]+`)}},
		},
		{
			desc:     "ShouldReturnConfiguredProviderWithMinLowercaseUppercaseSpecial",
			have:     schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, RequireLowercase: true, RequireUppercase: true, RequireSpecial: true}},
			expected: &StandardPasswordPolicyProvider{min: 8, patterns: []regexp.Regexp{*regexp.MustCompile(`[a-z]+`), *regexp.MustCompile(`[A-Z]+`), *regexp.MustCompile(`[^a-zA-Z0-9]+`)}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual := NewPasswordPolicyProvider(tc.have)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestPasswordPolicyProvider_Validate(t *testing.T) {
	testCases := []struct {
		desc     string
		config   schema.PasswordPolicyConfiguration
		have     []string
		expected []error
	}{
		{
			desc:     "ShouldValidateAllPasswords",
			config:   schema.PasswordPolicyConfiguration{},
			have:     []string{"a", "1", "a really str0ng pass12nm3kjl12word@@#4"},
			expected: []error{nil, nil, nil},
		},
		{
			desc:     "ShouldValidatePasswordMinLength",
			config:   schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8}},
			have:     []string{"a", "b123", "1111111", "aaaaaaaa", "1o23nm1kio2n3k12jn"},
			expected: []error{errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, nil, nil},
		},
		{
			desc:   "ShouldValidatePasswordMaxLength",
			config: schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MaxLength: 30}},
			have: []string{
				"a1234567894654wkjnkjasnskjandkjansdkjnas",
				"012345678901234567890123456789a",
				"0123456789012345678901234567890123456789",
				"012345678901234567890123456789",
				"1o23nm1kio2n3k12jn",
			},
			expected: []error{errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, nil, nil},
		},
		{
			desc:     "ShouldValidatePasswordAdvancedLowerUpperMin8",
			config:   schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, RequireLowercase: true, RequireUppercase: true}},
			have:     []string{"a", "b123", "1111111", "aaaaaaaa", "1o23nm1kio2n3k12jn", "ANJKJQ@#NEK!@#NJK!@#", "qjik2nkjAkjlmn123"},
			expected: []error{errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, errPasswordPolicyNoMet, nil},
		},
		{
			desc:   "ShouldValidatePasswordAdvancedAllMax100Min8",
			config: schema.PasswordPolicyConfiguration{Standard: schema.PasswordPolicyStandardParams{Enabled: true, MinLength: 8, MaxLength: 100, RequireLowercase: true, RequireUppercase: true, RequireNumber: true, RequireSpecial: true}},
			have: []string{
				"a",
				"b123",
				"1111111",
				"aaaaaaaa",
				"1o23nm1kio2n3k12jn",
				"ANJKJQ@#NEK!@#NJK!@#",
				"qjik2nkjAkjlmn123",
				"qjik2n@jAkjlmn123",
				"qjik2n@jAkjlmn123qjik2n@jAkjlmn123qjik2n@jAkjlmn123qjik2n@jAkjlmn123qjik2n@jAkjlmn123qjik2n@jAkjlmn123",
			},
			expected: []error{
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				errPasswordPolicyNoMet,
				nil,
				errPasswordPolicyNoMet,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			require.Equal(t, len(tc.have), len(tc.expected))
			for i := 0; i < len(tc.have); i++ {
				provider := NewPasswordPolicyProvider(tc.config)
				t.Run(tc.have[i], func(t *testing.T) {
					assert.Equal(t, tc.expected[i], provider.Check(tc.have[i]))
				})
			}
		})
	}
}
