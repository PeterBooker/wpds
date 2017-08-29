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

func unzip(content []byte, dest string) error {

	br := bytes.NewReader(content)

	zr, err := zip.NewReader(br, int64(br.Len()))
	if err != nil {
		return err
	}

	path := filepath.Join(wd, "plugins", dest)

	err = mkdir(path)
	if err != nil {

	}

	// Used to avoid open file descriptors. TODO: Check.
	writeFile := func(zf *zip.File) {

		// If this is a Directory, create it and move on.
		if zf.FileInfo().IsDir() {

			folder := filepath.Join(path, zf.Name)

			err := mkdir(folder)
			if err != nil {

			}

			return

		}

		fr, err := zf.Open()
		if err != nil {
			fmt.Printf("Unable to read file: %s\n", zf.Name)
			return
		}
		defer fr.Close()

		path := strings.Replace(filepath.Join(wd, "plugins", zf.Name), "/", string(filepath.Separator), -1)
		dir, _ := filepath.Split(path)

		// Make the directory required by this File.
		err = mkdir(dir)
		if err != nil {

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
			fmt.Printf("Issue writing file <%s>: %s\n", path, err)
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
