package envconf

import (
	"strings"
	"testing"
	"time"

	"github.com/go-courier/ptr"
	"github.com/stretchr/testify/require"
)

func TestEnvVar(t *testing.T) {
	type SubConfig struct {
		Duration Duration
		Password Password `env:""`
		Key      string   `env:""`
		Bool     bool
		Map      map[string]string
		Func     func() error
		ignore   bool
	}

	type Config struct {
		Map       map[string]string
		Slice     []string
		PtrString *string `env:""`
		Host      string  `env:",upstream"`
		SubConfig
		Config SubConfig
	}

	c := Config{}

	c.Duration = Duration(10 * time.Second)
	c.Password = Password("123123")
	c.Key = "123456"
	c.PtrString = ptr.String("123456")
	c.Slice = []string{"1", "2"}
	c.Config.Key = "key"

	envVars := NewEnvVars("S")

	{
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		require.Equal(t, `
SU__Host=
S__Config_Key=key
S__Config_Password=
S__Key=123456
S__Password=123123
S__PtrString=123456
S___Bool=false
S___Config_Bool=false
S___Config_Duration=0s
S___Duration=10s
S___Slice_0=1
S___Slice_1=2
`, "\n"+string(data))

		{
			envVars := EnvVarsFromEnviron(strings.Split(string(data), "\n"))

			c2 := Config{}
			err := NewDotEnvDecoder(envVars).Decode(&c2)
			require.NoError(t, err)

			require.Equal(t, c, c2)
		}
	}

	{
		data, _ := NewDotEnvEncoder(envVars).SecurityEncode(&c)

		require.Equal(t, `
SU__Host=
S__Config_Key=key
S__Config_Password=
S__Key=123456
S__Password=******
S__PtrString=123456
S___Bool=false
S___Config_Bool=false
S___Config_Duration=0s
S___Duration=10s
S___Slice_0=1
S___Slice_1=2
`, "\n"+string(data))
	}
}
