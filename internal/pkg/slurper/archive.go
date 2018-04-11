package slurper

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
)

// ExtractZip extracts the archive containing extension data.
func ExtractZip(content []byte, length int64, dest string, ctx *context.Context) (uint64, error) {

	zr, err := zip.NewReader(bytes.NewReader(content), length)
	if err != nil {
		return 0, err
	}

	path := filepath.Join(wd, ctx.ExtensionType, dest)

	if utils.DirExists(path) {
		err := utils.RemoveDir(path)
		if err != nil {
			log.Printf("Cannot delete extension folder: %s\n", path)
		}
	}

	err = utils.CreateDir(path)
	if err != nil {
		log.Printf("Cannot create extension folder: %s\n", path)
	}

	// Used to avoid open file descriptors.
	// TODO: Check it actually helps.
	writeFile := func(zf *zip.File) {

		// If this is a Directory, create it and move on.
		if zf.FileInfo().IsDir() {

			folder := filepath.Join(wd, ctx.ExtensionType, zf.Name)

			utils.CreateDir(folder)

			return

		}

		fr, err := zf.Open()
		if err != nil {
			log.Printf("Unable to read file: %s\n", zf.Name)
			return
		}
		defer fr.Close()

		path := strings.Replace(filepath.Join(wd, ctx.ExtensionType, zf.Name), "/", string(filepath.Separator), -1)
		dt := filepath.Dir(path)

		// Make the directory required by this File.
		utils.CreateDir(dt)

		// Create File.
		f, err := os.Create(path)
		if err != nil {
			log.Printf("Unable to create file: %s\n", path)
			return
		}

		defer f.Close()

		// Copy contents to the File.
		_, err = io.Copy(f, fr)
		if err != nil {
			log.Printf("Problem writing contents to file <%s>: %s\n", path, err)
			f.Close()
			return
		}

		err = f.Close()
		if err != nil {
			log.Printf("Problem writing contents to file <%s>: %s\n", path, err)
			return
		}

		return

	}

	var size uint64

	// Create each File in the Archive.
	for _, zf := range zr.File {
		writeFile(zf)
		size += zf.UncompressedSize64
		ctx.Stats.IncrementTotalFiles()
	}

	return size, nil

}
