package fix

import (
	"errors"
	"github.com/google/go-cmp/cmp"
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
		var diff1 string
		var diff2 string
		diff1 = cmp.Diff(tc.want, got)
		diff2 = cmp.Diff(tc.wantMac, got)
		if got != nil {
			diff1 = cmp.Diff(tc.wantMac.Error(), got.Error())
			diff2 = cmp.Diff(tc.want.Error(), got.Error())
		}
		if diff1 != "" && diff2 != "" {
			t.Logf("%s", tc.name)
			t.Fatalf(diff1, diff2)
		}
	}

}
