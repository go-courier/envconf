package envconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvVar(t *testing.T) {
	require.Equal(t, EnvVarFromKeyValue("SU__TEST", "V"), &EnvVar{
		KeyPath:    "TEST",
		Value:      "V",
		ShouldConf: true,
		IsUpstream: true,
	})

	require.Equal(t, EnvVarFromKeyValue("S___TEST", "V"), &EnvVar{
		KeyPath:    "TEST",
		Value:      "V",
		ShouldConf: false,
	})

	require.Equal(t, "SU__TEST", (&EnvVar{
		KeyPath:    "TEST",
		Value:      "V",
		ShouldConf: true,
		IsUpstream: true,
	}).Key("S"))

	require.Equal(t, "S___TEST", (&EnvVar{
		KeyPath:    "TEST",
		Value:      "V",
		ShouldConf: false,
	}).Key("S"))
}
