package cleaner

import (
	"github.com/pkg/errors"
)

// nolint: gochecknoglobals
var stopers []func()

// Add adds the provided stop to the list of functions
// that will be cleaned when Clean() is called.
func Add(stop func()) {
	stopers = append(stopers, stop)
}

// Clean calls each added functions, call them
// and handle recover.
func Clean(onFailure func(err error)) {
	if reason := recover(); reason != nil {
		onPanic(reason, onFailure)
	}
	defer func() {
		if reason := recover(); reason != nil {
			onPanic(reason, onFailure)
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

func onPanic(reason interface{}, onFailure func(error)) {
	var err error

	switch r := reason.(type) {
	case error:
		err = r
	case string:
		err = errors.New(r)
	default:
		err = errors.Errorf("%v", r)
	}

	onFailure(err)
}
