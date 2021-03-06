package sarama

import (
	"bytes"
	"reflect"
	"testing"
)

type testRequestBody struct {
}

func (s *testRequestBody) key() int16 {
	return 0x666
}

func (s *testRequestBody) version() int16 {
	return 0xD2
}

func (s *testRequestBody) encode(pe packetEncoder) error {
	return pe.putString("abc")
}

// not specific to request tests, just helper functions for testing structures that
// implement the encoder or decoder interfaces that needed somewhere to live

func testEncodable(t *testing.T, name string, in encoder, expect []byte) {
	packet, err := encode(in)
	if err != nil {
		t.Error(err)
	} else if !bytes.Equal(packet, expect) {
		t.Error("Encoding", name, "failed\ngot ", packet, "\nwant", expect)
	}
}

func testDecodable(t *testing.T, name string, out decoder, in []byte) {
	err := decode(in, out)
	if err != nil {
		t.Error("Decoding", name, "failed:", err)
	}
}

func testVersionDecodable(t *testing.T, name string, out versionedDecoder, in []byte, version int16) {
	err := versionedDecode(in, out, version)
	if err != nil {
		t.Error("Decoding", name, "version", version, "failed:", err)
	}
}

func testRequest(t *testing.T, name string, rb protocolBody, expected []byte) {
	// Encoder request
	req := &request{correlationID: 123, clientID: "foo", body: rb}
	packet, err := encode(req)
	headerSize := 14 + len("foo")
	if err != nil {
		t.Error(err)
	} else if !bytes.Equal(packet[headerSize:], expected) {
		t.Error("Encoding", name, "failed\ngot ", packet[headerSize:], "\nwant", expected)
	}
	// Decoder request
	decoded, n, err := decodeRequest(bytes.NewReader(packet))
	if err != nil {
		t.Error("Failed to decode request", err)
	} else if decoded.correlationID != 123 || decoded.clientID != "foo" {
		t.Errorf("Decoded header is not valid: %v", decoded)
	} else if !reflect.DeepEqual(rb, decoded.body) {
		t.Errorf("Decoded request does not match the encoded one\nencoded: %v\ndecoded: %v", rb, decoded.body)
	} else if n != len(packet) {
		t.Errorf("Decoded request bytes: %d does not match the encoded one: %d\n", n, len(packet))
	}
}

func testResponse(t *testing.T, name string, res protocolBody, expected []byte) {
	encoded, err := encode(res)
	if err != nil {
		t.Error(err)
	} else if expected != nil && !bytes.Equal(encoded, expected) {
		t.Error("Encoding", name, "failed\ngot ", encoded, "\nwant", expected)
	}

	decoded := reflect.New(reflect.TypeOf(res).Elem()).Interface().(versionedDecoder)
	if err := versionedDecode(encoded, decoded, res.version()); err != nil {
		t.Error("Decoding", name, "failed:", err)
	}

	if !reflect.DeepEqual(decoded, res) {
		t.Errorf("Decoded response does not match the encoded one\nencoded: %#v\ndecoded: %#v", res, decoded)
	}
}
