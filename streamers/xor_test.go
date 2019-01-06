package streamers

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXorifyReader(t *testing.T) {
	input := []byte{
		0x0d,
		0x0c,
		0x08,
		0x0d,
	}
	expected := []byte("dead")

	reader := bytes.NewReader(input)
	xorReader := XorifyReader(reader, 0x69)

	buf := new(bytes.Buffer)
	buf.ReadFrom(xorReader)

	assert.Equal(t, expected, buf.Bytes())
}

func TestXorifyReaderDouble(t *testing.T) {
	input := []byte{
		0x0d,
		0x0c,
		0x08,
		0x0d,
	}

	reader := bytes.NewReader(input)
	xorReader := XorifyReader(reader, 0x69)
	xorReader = XorifyReader(xorReader, 0x69)

	buf := new(bytes.Buffer)
	buf.ReadFrom(xorReader)

	assert.Equal(t, input, buf.Bytes())
}
