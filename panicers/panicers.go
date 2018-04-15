package panicers

import "github.com/pkg/errors"

// Remember kids only panic on unrecoverable errors

// PanicOn if cond is true will panic with err
func PanicOn(cond bool, err error) {
	if cond {
		panic(err)
	}
}

// WrapAndPanicIfErr if err is not nil will wrap err and panic
func WrapAndPanicIfErr(err error, format string, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			panic(errors.Wrapf(err, format, args))
		} else {
			panic(err)
		}
	}
}
