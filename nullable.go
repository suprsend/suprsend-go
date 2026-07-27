package suprsend

import "encoding/json"

// Nullable is a tri-state JSON field. It distinguishes between three cases
// that a plain pointer + omitempty cannot:
//
//  1. not set    => the key is omitted from the payload entirely
//  2. set to null => the key is emitted as `null`
//  3. set to value => the key is emitted with the value
//
// Use it with the `omitzero` json tag (Go 1.24+) so that an unset field is
// dropped from the marshalled output:
//
//	DigestSchedule Nullable[any] `json:"digest_schedule,omitzero"`
//
// Construct with NewValue / NewNull; the zero value is the "not set" state.
type Nullable[T any] struct {
	set   bool // whether the key should appear in the payload
	value *T   // nil => JSON null, non-nil => the value
}

// NewValue returns a Nullable that marshals to the given value.
func NewValue[T any](v T) Nullable[T] { return Nullable[T]{set: true, value: &v} }

// NewNull returns a Nullable that marshals to an explicit JSON null.
func NewNull[T any]() Nullable[T] { return Nullable[T]{set: true} }

// IsSet reports whether the field will be included in the payload.
func (n Nullable[T]) IsSet() bool { return n.set }

// Get returns the value and whether a (non-null) value is present.
func (n Nullable[T]) Get() (T, bool) {
	if n.value == nil {
		var zero T
		return zero, false
	}
	return *n.value, true
}

// IsZero drives the `omitzero` json tag: an unset field is omitted.
func (n Nullable[T]) IsZero() bool { return !n.set }

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(n.value)
}

func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.set = true
	if string(b) == "null" {
		n.value = nil
		return nil
	}
	var v T
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	n.value = &v
	return nil
}
