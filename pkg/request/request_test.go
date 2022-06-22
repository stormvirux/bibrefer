package request

import (
	"github.com/google/go-cmp/cmp"
	"strings"
	"testing"
)

func TestDoiCrossRef(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"Name-Elseiver", "An Analysis of Fault Detection Strategies in Wireless Sensor Networks", "10.1016/j.jnca.2016.10.019"},
		{"Name-IEEE", "Distributed Inference Acceleration with Adaptive DNN Partitioning and Offloading", "10.1109/infocom41043.2020.9155237"},
		{"Name-Empty", "", ""},
		{"name-wierd", "vmrfvs ndkc", ""},
	}
	for _, tc := range tests {
		got, _ := DoiCrossRef(tc.query)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}

func TestDoiDataCite(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{"Name-ArXiv", "Dockerization Impacts in Database Performance Benchmarking", "10.48550/ARXIV.1812.04362"},
		{"Name-ArXivOld", "Linear probing and graphs", "10.48550/ARXIV.CS/9801103"},
		{"Name-Empty", "", ""},
		{"name-wierd", "vmrfvs ndkc", ""},
	}
	for _, tc := range tests {
		got, _ := DoiDataCite(tc.query)
		diff := cmp.Diff(tc.want, got)
		if diff != "" {
			t.Logf(tc.name)
			t.Fatalf(diff)
		}
	}
}

func TestRefDoi(t *testing.T) {
	jsonOut := strings.Fields(`{"indexed":{"date-parts":[[2022,4,1]],"date-time":"2022-04-01T15:19:44Z","timestamp":1648826384992},"reference-count":0,"publisher":"Institute of Electrical and Electronics Engineers (IEEE)","license":[{"start":{"date-parts":[[2021,1,1]],"date-time":"2021-01-01T00:00:00Z","timestamp":1609459200000},"content-version":"vor","delay-in-days":0,"URL":"https:\/\/ieeexplore.ieee.org\/Xplorehelp\/downloads\/license-information\/IEEE.html"},{"start":{"date-parts":[[2021,1,1]],"date-time":"2021-01-01T00:00:00Z","timestamp":1609459200000},"content-version":"am","delay-in-days":0,"URL":"https:\/\/ieeexplore.ieee.org\/Xplorehelp\/downloads\/license-information\/IEEE.html"}],"content-domain":{"domain":[],"crossmark-restriction":false},"published-print":{"date-parts":[[2021]]},"DOI":"10.1109\/tnse.2021.3126021","type":"journal-article","created":{"date-parts":[[2021,11,9]],"date-time":"2021-11-09T20:42:30Z","timestamp":1636490550000},"page":"1-1","source":"Crossref","is-referenced-by-count":0,"title":"Resource-constrained Federated Learning with Heterogeneous Data: Formulation and Analysis","prefix":"10.1109","author":[{"given":"Yi","family":"Liu","sequence":"first","affiliation":[]},{"given":"Yuanshao","family":"Zhu","sequence":"additional","affiliation":[]},{"given":"James J.Q.","family":"Yu","sequence":"additional","affiliation":[]}],"member":"263","container-title":"IEEE Transactions on Network Science and Engineering","original-title":[],"link":[{"URL":"http:\/\/xplorestaging.ieee.org\/ielx7\/6488902\/6930788\/09609654.pdf?arnumber=9609654","content-type":"unspecified","content-version":"vor","intended-application":"similarity-checking"}],"deposited":{"date-parts":[[2021,11,17]],"date-time":"2021-11-17T23:25:13Z","timestamp":1637191513000},"score":1,"resource":{"primary":{"URL":"https:\/\/ieeexplore.ieee.org\/document\/9609654\/"}},"subtitle":[],"short-title":[],"issued":{"date-parts":[[2021]]},"references-count":0,"URL":"http:\/\/dx.doi.org\/10.1109\/TNSE.2021.3126021","relation":{},"ISSN":["2327-4697","2334-329X"],"subject":["Computer Networks and Communications","Computer Science Applications","Control and Systems Engineering"],"container-title-short":"IEEE Trans. Netw. Sci. Eng.","published":{"date-parts":[[2021]]}}`)
	xmlOut := strings.Fields(`<rdf:RDF
    xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
    xmlns:j.0="http://purl.org/dc/terms/"
    xmlns:j.1="http://prismstandard.org/namespaces/basic/2.1/"
    xmlns:owl="http://www.w3.org/2002/07/owl#"
    xmlns:j.2="http://purl.org/ontology/bibo/"
    xmlns:j.3="http://xmlns.com/foaf/0.1/">
  <rdf:Description rdf:about="http://dx.doi.org/10.1109/TNSE.2021.3126021">
    <j.0:publisher>Institute of Electrical and Electronics Engineers (IEEE)</j.0:publisher>
    <j.0:creator>
      <j.3:Person rdf:about="http://id.crossref.org/contributor/yuanshao-zhu-32yh171znboo3">
        <j.3:name>Yuanshao Zhu</j.3:name>
        <j.3:familyName>Zhu</j.3:familyName>
        <j.3:givenName>Yuanshao</j.3:givenName>
      </j.3:Person>
    </j.0:creator>
    <j.1:doi>10.1109/tnse.2021.3126021</j.1:doi>
    <j.0:identifier>10.1109/tnse.2021.3126021</j.0:identifier>
    <owl:sameAs rdf:resource="doi:10.1109/tnse.2021.3126021"/>
    <j.2:pageStart>1</j.2:pageStart>
    <j.2:pageEnd>1</j.2:pageEnd>
    <j.0:isPartOf>
      <j.2:Journal rdf:about="http://id.crossref.org/issn/2327-4697">
        <owl:sameAs>urn:issn:2327-4697</owl:sameAs>
        <j.1:issn>2334-329X</j.1:issn>
        <j.2:issn>2334-329X</j.2:issn>
        <owl:sameAs>urn:issn:2334-329X</owl:sameAs>
        <j.0:title>IEEE Transactions on Network Science and Engineering</j.0:title>
        <j.1:issn>2327-4697</j.1:issn>
        <j.2:issn>2327-4697</j.2:issn>
      </j.2:Journal>
    </j.0:isPartOf>
    <j.0:date rdf:datatype="http://www.w3.org/2001/XMLSchema#gYear"
    >2021</j.0:date>
    <owl:sameAs rdf:resource="http://dx.doi.org/10.1109/tnse.2021.3126021"/>
    <j.1:endingPage>1</j.1:endingPage>
    <j.2:doi>10.1109/tnse.2021.3126021</j.2:doi>
    <owl:sameAs rdf:resource="info:doi/10.1109/tnse.2021.3126021"/>
    <j.0:creator>
      <j.3:Person rdf:about="http://id.crossref.org/contributor/yi-liu-32yh171znboo3">
        <j.3:name>Yi Liu</j.3:name>
        <j.3:familyName>Liu</j.3:familyName>
        <j.3:givenName>Yi</j.3:givenName>
      </j.3:Person>
    </j.0:creator>
    <j.0:title>Resource-constrained Federated Learning with Heterogeneous Data: Formulation and Analysis</j.0:title>
    <j.0:creator>
      <j.3:Person rdf:about="http://id.crossref.org/contributor/james-j-q-yu-32yh171znboo3">
        <j.3:name>James J.Q. Yu</j.3:name>
        <j.3:familyName>Yu</j.3:familyName>
        <j.3:givenName>James J.Q.</j.3:givenName>
      </j.3:Person>
    </j.0:creator>
    <j.1:startingPage>1</j.1:startingPage>
  </rdf:Description>
</rdf:RDF>
`)
	tnsmOut := strings.Fields(`@article{Liu_2021,
        doi = {10.1109/tnse.2021.3126021},
        url = {https://doi.org/10.1109%2Ftnse.2021.3126021},
        year = 2021,
        publisher = {Institute of Electrical and Electronics Engineers ({IEEE})},
        pages = {1--1},
        author = {Yi Liu and Yuanshao Zhu and James J.Q. Yu},
        title = {Resource-constrained Federated Learning with Heterogeneous Data: Formulation and Analysis},
        journal = {{IEEE} Transactions on Network Science and Engineering}
}`)

	acmJOut := strings.Fields(`@article{Ma_2022,
        doi = {10.1145/3508461},
        url = {https://doi.org/10.1145%2F3508461},
        year = 2022,
        month = {oct},
        publisher = {Association for Computing Machinery ({ACM})},
        volume = {41},
        number = {5},
        pages = {1--18},
        author = {Karima Ma and Michael Gharbi and Andrew Adams and Shoaib Kamil and Tzu-Mao Li and Connelly Barnes and Jonathan Ragan-Kelley},
        title = {Searching for Fast Demosaicking Algorithms},
        journal = {{ACM} Transactions on Graphics}
}`)
	acmPOut := strings.Fields(`@inproceedings{Shin_2019,
        doi = {10.1145/3325413.3329788},
        url = {https://doi.org/10.1145%2F3325413.3329788},
        year = 2019,
        publisher = {{ACM} Press},
        author = {Kwang Yong Shin and Hyuk-Jin Jeong and Soo-Mook Moon},
        title = {Enhanced Partitioning of {DNN} Layers for Uploading from Mobile Devices to Edge Servers},
        booktitle = {The 3rd International Workshop on Deep Learning for Mobile Systems and Applications  - {EMDL} {\textquotesingle}19}
}`)
	ieeePOut := strings.Fields(`@inproceedings{Jedari_2018,
        doi = {10.1109/glocom.2018.8647205},
        url = {https://doi.org/10.1109%2Fglocom.2018.8647205},
        year = 2018,
        month = {dec},
        publisher = {{IEEE}},
        author = {Behrouz Jedari and Mario Di Francesco},
        title = {Delay Analysis of Layered Video Caching in Crowdsourced Heterogeneous Wireless Networks},
        booktitle = {2018 {IEEE} Global Communications Conference ({GLOBECOM})}
}`)
	elseiverOut := strings.Fields(`@article{Muhammed_2017,
         doi = {10.1016/j.jnca.2016.10.019},
         url = {https://doi.org/10.1016%2Fj.jnca.2016.10.019},
         year = 2017,
         month = {jan},
         publisher = {Elsevier {BV}},
         volume = {78},
         pages = {267--287},
         author = {Thaha Muhammed and Riaz Ahmed Shaikh},
        title = {An analysis of fault detection strategies in wireless sensor networks},
        journal = {Journal of Network and Computer Applications}
}`)
	springerOut := strings.Fields(`@article{Dolbeau_2017,
        doi = {10.1007/s11227-017-2177-5},
        url = {https://doi.org/10.1007%2Fs11227-017-2177-5},
        year = 2017,
        month = {nov},
        publisher = {Springer Science and Business Media {LLC}},
        volume = {74},
        number = {3},
        pages = {1341--1377},
        author = {Romain Dolbeau},
        title = {Theoretical peak {FLOPS} per instruction set: a tutorial},
        journal = {The Journal of Supercomputing}
}`)

	tests := []struct {
		name   string
		query  string
		output string
		want   []string
	}{
		{"IEEE-J", "10.1109/TNSE.2021.3126021", "bibtex", tnsmOut},
		{"IEEE-J-XML", "10.1109/TNSE.2021.3126021", "xml", xmlOut},
		{"IEEE-J-XML", "10.1109/TNSE.2021.3126021", "json", jsonOut},
		{"ACM-J", "10.1145/3508461", "bibtex", acmJOut},
		{"ACM-P", "10.1145/3325413.3329788", "bibtex", acmPOut},
		{"IEEE-PNamed", "10.1109/glocom.2018.8647205", "bibtex", ieeePOut},
		{"Elseiver", "10.1016/j.jnca.2016.10.019", "bibtex", elseiverOut},
		{"Springer", "10.1007/s11227-017-2177-5", "bibtex", springerOut},
	}
	for _, tc := range tests {
		got, _ := RefDoi(tc.query, tc.output)
		diff := cmp.Diff(tc.want, strings.Fields(got))
		if diff != "" {
			t.Fatalf(diff)
		}
	}
}
