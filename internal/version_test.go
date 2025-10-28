package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	version := Version()
	assert.NotEmpty(t, version, "Version should not be empty")
	assert.Equal(t, "0.1.0", version, "Version should be 0.1.0")
}

func TestAppName(t *testing.T) {
	appName := AppName()
	assert.NotEmpty(t, appName, "AppName should not be empty")
	assert.Equal(t, "linux-wave", appName, "AppName should be linux-wave")
}