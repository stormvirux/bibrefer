package ref

import (
	"fmt"
	"io"
	"strings"
)

func getIndexBib(bibLinebyline []string) map[string]uint {
	indexMap := map[string]uint{"journalIndex": 255,
		"authorIndex": 255,
		"urlIndex":    255,
		"monthIndex":  255,
		"doiIndex":    255,
		"yearIndex":   255,
	}

	for i := 0; i < len(bibLinebyline); i++ {
		if strings.Contains(bibLinebyline[i], "author") {
			indexMap["authorIndex"] = uint(i)
			continue
		}
		if strings.Contains(bibLinebyline[i], "journal") {
			indexMap["journalIndex"] = uint(i)
			continue
		}
		if strings.Contains(bibLinebyline[i], "url") {
			indexMap["urlIndex"] = uint(i)
			continue
		}
		if strings.Contains(bibLinebyline[i], "month") {
			indexMap["monthIndex"] = uint(i)
			continue
		}
		if strings.Contains(bibLinebyline[i], "doi") {
			indexMap["doiIndex"] = uint(i)
			continue
		}
		if strings.Contains(bibLinebyline[i], "year") {
			indexMap["yearIndex"] = uint(i)
			continue
		}
	}
	return indexMap
}

// Only use to avoid cyclomatic complexity
func verbosePrint(isVerbose bool, message string, stream io.Writer) {
	if isVerbose {
		_, _ = fmt.Fprintln(stream, message)
	}
}
