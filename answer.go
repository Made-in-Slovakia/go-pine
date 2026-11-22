package pine

import (
	"fmt"
	"math"
)

// Collection of types that [Answer] content can be.
type answerContentType interface {
	uint8 | uint16 | uint32 | uint64 | string
}

// Answer is a struct that represents single and independent response to one specific [Command].
type Answer struct {
	// Content of answer, defined as type any
	content any
}

// Returns the content as an uint8 or error if content is incompatible type.
func (a *Answer) ContentAsUint8() (uint8, error) {
	return cast[uint8](a.content)
}

// Returns the content as an uint16 or error if content is incompatible type.
func (a *Answer) ContentAsUint16() (uint16, error) {
	return cast[uint16](a.content)
}

// Returns the content as an uint32 or error if content is incompatible type.
func (a *Answer) ContentAsUint32() (uint32, error) {
	return cast[uint32](a.content)
}

// Returns the content as an uint64 or error if content is incompatible type.
func (a *Answer) ContentAsUint64() (uint64, error) {
	return cast[uint64](a.content)
}

// Returns the content as a string or error if content is incompatible type.
func (a *Answer) ContentAsString() (string, error) {
	return cast[string](a.content)
}

// Returns the content as a float or error if content is incompatible type.
func (a *Answer) ContentAsFloat() (float32, error) {
	value, err := a.ContentAsUint32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(value), nil
}

// Returns input casted to type defined by generics. Intended only for internal use of [Answer] struct.
func cast[T answerContentType](input any) (T, error) {
	value, ok := input.(T)
	if !ok {
		return value, fmt.Errorf("answer is not of required type %T, it is %T type", *new(T), input)
	}

	return value, nil
}
