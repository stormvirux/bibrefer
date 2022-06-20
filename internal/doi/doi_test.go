package doi

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
		{"perfect arxiv", []string{"testdata/test1.pdf"}, []bool{true, true, false}, "10.1109/TNSE.2021.3126021"},
		{"repair pdf", []string{"testdata/test2.pdf"}, []bool{false, false, false}, "10.1109/TNSE.2021.3126021"},
		{"repair arxiv", []string{"testdata/test3.pdf"}, []bool{false, true, false}, "10.48550/arXiv.1602.05629"},
		{"repair ieee", []string{"testdata/test4.pdf"}, []bool{false, true, false}, "10.1007/978-3-030-13705-2_18"},
		{"repair springer", []string{"testdata/test5.pdf"}, []bool{false, true, false}, "10.1145/3410530.3414593"},
		{"repair mdpi", []string{"testdata/test7.pdf"}, []bool{false, true, false}, "10.1109/ACCESS.2018.2846609"},
		{"repair access", []string{"testdata/test8.pdf"}, []bool{false, true, false}, "10.1016/j.jnca.2016.10.019"},
		{"repair elseiver", []string{"testdata/test6.pdf"}, []bool{false, true, false}, "10.3390/app9050947"},
	}

	for _, tc := range tests {
		pdfTest := Doi{Verbose: tc.flags[2], Arxiv: tc.flags[0], Clip: tc.flags[1]}
		got, _ := pdfTest.Run(tc.query)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}

func TestApp_RunCrossRef(t *testing.T) {
	tests := []struct {
		name  string
		query []string
		flags []bool
		want  string
	}{
		{"perfect", []string{"Resource-constrained Federated Edge Learning with Heterogeneous Data: Formulation and Analysis"},
			[]bool{false, false, false}, "10.1109/tnse.2021.3126021"},
		{"perfect verbose", []string{"Resource-constrained Federated Edge Learning with Heterogeneous Data: Formulation and Analysis"},
			[]bool{false, true, false}, "10.1109/tnse.2021.3126021"},
		{"repair ieee", []string{"Collaborative Edge-Network Content Replication: A Joint User Preference and Mobility Approach"},
			[]bool{false, false, false}, "10.1145/3410530.3414593"},
		{"repair springer", []string{"HPC-Smart Infrastructures: A Review and Outlook on Performance Analysis Methods and Tools"},
			[]bool{false, false, false}, "10.1007/978-3-030-13705-2_18"},
		{"repair mdpi", []string{"UbeHealth: A Personalized Ubiquitous Cloud and Edge-Enabled Networked Healthcare System for Smart Cities"},
			[]bool{false, false, false}, "10.1109/access.2018.2846609"},
		{"repair access", []string{"An analysis of fault detection strategies in wireless sensor networks"},
			[]bool{false, false, false}, "10.1016/j.jnca.2016.10.019"},
		{"repair elseiver", []string{"SURAA: A Novel Method and Tool for Loadbalanced and Coalesced SpMV Computations on GPUs"},
			[]bool{false, false, false}, "10.3390/app9050947"},
	}
	for _, tc := range tests {
		pdfTest := Doi{Verbose: tc.flags[2], Arxiv: tc.flags[0], Clip: tc.flags[1]}
		got, _ := pdfTest.Run(tc.query)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}

func TestApp_RunDataCite(t *testing.T) {
	tests := []struct {
		name  string
		query []string
		flags []bool
		want  string
	}{
		{"repair arxiv1", []string{"Orbital Semilattices"},
			[]bool{true, false, false}, "10.48550/ARXIV.2206.07790"},
		{"repair arxiv1", []string{"Deep Learning and Handheld Augmented Reality Based System for Optimal Data Collection in Fault Diagnostics Domain"},
			[]bool{true, false, false}, "10.48550/ARXIV.2206.07772"},
		{"repair arxiv1", []string{"An Infinite Dimensional Model for a Many Server Priority Queue"},
			[]bool{true, false, false}, "10.48550/ARXIV.1701.01328"},
		{"repair arxiv1", []string{"Process, Bias and Temperature Scalable CMOS Analog Computing Circuits for Machine Learning"},
			[]bool{true, false, false}, "10.48550/ARXIV.2205.05664"},
	}
	for _, tc := range tests {
		pdfTest := Doi{Verbose: tc.flags[2], Arxiv: tc.flags[0], Clip: tc.flags[1]}
		got, _ := pdfTest.Run(tc.query)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}
