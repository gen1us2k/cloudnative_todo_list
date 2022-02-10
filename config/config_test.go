package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	os.Setenv("SUPABASE_URL", "http://localhost:3000")
	c, err := Parse()
	assert.NoError(t, err)
	assert.Equal(t, EnvDevelopment, c.Env)
	assert.Equal(t, "http://127.0.0.1:4433", c.KratosAPIURL)
	assert.Equal(t, "http://127.0.0.1:4455", c.KratosUIURL)
	assert.Equal(t, 8080, c.GRPCPort)
}
