package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/AlbinoDrought/creamy-videos/files"
	"github.com/AlbinoDrought/creamy-videos/streamers"
)

func main() {
	inputDir := "./../data"
	outputDir := "./../data-xor"
	transformed := files.TransformFileSystem(
		http.Dir(inputDir),
		func(file http.File) io.Reader {
			return streamers.XorifyReader(file, 0x69)
		},
	)

	files, err := ioutil.ReadDir(inputDir)

	if err != nil {
		panic(err)
	}

	fmt.Printf("files: %+v", files)

	for _, f := range files {
		// fullInputPath := path.Join(inputDir, f.Name())
		fullInputPath := f.Name()

		fullOutputPath := path.Join(outputDir, f.Name())

		createdFile, err := os.Create(fullOutputPath)
		if err != nil {
			panic(err)
		}

		defer createdFile.Close()

		openedFile, err := transformed.Open(fullInputPath)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(createdFile, openedFile)

		if err != nil {
			panic(err)
		}
	}
}
