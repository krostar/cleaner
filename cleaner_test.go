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

func TestClean_panic(t *testing.T) {
	t.Run("recover with string", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { require.Equal(t, "hello", err.Error()) })
		Reset()
		panic("hello")
	})

	t.Run("recover with error", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { require.Equal(t, "hello", err.Error()) })
		Reset()
		panic(errors.New("hello"))
	})

	t.Run("recover with unknown", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { require.Error(t, err) })
		Reset()
		panic(42)
	})
}

func TestClean_panic_in_handler(t *testing.T) {
	t.Run("successful clean", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) {})
		Reset()
		Add(func() {})
	})

	t.Run("recover with string", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { t.FailNow() })
		Reset()
		Add(func() { panic("hello") })
	})

	t.Run("recover with error", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { t.FailNow() })
		Reset()
		Add(func() { panic(errors.New("hello")) })
	})

	t.Run("recover with unknown", func(t *testing.T) {
		defer Clean(func(err error, stack []byte) { t.FailNow() })
		Reset()
		Add(func() { panic(42) })
	})
}
