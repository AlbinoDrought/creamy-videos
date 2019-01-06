package streamers

import "io"

type byteTransformerFunc = func(p []byte, n int)

type byteTransformer struct {
	io.Reader
	reader      io.Reader
	transformer byteTransformerFunc
}

func makeByteTransformer(base io.Reader, transformer byteTransformerFunc) io.Reader {
	return byteTransformer{
		reader:      base,
		transformer: transformer,
	}
}

func (a byteTransformer) Read(p []byte) (int, error) {
	n, err := a.reader.Read(p)

	if err != nil {
		return n, err
	}

	a.transformer(p, n)

	return n, nil
}
