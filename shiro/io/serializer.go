// Package io ports the subset of org.apache.shiro.io used by the bridge: a
// binary serializer used to round-trip principals (as exercised by the principal
// serialization test).
package io

import "encoding"

// Serializable is satisfied by any value that can encode and decode itself to a
// byte slice, mirroring the java.io.Serializable contract exercised by Shiro's
// DefaultSerializer for the types the bridge serializes.
type Serializable interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// DefaultSerializer ports org.apache.shiro.io.DefaultSerializer: it converts a
// serializable object to bytes and back. Java's implementation uses native JVM
// serialization; here each serializable type provides its own binary form via
// MarshalBinary/UnmarshalBinary, preserving value equality across the round-trip.
type DefaultSerializer struct{}

// NewDefaultSerializer constructs a serializer.
func NewDefaultSerializer() *DefaultSerializer {
	return &DefaultSerializer{}
}

// Serialize encodes a serializable value to bytes
// (DefaultSerializer.serialize(Object)).
func (DefaultSerializer) Serialize(v Serializable) ([]byte, error) {
	return v.MarshalBinary()
}

// Deserialize decodes bytes into the provided serializable target
// (DefaultSerializer.deserialize(byte[])). The caller supplies a target of the
// correct concrete type into which the bytes are decoded.
func (DefaultSerializer) Deserialize(data []byte, target Serializable) error {
	return target.UnmarshalBinary(data)
}
