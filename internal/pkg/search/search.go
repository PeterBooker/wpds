package search

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/fatih/color"
	"github.com/peterbooker/wpds/internal/pkg/stats"
	"github.com/peterbooker/wpds/internal/pkg/utils"
)

const (
	// MB holds the number of bytes in a megabyte
	MB = 1000000
)

// Search ...
type Search struct {
	input   *regexp.Regexp
	path    string
	Results []Match
	context *Context
}

// Match contains details for a search match.
type Match struct {
	Extension string
	Filename  string
	Path      string
	Line      int
	Text      string
}

// Context ...
type Context struct {
	ExtensionType    string
	FileType         string
	ExtWhitelist     []string
	WorkingDirectory string
	Stats            *stats.Stats
}

// Setup ...
func Setup(input string, ctx *Context) *Search {

	searchPath := filepath.Join(ctx.WorkingDirectory, ctx.ExtensionType)

	if !utils.DirExists(searchPath) || utils.IsDirEmpty(searchPath) {
		log.Fatalf("Nothing to search at specified location: %s", searchPath)
	}

	// TODO: Implement this as context option
	ignoreCase := false

	var regex *regexp.Regexp
	var err error

	if ignoreCase {
		regex, err = regexp.Compile(`(?i)(` + input + `)`)
	} else {
		regex, err = regexp.Compile(`(` + input + `)`)
	}

	if err != nil {
		log.Fatalf("Cannot compile regex, invalid syntax: %s\n", input)
	}

	return &Search{
		input:   regex,
		path:    searchPath,
		context: ctx,
	}

}

// Run ...
func (s *Search) Run() error {

	extensions := ListDirNames(s.path)

	fmt.Printf("Total Extensions: %d\n", len(extensions))

	var wg sync.WaitGroup

	limiter := make(chan struct{}, 18)

	for _, extension := range extensions {

		limiter <- struct{}{}

		path := filepath.Join(s.path, extension)

		wg.Add(1)
		// Start a goroutine to fetch the folder.
		go func(root string, extension string) {

			// Decrement the counter when the goroutine completes.
			defer wg.Done()

			filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

				pathOffset := len(root) - len(extension)

				if err != nil {
					//fmt.Println(err)
				}

				// Directories cannot be searched
				if info.IsDir() {
					return nil
				}

				if info.Size() > (20 * MB) {
					return fmt.Errorf("File exceeds valid size for searching. Filename: %s Size: %d", path, info.Size())
				}

				data, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}

				// Ignore files which are not utf8 encoded e.g. binary files- images, etc.
				if !isUTF8(data) {
					return nil
				}

				matches := s.input.FindAllIndex(data, -1)

				if matches != nil {

					for _, match := range matches {

						start, end := getMatchLineIndexes(data, match)

						line := getLineNum(data[:end])

						var preText, matchText, postText, fullText string
						preText = string(data[start:match[0]])
						matchText = string(data[match[0]:match[1]])
						postText = string(data[match[1]:end])
						fullText = string(data[start:end])

						yellow := color.New(color.FgYellow).SprintFunc()
						red := color.New(color.FgRed).SprintFunc()
						green := color.New(color.FgGreen).SprintFunc()

						fmt.Println("")
						fmt.Printf("%s\n", green(path[pathOffset:len(path)]))
						fmt.Printf("%s: %s%s%s\n", red(line), preText, yellow(matchText), postText)

						m := Match{
							Extension: extension,
							Filename:  extension + string(os.PathSeparator) + info.Name(),
							Path:      path,
							Line:      line,
							Text:      fullText,
						}

						s.Results = append(s.Results, m)

					}

				}

				return nil

			})

			<-limiter

		}(path, extension)

	}

	wg.Wait()

	printSummary(s.context)

	// TODO: Write search results to file.

	return nil

}

// ListDirNames lists all Directories for a type of extension.
func ListDirNames(path string) []string {

	var dirs []string

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}

	}

	return dirs

}
