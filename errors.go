package gserrors

import (
	"bytes"
	"fmt"
	"runtime"
)

//GSError gsdocker error interface
type GSError interface {
	error          //inher from system error interface
	Stack() string //get stack trace message
	Origin() error //get origin error object
}

type errorHost struct {
	origin  error  //origin error
	stack   string //stack trace message
	message string //error message
}

func (err *errorHost) Error() string {
	if err.message == "" {
		return fmt.Sprintf("%s\n%s", err.origin.Error(), err.stack)
	}

	return fmt.Sprintf("%s\n%s", err.message, err.stack)
}

func (err *errorHost) Stack() string {
	return err.stack
}

func (err *errorHost) Origin() error {
	return err.origin
}

func stack() []byte {
	var buff bytes.Buffer
	for skip := 2; ; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		buff.WriteString(fmt.Sprintf("\tfile = %s, line = %d\n", file, line))
	}

	return buff.Bytes()
}

//Panic throw GSError
func Panic(err error) {
	panic(New(err))
}

//Panicf throw GSError
func Panicf(err error, fmtstring string, args ...interface{}) {
	panic(Newf(err, fmtstring, args...))
}

//New create new GSError object
func New(err error) GSError {
	return &errorHost{
		origin: err,
		stack:  string(stack()),
	}
}

//Newf create new GSError object
func Newf(err error, fmtstring string, args ...interface{}) GSError {
	return &errorHost{
		origin:  err,
		stack:   string(stack()),
		message: fmt.Sprintf(fmtstring, args...),
	}
}
