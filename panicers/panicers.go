package panicers

import "github.com/pkg/errors"

func PanicOn(cond bool, err error) {
	if cond {
		panic(err)
	}
}

func WrapAndPanicIfErr(err error, format string, args ...interface{}) {
	if err != nil {
		if len(args) > 0 {
			panic(errors.Wrapf(err, format, args))
		} else {
			panic(err)
		}
	}
}
