package xvalidator

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
)

type NoStringMethod int

const (
	NoStringMethodError   NoStringMethod = -100
	NoStringMethodNone    NoStringMethod = 0
	NoStringMethodSuccess NoStringMethod = 100
)

type nsm struct {
	Status NoStringMethod `validate:"is_enum"`
}

func TestEnum_NoStringMethod(t *testing.T) {
	value := nsm{Status: NoStringMethodError}
	require.Error(t, Validate.Struct(value), "no string method cause error")
}

type StringMethodWithoutReturning int

func (s StringMethodWithoutReturning) String() {}

const (
	StringMethodWithoutReturningError   StringMethodWithoutReturning = -100
	StringMethodWithoutReturningNone    StringMethodWithoutReturning = 0
	StringMethodWithoutReturningSuccess StringMethodWithoutReturning = 100
)

type smwr struct {
	Status StringMethodWithoutReturning `validate:"is_enum"`
}

func TestEnum_StringMethodWithoutReturning(t *testing.T) {
	value := smwr{Status: StringMethodWithoutReturningError}
	require.Error(t, Validate.Struct(value), "no returning value string method cause error")
}

type Enumerable int

func (s Enumerable) String() string {
	switch s {
	case EnumerableError:
		return "error"
	case EnumerableNone:
		return "none"
	case EnumerableSuccess:
		return "success"
	default:
		return ""
	}
}

const (
	EnumerableError   Enumerable = -100
	EnumerableNone    Enumerable = 0
	EnumerableSuccess Enumerable = 100
)

type enumerable struct {
	Status Enumerable `validate:"is_enum"`
}

func TestEnum(t *testing.T) {
	value := enumerable{Status: EnumerableError}
	require.NoError(t, Validate.Struct(value))

	value.Status = Enumerable(math.MaxInt)
	require.Error(t, Validate.Struct(value))
}
