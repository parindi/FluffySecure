package validator

import (
	"fmt"

	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/utils"
)

// ValidateLog validates the logging configuration.
func ValidateLog(config *schema.Configuration, validator *schema.StructValidator) {
	if config.Log.Level == "" {
		config.Log.Level = schema.DefaultLoggingConfiguration.Level
	}

	if config.Log.Format == "" {
		config.Log.Format = schema.DefaultLoggingConfiguration.Format
	}

	if !utils.IsStringInSlice(config.Log.Level, validLogLevels) {
		validator.Push(fmt.Errorf(errFmtLoggingLevelInvalid, strJoinOr(validLogLevels), config.Log.Level))
	}
}
