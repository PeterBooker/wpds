package search

import (
	"bytes"
	"unicode/utf8"
)

const maxLen = 512

// isUTF8 checks upto the first 512 bytes for valid UTF8 encoding.
func isUTF8(data []byte) bool {
	if len(data) > maxLen {
		data = data[:maxLen]
	}
	return utf8.Valid(data)
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) ([]byte, bool) {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1], true
	}
	return data, false
}

// getLineNum counts the lines in data.
// TODO: Check this is always accurate.
func getLineNum(data []byte) int {
	return 1 + bytes.Count(data, []byte("\n"))
}

// getMatchLines finds the index for the start of the first line and the end of the last line.
func getMatchLineIndexes(data []byte, match []int) (int, int) {
	start := getStartIndex(data, match[0])
	end := getEndIndex(data, match[1])
	return start, end
}

// getStartIndex begins from the start of the match and finds the start of that line.
func getStartIndex(data []byte, match int) int {
	start := match
	if match > 0 {

		for i := match - 1; i > 0; i-- {
			if data[i] == '\n' {
				start = i + 1
				break
			}
		}

	}
	return start
}

// getEndIndex begins from the end of the match and finds the end of that line.
func getEndIndex(data []byte, match int) int {
	end := match
	max := len(data)
	if match < max {

		for i := match + 1; i <= max; i++ {
			if i == max {
				end = i
				break
			}
			if data[i] == '\n' || data[i] == '\r' {
				end = i - 1
				break
			}
		}

	}
	return end
}
