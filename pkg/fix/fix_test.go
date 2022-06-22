package fix

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestPdf_Fix(t *testing.T) {
	var err = errors.New(`failure repairing the pdf: fork/exec : no such file or directory`)
	var err2 = errors.New("failure repairing the pdf: exec: \"gssfedfe\": executable file not found in $PATH")
	tests := []struct {
		name    string
		pdfName string
		binary  string
		want    error
	}{
		{"gsPath", "testdata/test1.pdf", "/usr/bin/gs", nil},
		{"noGsPath", "testdata/test1.pdf", "", err},
		{"wrongGsPath", "testdata/test1.pdf", "/usr/bin/gssfedfe", err2},
		{"boPdfPath", "", "/usr/bin/gs", errors.New("failure repairing the pdf: exit status 1")},
	}
	for _, tc := range tests {
		pdfTest := Pdf{Name: tc.pdfName}
		got := pdfTest.Fix(tc.binary)
		var diff string
		diff = cmp.Diff(tc.want, got)
		if got != nil {
			diff = cmp.Diff(tc.want.Error(), got.Error())
		}
		if diff != "" {
			t.Logf("%s", tc.name)
			t.Fatalf(diff)
		}
	}

}
