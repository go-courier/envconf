package envconf

import (
	"fmt"
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
	Slice     []string
	PtrString *string `env:""`
	Host      string  `env:",upstream"`
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
	c.defaultValue = true
	c.Config.defaultValue = true

	envVars := NewEnvVars("S")

	t.Run("Decoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).Encode(&c)

		require.NotEqual(t, `
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

		require.Equal(t, `
SU__HOST=
S__CONFIG_KEY=key
S__CONFIG_PASSWORD=
S__KEY=123456
S__PASSWORD=123123
S__PTRSTRING=123456
S___BOOL=false
S___CONFIG_BOOL=false
S___CONFIG_DURATION=0s
S___DURATION=10s
S___SLICE_0=1
S___SLICE_1=2
`, "\n"+string(data))

		envVars := EnvVarsFromEnviron("S", strings.Split(string(data), "\n"))

		fmt.Println(envVars)

		c2 := Config{}
		err := NewDotEnvDecoder(envVars).Decode(&c2)
		require.NoError(t, err)

		require.Equal(t, c, c2)
	})

	t.Run("Encoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).SecurityEncode(&c)

		require.NotEqual(t, `
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
	})
	t.Run("Encoding", func(t *testing.T) {
		data, _ := NewDotEnvEncoder(envVars).SecurityEncode(&c)

		require.Equal(t, `
SU__HOST=
S__CONFIG_KEY=key
S__CONFIG_PASSWORD=
S__KEY=123456
S__PASSWORD=******
S__PTRSTRING=123456
S___BOOL=false
S___CONFIG_BOOL=false
S___CONFIG_DURATION=0s
S___DURATION=10s
S___SLICE_0=1
S___SLICE_1=2
`, "\n"+string(data))
	})

}
