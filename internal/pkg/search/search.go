package search

import (
	"os"
	"regexp"
)

// StringSearch ...
type StringSearch struct {
	input   string
	Results Results
}

// RegexSearch ...
type RegexSearch struct {
	input   string
	regex   *regexp.Regexp
	Results Results
}

// Results ...
type Results struct {
	Matches []Result
}

// Result ...
type Result struct {
	File string
	Line int
	Text string
}

// NewString ...
func NewString(input string) *StringSearch {

	return &StringSearch{
		input: input,
	}

}

// NewRegex ...
func NewRegex(input string) *RegexSearch {

	regex, err := regexp.Compile(input)
	if err != nil {
		os.Exit(1)
	}

	return &RegexSearch{
		input: input,
		regex: regex,
	}

}
