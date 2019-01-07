package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
)

func main() {
	inputDir := "./../data"

	files, err := ioutil.ReadDir(inputDir)

	if err != nil {
		panic(err)
	}

	fmt.Printf("files: %+v", files)

	for _, f := range files {
		fullInputPath := path.Join(inputDir, f.Name())
		fullOutputPath := path.Join(inputDir, f.Name()+".jpg")

		cmd := exec.Command("ffmpeg", "-i", fullInputPath, "-vframes", "1", fullOutputPath)

		err = cmd.Run()

		if err != nil {
			panic(err)
		}
	}
}
