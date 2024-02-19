package xerrors

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
)

type ErrorWrapper struct {
	Err      error
	Filename string
	Line     int
	Params   []any
}

func (e ErrorWrapper) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%v:%v %v\n", e.Filename, e.Line, e.Err.Error()))

	for _, param := range e.Params {
		// XXX: dumper
		buffer.WriteString(fmt.Sprintf("%v", param))
		buffer.WriteString("\n")
	}

	return buffer.String()
}

type ErrorList []error

func (e ErrorList) Error() string {
	var buffer bytes.Buffer

	for _, err := range e {
		buffer.WriteString(err.Error())
	}

	return buffer.String()
}

func NewError(err error, args ...any) ErrorWrapper {
	_, filename, line, _ := runtime.Caller(1)

	filename = filepath.Base(filename)

	return ErrorWrapper{
		Err:      err,
		Filename: filename,
		Line:     line,
		Params:   append([]any{}, args...),
	}
}

func Errors(inputList ...any) ErrorList {
	_, filename, line, _ := runtime.Caller(1)

	filename = filepath.Base(filename)

	var outputList ErrorList

	for _, input := range inputList {
		switch typedInput := input.(type) {

		case ErrorList:

			outputList = append(outputList, typedInput...)

		case ErrorWrapper:

			outputList = append(outputList, typedInput)

		case error:

			outputList = append(outputList,
				ErrorWrapper{
					Err:      typedInput,
					Filename: filename,
					Line:     line,
				},
			)

		case string:

			outputList = append(outputList,
				ErrorWrapper{
					Err:      errors.New(typedInput),
					Filename: filename,
					Line:     line,
				},
			)

		}
	}

	return outputList
}
