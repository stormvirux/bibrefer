package ref

import (
	"bytes"
	"fmt"
	"github.com/nickng/bibtex"
	"io"
	"strconv"
	"strings"
)

// Only use to avoid cyclomatic complexity.
func verbosePrint(isVerbose bool, message string, stream io.Writer) {
	if isVerbose {
		_, _ = fmt.Fprintln(stream, message)
	}
}

func prettyPrint(entry *bibtex.BibEntry) string {
	var bt bytes.Buffer
	bt.WriteString(fmt.Sprintf("@%s{%s,\n", entry.Type, entry.CiteName))
	for key, val := range entry.Fields {
		if i, err := strconv.Atoi(strings.TrimSpace(val.String())); err == nil {
			bt.WriteString(fmt.Sprintf("    %-10s = %d,\n", key, i))
		} else {
			bt.WriteString(fmt.Sprintf("    %-10s = {%s},\n", key, strings.TrimSpace(val.String())))
		}
	}
	bt.Truncate(bt.Len() - 2)
	bt.WriteString("\n}\n")
	return bt.String()
}
