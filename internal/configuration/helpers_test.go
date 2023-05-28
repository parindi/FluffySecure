package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSecretKey(t *testing.T) {
	assert.True(t, IsSecretKey("my_fake_token"))
	assert.False(t, IsSecretKey("my_fake_tokenz"))
	assert.True(t, IsSecretKey("my_.fake.secret"))
	assert.True(t, IsSecretKey("my.password"))
	assert.False(t, IsSecretKey("my.passwords"))
	assert.False(t, IsSecretKey("my.passwords"))
}

func TestGetEnvConfigMaps(t *testing.T) {
	var (
		key string
		ok  bool
	)

	input := []string{
		"my.non_secret.config_item",
		"myother.configkey",
		"mysecret.password",
		"mysecret.user_password",
	}

	keys, ignoredKeys := getEnvConfigMap(input, DefaultEnvPrefix, DefaultEnvDelimiter)

	key, ok = keys[DefaultEnvPrefix+"MY_NON_SECRET_CONFIG_ITEM"]
	assert.True(t, ok)
	assert.Equal(t, key, "my.non_secret.config_item")

	key, ok = keys[DefaultEnvPrefix+"MYSECRET_USER_PASSWORD"]
	assert.True(t, ok)
	assert.Equal(t, key, "mysecret.user_password")

	key, ok = keys[DefaultEnvPrefix+"MYOTHER_CONFIGKEY"]
	assert.False(t, ok)
	assert.Equal(t, key, "")

	key, ok = keys[DefaultEnvPrefix+"MYSECRET_PASSWORD"]
	assert.False(t, ok)
	assert.Equal(t, key, "")

	assert.Len(t, ignoredKeys, 3)
	assert.Contains(t, ignoredKeys, DefaultEnvPrefix+"MYOTHER_CONFIGKEY_FILE")
	assert.Contains(t, ignoredKeys, DefaultEnvPrefix+"MYSECRET_PASSWORD_FILE")
	assert.Contains(t, ignoredKeys, DefaultEnvPrefix+"MYSECRET_USER_PASSWORD_FILE")
}

func TestGetSecretConfigMap(t *testing.T) {
	var (
		key string
		ok  bool
	)

	input := []string{
		"my.non_secret.config_item",
		"myother.configkey",
		"mysecret.password",
		"mysecret.user_password",
	}

	keys := getSecretConfigMap(input, DefaultEnvPrefix, DefaultEnvDelimiter)

	key, ok = keys[DefaultEnvPrefix+"MY_NON_SECRET_CONFIG_ITEM_FILE"]
	assert.False(t, ok)
	assert.Equal(t, key, "")

	key, ok = keys[DefaultEnvPrefix+"MYOTHER_CONFIGKEY_FILE"]
	assert.True(t, ok)
	assert.Equal(t, key, "myother.configkey")

	key, ok = keys[DefaultEnvPrefix+"MYSECRET_PASSWORD_FILE"]
	assert.True(t, ok)
	assert.Equal(t, key, "mysecret.password")

	key, ok = keys[DefaultEnvPrefix+"MYSECRET_USER_PASSWORD_FILE"]
	assert.True(t, ok)
	assert.Equal(t, key, "mysecret.user_password")
}
