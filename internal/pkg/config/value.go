package config

import (
	"fmt"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

// Value описывает значение конфига.
type Value interface {
	IsNil() bool
	IsEqual(Value) bool

	Bool() bool
	MaybeBool() (bool, error)

	Int64() int64
	MaybeInt64() (int64, error)

	Float64() float64
	MaybeFloat64() (float64, error)

	Duration() time.Duration
	MaybeDuration() (time.Duration, error)

	String() string
	MaybeString() (string, error)
}

// concreteValue реализация Value
type concreteValue struct {
	Type  string `json:"type"`
	Value any    `json:"value"`

	asStringPtr atomic.Pointer[string]
}

func (v *concreteValue) Int64() int64 {
	val, _ := v.MaybeInt64()
	return val
}

func (v *concreteValue) MaybeInt64() (int64, error) {
	raw := v.String()

	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value cannot be parsed as int64: %v", err)
	}
	return val, nil
}

func (v *concreteValue) Float64() float64 {
	val, _ := v.MaybeFloat64()
	return val
}

func (v *concreteValue) MaybeFloat64() (float64, error) {
	raw := v.String()

	val, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("value cannot be parsed as float64: %v", err)
	}
	return val, nil
}

func (v *concreteValue) Duration() time.Duration {
	val, _ := v.MaybeDuration()
	return val
}

func (v *concreteValue) MaybeDuration() (time.Duration, error) {
	raw := v.String()

	val, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("value cannot be parsed as time.Duration: %v", err)
	}
	return val, nil
}

func (v *concreteValue) String() string {
	raw, _ := v.MaybeString()
	return raw
}

func (v *concreteValue) MaybeString() (string, error) {
	if v.IsNil() {
		return "", nil
	}
	if sPtr := v.asStringPtr.Load(); sPtr != nil && *sPtr != "" {
		return *sPtr, nil
	}

	var s string
	s = fmt.Sprintf("%v", v.Value)

	v.asStringPtr.Store(&s)
	return s, nil
}

func (v *concreteValue) IsNil() bool {
	return v == nil || v.Value == nil
}

func (v *concreteValue) IsEqual(other Value) bool {
	if v.IsNil() && other.IsNil() {
		return true
	}
	return v.String() == other.String()
}

func (v *concreteValue) Bool() bool {
	val, _ := v.MaybeBool()
	return val
}

func (v *concreteValue) MaybeBool() (bool, error) {
	if v.IsNil() {
		return false, errors.New("value is nil")
	}
	switch val := v.Value.(type) {
	case bool:
		return val, nil
	case string:
		return val == "true", nil
	default:
		return false, errors.Errorf("cannot convert %T to bool", v.Value)
	}
}
