package files

type ByteTransformer = func(original []byte)

// readableTransformedFile implements http.File but supports
// transforming the read stream
// If there's a better way to do this, please let me know
type readableTransformedFile struct {
	ReadableFile
	transformer ByteTransformer
}

// Read transforms the parent implementation of Read
func (tf readableTransformedFile) Read(p []byte) (n int, err error) {
	n, err = tf.ReadableFile.Read(p)
	tf.transformer(p)
	return n, err
}

type writeableTransformedFile struct {
	WriteableFile
	transformer ByteTransformer
}

func (wtf writeableTransformedFile) Write(b []byte) (n int, err error) {
	wtf.transformer(b)
	return wtf.WriteableFile.Write(b)
}

type transformedFileSystem struct {
	FileSystem
	transformer ByteTransformer
}

func TransformFileSystem(fs FileSystem, transformer ByteTransformer) FileSystem {
	return transformedFileSystem{
		fs,
		transformer,
	}
}

func (fs transformedFileSystem) Create(name string) (WriteableFile, error) {
	file, err := fs.FileSystem.Create(name)
	if err != nil {
		return file, err
	}

	return writeableTransformedFile{
		WriteableFile: file,
		transformer:   fs.transformer,
	}, nil
}

func (fs transformedFileSystem) Open(path string) (ReadableFile, error) {
	file, err := fs.FileSystem.Open(path)

	if err != nil {
		return file, err
	}

	return readableTransformedFile{
		ReadableFile: file,
		transformer:  fs.transformer,
	}, nil
}
