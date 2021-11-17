package cache

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type tlvWriter struct {
	writer io.Writer
	err    error
}

// WriteHeader will write the given header and version number.
func (tlv *tlvWriter) WriteHeader(header []byte, version uint32) {
	if tlv.err != nil {
		return
	}

	_, err := tlv.writer.Write(header)
	if err != nil {
		tlv.err = err
		return
	}

	err = binary.Write(tlv.writer, binary.LittleEndian, version)
	if err != nil {
		tlv.err = err
	}
}

// WriteTLV will write the given type and value to the underlying writer.
func (tlv *tlvWriter) WriteTLV(t uint32, value []byte) {
	if tlv.err != nil {
		return
	}

	// Write the type.
	err := binary.Write(tlv.writer, binary.LittleEndian, t)
	if err != nil {
		tlv.err = err
		return
	}

	// Write the length.
	err = binary.Write(tlv.writer, binary.LittleEndian, uint64(len(value)))
	if err != nil {
		tlv.err = err
		return
	}

	// Write the value.
	_, err = tlv.writer.Write(value)
	tlv.err = err
}

// End will finish writing the TLV stream, and will return an error if anything failed.
func (tlv *tlvWriter) End() error {
	return tlv.err
}

type tlvReader struct {
	reader io.Reader
}

// ReadHEader will read the header from the underlying reader.
func (tlv *tlvReader) ReadHeader(header []byte) (version uint32, err error) {
	// Read the header.
	headerBytes := make([]byte, len(header))
	_, err = io.ReadFull(tlv.reader, headerBytes)
	if err != nil {
		return 0, err
	}
	if !bytes.Equal(headerBytes, []byte(header)) {
		return 0, fmt.Errorf("invalid cache header")
	}

	// Read the version.
	err = binary.Read(tlv.reader, binary.LittleEndian, &version)
	if err != nil {
		return 0, err
	}

	return version, nil
}

// ReadTLV reads a TLV from the underlying reader.
func (tlv *tlvReader) ReadTLV() (t uint32, value []byte, err error) {
	// Read the type.
	err = binary.Read(tlv.reader, binary.LittleEndian, &t)
	if err != nil {
		return 0, nil, err
	}

	// Read the length.
	var size uint64
	err = binary.Read(tlv.reader, binary.LittleEndian, &size)
	if err != nil {
		return 0, nil, err
	}

	// Read the value.
	value = make([]byte, size)
	_, err = io.ReadFull(tlv.reader, value)
	if err != nil {
		return 0, nil, err
	}

	return t, value, nil
}

// ReadTLVOfType will read a TLV from the underlying reader, and error if the
// type is not the expected type.
func (tlv *tlvReader) ReadTLVOfType(t uint32) ([]byte, error) {
	readType, value, err := tlv.ReadTLV()
	if err != nil {
		return nil, err
	}
	if readType != t {
		return nil, fmt.Errorf("expected type %v, found %v", t, readType)
	}

	return value, nil
}

// End will finish reading the underlying file, and error if there is any data left over.
func (tlv *tlvReader) End() error {
	data := make([]byte, 1)
	n, err := tlv.reader.Read(data)
	if err != io.EOF && err != nil {
		return err
	}
	if n > 0 {
		return fmt.Errorf("unexpected data at end of record")
	}
	return nil
}
