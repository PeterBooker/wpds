package slurper

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/peterbooker/wpds/internal/pkg/context"
)

const (
	filename = ".failed-downloads"
)

// extensionFailure appends an extension name to the .failed-downloads file.
func extensionFailure(name string, ctx *context.Context) {

	ctx.Stats.IncrementTotalExtensionsFailed()

	path := filepath.Join(wd, ctx.ExtensionType, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {

		f, err := os.Create(path)
		defer f.Close()
		if err != nil {
			return
		}

	}

	f, err := os.OpenFile(path, os.O_APPEND, 0777)
	if err != nil {
		return
	}

	defer f.Close()

	_, err = f.WriteString(name + "\n")
	if err != nil {
		return
	}

	return

}

// getFailedList reads the .failed-downloads file and returns a list of extensions.
func getFailedList(extType string) []string {

	var list []string

	wd, _ := os.Getwd()

	path := filepath.Join(wd, extType, filename)

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return list
	}
	defer file.Close()

	r := bufio.NewReader(file)

	for {

		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		list = append(list, string(line))

	}

	return list

}

// removeFailedList deletes the .failed-downloads file.
func removeFailedList(extType string) error {

	wd, _ := os.Getwd()

	path := filepath.Join(wd, extType, filename)

	return os.Remove(path)

}
