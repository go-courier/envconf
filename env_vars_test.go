package envconf

import (
	"strings"
	"testing"
	"time"

	"github.com/go-courier/ptr"
	"github.com/stretchr/testify/require"
)

type SubConfig struct {
	Duration     Duration
	Password     Password `env:""`
	Key          string   `env:""`
	Bool         bool
	Map          map[string]string
	Func         func() error
	ignore       bool
	defaultValue bool
}

func (c *SubConfig) SetDefaults() {
	c.defaultValue = true
}

type Config struct {
	Map       map[string]string
	Slice     []string `env:""`
	PtrString *string  `env:""`
	Host      string   `env:",upstream"`
	SubConfig
	Config SubConfig
}

func TestEnvVars(t *testing.T) {
	c := Config{}

	c.Duration = Duration(10 * time.Second)
	c.Password = Password("123123")
	c.Key = "123456"
	c.PtrString = ptr.String("123456")
	c.Slice = []string{"1", "2"}
	c.Config.Key = "key"
	c.Config.defaultValue = true
	c.defaultValue = true

	envVars := NewEnvVars("S")

	t.Run("Encoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		require.Equal(t, `
S__Bool=false
S__Config_Bool=false
S__Config_Duration=0s
S__Config_Key=key
S__Config_Password=
S__Duration=10s
S__Host=
S__Key=123456
S__Password=123123
S__PtrString=123456
S__Slice_0=1
S__Slice_1=2
`, "\n"+string(data))
	})

	t.Run("Decoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		require.Equal(t, `
S__Bool=false
S__Config_Bool=false
S__Config_Duration=0s
S__Config_Key=key
S__Config_Password=
S__Duration=10s
S__Host=
S__Key=123456
S__Password=123123
S__PtrString=123456
S__Slice_0=1
S__Slice_1=2
`, "\n"+string(data))

		envVars := EnvVarsFromEnviron("S", strings.Split(string(data), "\n"))

		c2 := Config{}
		err := NewDotEnvDecoder(envVars).Decode(&c2)
		require.NoError(t, err)

		require.Equal(t, c, c2)
	})
}
