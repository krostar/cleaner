package cleaner

import (
	"io"
	"os"
	"runtime/debug"

	"github.com/pkg/errors"
)

// nolint: gochecknoglobals
var stopers []func()

// OnFailure is the signature of the function called on failure.
type OnFailure func(err error, stack []byte)

// Add adds the provided stop to the list of functions
// that will be cleaned when Clean() is called.
func Add(stop func()) {
	stopers = append(stopers, stop)
}

// Clean calls each added functions, call them
// and handle recover.
func Clean(onFailure OnFailure) {
	if reason := recover(); reason != nil {
		onPanic(reason, debug.Stack(), onFailure)
	}
	defer func() {
		// catch panic again in case we generate a panic in onFailure handlers
		// but don't call onFailure this time
		if reason := recover(); reason != nil {
			onPanic(reason, debug.Stack(), func(err error, stack []byte) {
				io.WriteString(os.Stderr, errors.Wrap(err, "panic catched in panic handler").Error())
				os.Stderr.Write(stack)
			})
		}
	}()

	for _, stop := range stopers {
		if stop != nil {
			stop()
		}
	}
}

// Reset empties the stopers list.
func Reset() {
	stopers = nil
}

func onPanic(reason interface{}, stack []byte, onFailure OnFailure) {
	var err error

	switch r := reason.(type) {
	case error:
		err = r
	case string:
		err = errors.New(r)
	default:
		err = errors.Errorf("%v", r)
	}

	onFailure(err, stack)
}
