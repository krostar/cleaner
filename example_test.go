package cleaner_test

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/krostar/cleaner"
)

func Example() {
	cleaner.Reset() // this is useful only in tests
	defer cleaner.Clean(onCleanFailure)

	var (
		something     = initSomething()
		somethingElse = initSomethingElse()
	)

	// do something with it
	_ = something
	_ = somethingElse

	// Output:
	// something successfully init
	// oops, there is a failure: oh no somethingelse could not be init
	// close something
}

func initSomething() *something {
	var s something
	cleaner.Add(s.Close)
	fmt.Println("something successfully init")
	return &s
}

func initSomethingElse() *somethingElse {
	var s somethingElse
	if true {
		panic("oh no somethingelse could not be init")
	}

	cleaner.Add(s.Flush)
	fmt.Println("somethingelse successfully init")
	return &s

}

func onCleanFailure(err error) {
	fmt.Println(errors.Wrap(err, "oops, there is a failure"))
}

type something struct{}

func (something) Close() {
	fmt.Println("close something")
}

type somethingElse struct{}

func (somethingElse) Flush() {
	fmt.Println("flush somethingelse")
}
