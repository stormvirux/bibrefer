package getdoi

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

// https://eli.thegreenplace.net/2022/file-driven-testing-in-go/
func TestApp_RunFile(t *testing.T) {
	tests := []struct {
		name  string
		query []string
		flags []bool
		want  string
	}{
		{"perfect", []string{"testdata/test1.pdf"}, []bool{false, false, false}, "10.1109/TNSE.2021.3126021"},
		{"perfect verbose", []string{"testdata/test1.pdf"}, []bool{false, true, false}, "10.1109/TNSE.2021.3126021"},
		{"perfect arxiv", []string{"testdata/test1.pdf"}, []bool{true, false, false}, "10.1109/TNSE.2021.3126021"},
		{"repair pdf", []string{"testdata/test2.pdf"}, []bool{false, false, false}, "10.1109/TNSE.2021.3126021"},
		{"repair arxiv pdf", []string{"testdata/test3.pdf"}, []bool{false, false, false}, "10.48550/arXiv.1602.05629"},
		{"repair ieee pdf", []string{"testdata/test4.pdf"}, []bool{false, false, false}, "10.1007/978-3-030-13705-2_18"},
		{"repair pdf", []string{"testdata/test5.pdf"}, []bool{false, false, false}, "10.1145/3410530.3414593"},
		{"repair pdf", []string{"testdata/test7.pdf"}, []bool{false, false, false}, "10.1109/ACCESS.2018.2846609"},
		{"repair pdf", []string{"testdata/test8.pdf"}, []bool{false, false, false}, "10.1016/j.jnca.2016.10.019"},
		{"repair pdf", []string{"testdata/test6.pdf"}, []bool{false, false, false}, "10.3390/app9050947"},
	}

	for _, tc := range tests {
		pdfTest := App{}
		got, _ := pdfTest.Run(tc.query, tc.flags)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}

func TestApp_RunCrossRef(t *testing.T) {
}

func TestApp_RunDataCite(t *testing.T) {

}
