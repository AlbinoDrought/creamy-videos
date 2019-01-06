package streamers

import "io"

// XorifyReader xors all read bytes by the specified key
func XorifyReader(reader io.Reader, key byte) io.Reader {
	return makeByteTransformer(
		reader,
		func(p []byte, n int) {
			for i := 0; i < n; i++ {
				p[i] = p[i] ^ key
			}
		},
	)
}
