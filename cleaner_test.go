package cleaner

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	Reset()
	Add(func() {})
	Add(func() {})

	require.Len(t, stopers, 2)
}

func TestClean(t *testing.T) {
	t.Run("successful clean", func(t *testing.T) {
		defer Clean(func(err error) {})
		Reset()
		Add(func() {})
	})

	t.Run("recover with string", func(t *testing.T) {
		defer Clean(func(err error) { require.Equal(t, "hello", err.Error()) })
		Reset()
		Add(func() { panic("hello") })
	})

	t.Run("recover with error", func(t *testing.T) {
		defer Clean(func(err error) { require.Equal(t, "hello", err.Error()) })
		Reset()
		Add(func() { panic(errors.New("hello")) })
	})

	t.Run("recover with unknown", func(t *testing.T) {
		defer Clean(func(err error) { require.Error(t, err) })
		Reset()
		Add(func() { panic(42) })
	})
}
