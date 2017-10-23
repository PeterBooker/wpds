package slurper

import (
	"os"
	"path/filepath"

	"github.com/peterbooker/wpds2/internal/pkg/context"
)

func extensionFailure(item string, ctx *context.Context) {

	ctx.Stats.IncrementTotalExtensionsFailed()

	fname := ".failed-downloads"

	path := filepath.Join(wd, ctx.ExtensionType, fname)

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

	_, err = f.WriteString(item + "\n")
	if err != nil {
		return
	}

	return

}

//func getFailedList(extType string) []string {

//}
