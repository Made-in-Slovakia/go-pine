package pine

import (
	"fmt"
)

// Collection of types that answer content can be.
type answerType interface {
	uint8 | uint16 | uint32 | uint64 | string
}

// Answer is a struct that represents single and independent response to one specific command.
type Answer struct {
	// Content of answer defined as type any.
	content any
}

func (a *Answer) ContentAsUint8() (uint8, error) {
	return cast[uint8](a.content)
}

func (a *Answer) ContentAsUint16() (uint16, error) {
	return cast[uint16](a.content)
}

func (a *Answer) ContentAsUint32() (uint32, error) {
	return cast[uint32](a.content)
}

func (a *Answer) ContentAsUint64() (uint64, error) {
	return cast[uint64](a.content)
}

func (a *Answer) ContentAsString() (string, error) {
	return cast[string](a.content)
}

// Returns input casted to type defined by generics. Intended only for internal use of Answer struct.
func cast[T answerType](input any) (T, error) {
	v, ok := input.(T)
	if !ok {
		return v, fmt.Errorf("answer is not of required type %T, it is %T type", *new(T), input)
	}

	return v, nil
}
