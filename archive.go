package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extract(content []byte, length int64, dest string, dir string) error {

	zr, err := zip.NewReader(bytes.NewReader(content), length)
	if err != nil {
		return err
	}

	path := filepath.Join(wd, dir, dest)

	err = mkdir(path)
	if err != nil {

	}

	// Used to avoid open file descriptors. TODO: Check.
	writeFile := func(zf *zip.File) {

		// If this is a Directory, create it and move on.
		if zf.FileInfo().IsDir() {

			folder := filepath.Join(wd, dir, zf.Name)

			err := mkdir(folder)
			if err != nil {
				// ignore errors
			}

			return

		}

		fr, err := zf.Open()
		if err != nil {
			fmt.Printf("Unable to read file: %s\n", zf.Name)
			return
		}
		defer fr.Close()

		path := strings.Replace(filepath.Join(wd, dir, zf.Name), "/", string(filepath.Separator), -1)
		//dir, _ := filepath.Split(path)
		dt := filepath.Dir(path)

		// Make the directory required by this File.
		err = mkdir(dt)
		if err != nil {
			// ignore error
		}

		// Create File.
		f, err := os.Create(path)
		if err != nil {
			fmt.Printf("Unable to create file: %s\n", path)
			return
		}
		defer f.Close()

		// Copy contents to File.
		_, err = io.Copy(f, fr)
		if err != nil {
			fmt.Printf("Problem writing contents to file <%s>: %s\n", path, err)
			return
		}

		return

	}

	// Create each File in the Archive.
	for _, zf := range zr.File {
		writeFile(zf)
	}

	return nil
}

func mkdir(dirPath string) error {

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("%s: making directory: %v", dirPath, err)
	}

	return nil

}
