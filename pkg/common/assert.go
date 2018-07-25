package common

import log "github.com/Sirupsen/logrus"

// Assertf checks cond and panic with passed error message if cond is not fulfilled
func Assertf(cond bool, format string, v ...interface{}) {
	if !cond {
		log.Panicf(format, v)
	}
}

// Assert checks cond and panic if cond is not fulfilled
func Assert(cond bool) {
	if !cond {
		panic("assert failed")
	}
}
