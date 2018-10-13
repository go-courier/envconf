package envconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPathWalker(t *testing.T) {
	tt := require.New(t)

	pw := NewPathWalker()
	pw.Enter("KEY")
	tt.Equal([]interface{}{"KEY"}, pw.Paths())
	tt.Equal("KEY", pw.String())

	pw.Enter(1)
	tt.Equal([]interface{}{"KEY", 1}, pw.Paths())
	tt.Equal("KEY_1", pw.String())

	pw.Enter("PROP")
	tt.Equal([]interface{}{"KEY", 1, "PROP"}, pw.Paths())
	tt.Equal("KEY_1_PROP", pw.String())

	pw.Exit()
	pw.Exit()
	tt.Equal([]interface{}{"KEY"}, pw.Paths())
	tt.Equal("KEY", pw.String())
}
