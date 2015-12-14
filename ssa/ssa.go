package ssa

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/martinlindhe/subber/caption"
	"github.com/martinlindhe/subber/txtformat"
)

func ParseSsa(b []byte) []caption.Caption {

	s := txtformat.ConvertToUTF8(b)

	var res []caption.Caption

	chunk, err := extractSsaChunk("[Events]", s)
	if err != nil {
		fmt.Println("[ssa]", err)
		return nil
	}

	lines := strings.Split(chunk, "\n")

	// first line must now be a  "Format:"
	format := lines[0]

	// parse format into column numbers
	startCol := parseSsaFormat(format, "Start")
	endCol := parseSsaFormat(format, "End")
	textCol := parseSsaFormat(format, "Text")

	dialogueColumns := columnCountFromSsaFormat(format)

	outSeq := 1

	for i := 1; i < len(lines); i++ {

		seq := strings.Trim(lines[i], "\r ")
		if seq == "" {
			break
		}

		start := parseSsaDialogue(lines[i], startCol, dialogueColumns)
		end := parseSsaDialogue(lines[i], endCol, dialogueColumns)
		text := parseSsaDialogue(lines[i], textCol, dialogueColumns)

		var o caption.Caption
		o.Seq = outSeq
		o.Start, err = parseSsaTime(start)
		if err != nil {
			fmt.Println(err)
			continue
		}
		o.End, err = parseSsaTime(end)
		if err != nil {
			fmt.Println(err)
			continue
		}

		o.Text = strings.Split(text, "\\n")

		if len(o.Text) > 0 {
			res = append(res, o)
			outSeq++
		}
	}

	return res
}

// return column idx from s
func parseSsaDialogue(s string, idx int, columns int) string {

	pos := strings.Index(s, ": ")
	if pos == -1 {
		return ""
	}

	s = s[pos+2:]
	cols := strings.SplitN(s, ",", columns)

	return strings.TrimSpace(cols[idx])
}

func parseSsaFormat(s string, colName string) int {

	pos := strings.Index(s, ": ")
	if pos == -1 {
		return -1
	}

	s = s[pos+2:]

	cols := strings.Split(s, ",")

	for idx, col := range cols {
		col = strings.TrimSpace(col)
		if col == colName {
			return idx
		}
	}

	// -1 if not found
	return -1
}

func columnCountFromSsaFormat(s string) int {
	pos := strings.Index(s, ": ")
	if pos == -1 {
		return -1
	}

	s = s[pos+2:]

	cols := strings.Split(s, ",")
	return len(cols)
}

func extractSsaChunk(chunk string, s string) (string, error) {

	pos := strings.Index(s, chunk)
	if pos == -1 {
		return "", fmt.Errorf("Parse error in chunk")
	}

	// XXX this will break if there is a chunk after [Events]
	res := s[pos+len(chunk):]

	return strings.Trim(res, "\r\n "), nil
}

func makeTime(h int, m int, s int, ms int) time.Time {
	return time.Date(0, 1, 1, h, m, s, ms*1000*1000, time.UTC)
}

func parseSsaTime(in string) (time.Time, error) {
	// "0:01:06.37" => h:mm:ss.cc (centisec)

	r1 := regexp.MustCompile("([0-9]+):([0-9]+):([0-9]+)[.]([0-9]+)")

	matches := r1.FindStringSubmatch(in)

	if len(matches) < 5 {
		return time.Now(), fmt.Errorf("[srt] Regexp didnt match: %s", in)
	}

	h, err := strconv.Atoi(matches[1])
	if err != nil {
		return time.Now(), err
	}

	m, err := strconv.Atoi(matches[2])
	if err != nil {
		return time.Now(), err
	}

	s, err := strconv.Atoi(matches[3])
	if err != nil {
		return time.Now(), err
	}

	cs, err := strconv.Atoi(matches[4])
	if err != nil {
		return time.Now(), err
	}

	return makeTime(h, m, s, cs*10), nil
}