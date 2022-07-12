## Bibrefer

[![Release](https://img.shields.io/github/release/stormvirux/bibrefer.svg?style=flat-square)](https://github.com/stormvirux/bibrefer/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/stormvirux/bibrefer)](https://goreportcard.com/report/github.com/stormvirux/bibrefer)
[![Coverage Status](https://coveralls.io/repos/github/stormvirux/bibrefer/badge.svg)](https://coveralls.io/github/stormvirux/bibrefer)
![Build Status](https://github.com/stormvirux/bibrefer/actions/workflows/bibrefer.yml/badge.svg)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B32136%2Fgithub.com%2Fstormvirux%2Fbibrefer.svg?type=shield)](https://app.fossa.com/projects/custom%2B32136%2Fgithub.com%2Fstormvirux%2Fbibrefer?ref=badge_shield)

Bibrefer is a CLI application written in pure Go that can fetch consistent references of publications in bibtex, json, 
or xml formats. Influenced by [scholarref](https://adamsgaard.dk/scholarref.html), a collection of POSIX shell scripts, Bibref is a faster cross-platform solution 
that enhances the functionalities of scholarref. Bibrefer is almost free of external dependencies except in small use-cases 
where the reference need to be extracted from a pdf file with corrupted xref table. In such cases, bibrefer shall attempt to
use ghostscript if installed, to repair the pdf and extract the references.


Bifrefer consists of two sub commands:  `doi` and `ref`. The `doi` sub command is used to fetch DOI using a given publication name
or a pdf file. The references by default is fetched from CrossRef. ArXiv publication details can be obtained from DataCite 
with an additional flag. The `ref` sub command is used to fetch references from doi.org using a doi.

### Installation
The binaries for bibrefer are available on [Releases](https://github.com/stormvirux/bibrefer/releases/tag/v1.0.0).

You can also download and build the source code from the [release](https://github.com/stormvirux/bibrefer/releases/tag/v1.0.0) page.

### `bibrefer doi`

```
bibrefer doi [flags] <query>
```

#### Flags

```
  -c, --clip       copy the DOI to clipboard
  -d, --datacite   retrieve the DOI from DataCite (for ArXiv)
  -h, --help       help for doi
  -V, --verbose    show verbose information
```

query is the name of a publication or a pdf file.

### `bibrefer ref`

```
bibrefer ref [flags] <query>
```
`query` is the DOI of a publication.

#### Flags

```
  -a, --full-author     use full author names in reference
  -j, --full-journal    use full  journal name in reference
  -h, --help            help for ref
  -k, --keep-key        keep the bib entry key format from doi.org
  -n, --no-newline      suppress trailing newline but prepend with newline
  -o, --output string   sets the output format. Supported values are json, bibtex, and rdf-xml. (default "bibtex")
  -V, --verbose         show verbose information
```


### Examples

Fetch DOI for the publication titled "An analysis of fault detection strategies in wireless sensor networks"
```
$ bibrefer doi An analysis of fault detection strategies in wireless sensor networks

10.1016/j.jnca.2016.10.019
```

Fetch DOI from a pdf file with the DOI `10.1016/j.jnca.2016.10.019` or the title "An analysis of fault detection strategies in wireless sensor networks"
```
$ bibrefer doi publication.pdf

10.1016/j.jnca.2016.10.019
```

Fetch the reference of a publication with DOI: 10.1016/j.jnca.2016.10.019 having full author name.
```
$ bibrefer ref -a 10.1016/j.jnca.2016.10.019

@article{Muhammed:2017:JNCA,
        doi = {10.1016/j.jnca.2016.10.019},
        year = 2017,
        publisher = {Elsevier {BV}},
        volume = {78},
        pages = {267--287},
        author = {Thaha Muhammed and Riaz Ahmed Shaikh},
        title = {An analysis of fault detection strategies in wireless sensor networks},
        journal = {J.  Network  Comput. Appl.}
}
```

Fetch the reference of a publication with DOI: 10.1016/j.jnca.2016.10.019 in xml format with default formatting.
```
```rdf-xml
$ bibrefer ref -o xml 10.1016/j.jnca.2016.10.019

<rdf:RDFefer ref -o xml 10.1016/j.jnca.2016.10.019
        xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
        xmlns:j.0="http://purl.org/dc/terms/"
        xmlns:j.1="http://prismstandard.org/namespaces/basic/2.1/"
        xmlns:owl="http://www.w3.org/2002/07/owl#"
        xmlns:j.2="http://purl.org/ontology/bibo/"
        xmlns:j.3="http://xmlns.com/foaf/0.1/">
<rdf:Description rdf:about="http://dx.doi.org/10.1016/j.jnca.2016.10.019">
<j.0:publisher>Elsevier BV</j.0:publisher>
<j.1:startingPage>267</j.1:startingPage>
<j.1:doi>10.1016/j.jnca.2016.10.019</j.1:doi>
<owl:sameAs rdf:resource="doi:10.1016/j.jnca.2016.10.019"/>
<owl:sameAs rdf:resource="info:doi/10.1016/j.jnca.2016.10.019"/>
<j.0:title>An analysis of fault detection strategies in wireless sensor networks</j.0:title>
<j.2:pageStart>267</j.2:pageStart>
<j.2:pageEnd>287</j.2:pageEnd>
<j.0:creator>
    <j.3:Person rdf:about="http://id.crossref.org/contributor/riaz-ahmed-shaikh-19gce4u4co3gj">
        <j.3:name>Riaz Ahmed Shaikh</j.3:name>
        <j.3:familyName>Shaikh</j.3:familyName>
        <j.3:givenName>Riaz Ahmed</j.3:givenName>
    </j.3:Person>
</j.0:creator>
<j.1:endingPage>287</j.1:endingPage>
<j.0:date rdf:datatype="http://www.w3.org/2001/XMLSchema#gYearMonth"
>2017-01</j.0:date>
<j.0:isPartOf>
    <j.2:Journal rdf:about="http://id.crossref.org/issn/1084-8045">
        <owl:sameAs>urn:issn:1084-8045</owl:sameAs>
        <j.0:title>Journal of Network and Computer Applications</j.0:title>
        <j.1:issn>1084-8045</j.1:issn>
        <j.2:issn>1084-8045</j.2:issn>
    </j.2:Journal>
</j.0:isPartOf>
<j.1:volume>78</j.1:volume>
<j.2:volume>78</j.2:volume>
<j.0:identifier>10.1016/j.jnca.2016.10.019</j.0:identifier>
<j.2:doi>10.1016/j.jnca.2016.10.019</j.2:doi>
<j.0:creator>
    <j.3:Person rdf:about="http://id.crossref.org/contributor/thaha-muhammed-19gce4u4co3gj">
        <j.3:name>Thaha Muhammed</j.3:name>
        <j.3:familyName>Muhammed</j.3:familyName>
        <j.3:givenName>Thaha</j.3:givenName>
    </j.3:Person>
</j.0:creator>
</rdf:Description>
        </rdf:RDF>
```

Fetch the reference of a publication with DOI: 10.1016/j.jnca.2016.10.019 with verbose output, 
original key format and unabbrevated journal name.
```
$ bibrefer ref -kjV 10.1016/j.jnca.2016.10.019

Provided valid DOI: 10.1016/j.jnca.2016.10.01919 
Refernce for DOI: 10.1016/j.jnca.2016.10.019 obtained
Cleaning the obtained reference
Abbreviating the author names
Removing url and month

@article{Muhammed_2017,
        doi = {10.1016/j.jnca.2016.10.019},
        year = 2017,
        publisher = {Elsevier {BV}},
        volume = {78},
        pages = {267--287},
        author = {T. Muhammed and R. A. Shaikh},
        title = {An analysis of fault detection strategies in wireless sensor networks},
        journal = {Journal of Network and Computer Applications}
}
```

### Additional Documentation
* [bibrefer doi](./doc/markdown/bibrefer_doi.md)	 - Returns the DOI for a given publication name or pdf file
* [bibrefer ref](./doc/markdown/bibrefer_ref.md)	 - Returns the reference of the given DOI, name, or pdf file

### TODO
- [ ] Clipping functionality is missing
- [ ] Add support for direct fetching with publication name/pdf file
- [ ] Improve documentation
