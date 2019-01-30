package grep

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
)

type GrepResults []*GrepResult

type GrepResult struct {
	Filename   string
	LineNumber int
	Line       string
}

func (r *GrepResult) TSV() string {
	return r.join("\t")
}

func (r *GrepResult) CSV() string {
	return r.join(",")
}

func (r *GrepResult) join(sep string) string {
	return strings.Join([]string{
		r.Filename,
		strconv.Itoa(r.LineNumber),
		r.Line,
	}, sep)
}

func ParseGrepResult(result []byte) GrepResults {
	var grepResults []*GrepResult

	r := bytes.NewReader(result)
	reader := bufio.NewReader(r)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}

		grepResults = append(grepResults, parseGrepResultLine(line))
	}

	return grepResults
}

func parseGrepResultLine(line []byte) *GrepResult {
	parts := strings.Split(string(line), ":")
	switch len(parts) {
	case 1:
		return &GrepResult{
			Filename: parts[0],
		}
	case 2:
		return &GrepResult{
			Filename: parts[0],
			Line:     parts[1],
		}
	default:
		lineNumber, _ := strconv.Atoi(parts[1])
		return &GrepResult{
			Filename:   parts[0],
			LineNumber: lineNumber,
			Line:       strings.Join(parts[2:], " "),
		}
	}
}
