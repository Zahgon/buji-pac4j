package io

import (
	"testing"
)

// stubSerializable is a minimal Serializable used to exercise the serializer
// round-trip without importing higher-level packages.
type stubSerializable struct {
	payload []byte
}

func (s *stubSerializable) MarshalBinary() ([]byte, error) {
	out := make([]byte, len(s.payload))
	copy(out, s.payload)
	return out, nil
}

func (s *stubSerializable) UnmarshalBinary(data []byte) error {
	s.payload = make([]byte, len(data))
	copy(s.payload, data)
	return nil
}

func TestDefaultSerializerRoundTrip(t *testing.T) {
	ser := NewDefaultSerializer()
	src := &stubSerializable{payload: []byte("hello-world")}

	data, err := ser.Serialize(src)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}
	if string(data) != "hello-world" {
		t.Fatalf("unexpected serialized bytes: %q", string(data))
	}

	dst := &stubSerializable{}
	if err := ser.Deserialize(data, dst); err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}
	if string(dst.payload) != "hello-world" {
		t.Fatalf("round-trip mismatch: %q", string(dst.payload))
	}
}
