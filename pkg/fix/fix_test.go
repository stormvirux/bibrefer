package fix

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"runtime"
	"testing"
)

func TestPdf_Fix(t *testing.T) {
	var err = errors.New(`failure repairing the pdf: fork/exec : no such file or directory`)
	var errMacOs = errors.New(`failure repairing the pdf: exec: no command`)
	var err2 = errors.New("failure repairing the pdf: exec: \"gssfedfe\": executable file not found in $PATH")
	tests := []struct {
		name    string
		pdfName string
		binary  string
		want    error
		wantMac error
	}{
		{"gsPath", "testdata/test1.pdf", "/usr/bin/gs", nil, nil},
		{"noGsPath", "testdata/test1.pdf", "", err, errMacOs},
		{"wrongGsPath", "testdata/test1.pdf", "/usr/bin/gssfedfe", err2, err2},
		{"boPdfPath", "", "/usr/bin/gs", errors.New("failure repairing the pdf: exit status 1"), errors.New("failure repairing the pdf: exit status 1")},
	}
	for _, tc := range tests {
		pdfTest := Pdf{Name: tc.pdfName}
		got := pdfTest.Fix(tc.binary)
		var diff string
		diff = cmp.Diff(tc.want, got)
		if got != nil && runtime.GOOS == "darwin" {
			diff = cmp.Diff(tc.wantMac.Error(), got.Error())
		} else if got != nil {
			diff = cmp.Diff(tc.want.Error(), got.Error())
		}

		if diff != "" {
			t.Logf("%s", tc.name)
			t.Fatalf(diff)
		}
	}

}
