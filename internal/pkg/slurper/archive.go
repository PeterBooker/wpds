package slurper

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/peterbooker/wpds/internal/pkg/context"
	"github.com/peterbooker/wpds/internal/pkg/utils"
)

// ExtractZip extracts the archive containing extension data.
func ExtractZip(name string, ctx *context.Context) (uint64, error) {

	var err error

	fname := filepath.Join(ctx.WorkingDirectory, ctx.ExtensionType, name, "plugin.zip")
	file, err := os.Open(fname)
	if err != nil {
		return 0, err
	}
	defer utils.CheckClose(file, &err)

	fi, err := file.Stat()
	if err != nil {
		file.Close()
		return 0, err
	}

	zr, err := zip.NewReader(file, fi.Size())
	if err != nil {
		file.Close()
		return 0, err
	}

	// Used to avoid open file descriptors.
	// TODO: Check it actually helps.
	writeFile := func(zf *zip.File) error {

		// If this is a Directory, create it and move on.
		if zf.FileInfo().IsDir() {
			folder := filepath.Join(wd, ctx.ExtensionType, zf.Name)
			utils.CreateDir(folder)
			return nil
		}

		fr, err := zf.Open()
		if err != nil {
			log.Printf("Unable to read file: %s\n", zf.Name)
			return err
		}
		defer utils.CheckClose(fr, &err)

		path := filepath.FromSlash(filepath.Join(wd, ctx.ExtensionType, zf.Name))
		dt := filepath.Dir(path)

		// Make the directory required by this File.
		utils.CreateDir(dt)

		// Create File.
		f, err := os.Create(path)
		if err != nil {
			log.Printf("Unable to create file: %s\n", path)
			return err
		}
		defer utils.CheckClose(f, &err)

		// Copy contents to the File.
		_, err = io.Copy(f, fr)
		if err != nil {
			log.Printf("Problem writing contents to file <%s>: %s\n", path, err)
			return err
		}

		return err

	}

	var size uint64

	// Create each File in the Archive.
	for _, zf := range zr.File {
		err := writeFile(zf)
		if err != nil {
			log.Printf("Error writing file: %s\n", err)
		}
		size += zf.UncompressedSize64
		ctx.Stats.IncrementTotalFiles()
	}

	return size, nil

}
