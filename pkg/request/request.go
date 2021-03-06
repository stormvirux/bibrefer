/*
Copyright © 2022 Thaha Mohammed <thaha.mohammed@aalto.fi>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package request

import (
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
	"time"
)

const appJSON = "application/json"

// DoiDataCite accepts the title name as a query and returns the DOI as string from DataCite API along with error.
func DoiDataCite(query string) (string, error) {
	var host = "https://api.datacite.org/dois"
	urlBuilder := &strings.Builder{}
	urlBuilder.WriteString(host)
	header := make(map[string]string)
	header["Accept"] = "application/citeproc+json"
	header["Content-Type"] = appJSON
	req, err := bibRequest("GET", urlBuilder.String(), nil, header)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	client := &http.Client{Timeout: 100 * time.Second}
	// verbosePrint(verbose, fmt.Sprintf("Getting reference for DOI: %s from host: %s", strippedDoi, host))
	urlQuery := make(map[string]string)
	query = strings.ReplaceAll(query, ":", " ")
	query = strings.ReplaceAll(query, "{", "")
	query = strings.ReplaceAll(query, "}", "")
	urlQuery["query"] = "titles.title:" + query
	res, err := bibDo(client, req, urlQuery)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	var data Doi
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if len(data) < 1 {
		return "", nil
	}
	return data[0].DOI, err
}

// DoiCrossRef accepts the title name as a query and returns the DOI as string from CrossRef API along with error.
func DoiCrossRef(query string) (string, error) {
	var host = "https://api.crossref.org/works"
	urlBuilder := &strings.Builder{}
	urlBuilder.WriteString(host)
	header := make(map[string]string)
	header["Accept"] = appJSON
	header["Content-Type"] = appJSON

	req, err := bibRequest("GET", urlBuilder.String(), nil, header)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	client := &http.Client{Timeout: 100 * time.Second}
	// verbosePrint(verbose, fmt.Sprintf("Getting reference for DOI: %s from host: %s", strippedDoi, host))
	urlQuery := make(map[string]string)
	urlQuery["rows"] = "1"
	urlQuery["select"] = "DOI"
	urlQuery["query.title"] = query
	res, err := bibDo(client, req, urlQuery)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	var data queryJSON
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if len(data.Message.Items) == 0 {
		return "", fmt.Errorf("could not find any article with that name")
	}
	return data.Message.Items[0].DOI, err
}

// RefDoi accepts the DOI as a query and the output format and returns the reference as string from doi.org API along with error.
// The output format can be json, xml or bibtex.
func RefDoi(query string, output string) (string, error) {
	var host = "https://doi.org/"
	urlBuilder := &strings.Builder{}
	urlBuilder.WriteString(host)
	urlBuilder.WriteString(query)
	header := make(map[string]string)
	header["Accept"] = accept(output)
	req, err := bibRequest("GET", urlBuilder.String(), nil, header)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	client := &http.Client{Timeout: 100 * time.Second}
	// verbosePrint(verbose, fmt.Sprintf("Getting reference for DOI: %s from host: %s", query, host), os.Stdout)
	res, err := bibDo(client, req, map[string]string{})
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = Body.Close()
		}
	}(res.Body)
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return string(body), err
}

func accept(output string) string {
	switch output {
	case "json":
		return "application/citeproc+json"
	case "xml":
		return "application/rdf+xml"
	}
	return "application/x-bibtex"
}

func bibRequest(method, path string, body io.Reader, query map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range query {
		req.Header.Add(k, v)
	}
	var email = "bibrefer@gmail.com"
	req.Header.Add("User-Agent", fmt.Sprintf("bibrefer, (github.com/stormvirux/bibrefer/; mailto: %s)", email))

	return req, nil
}

func bibDo(client *http.Client, req *http.Request, query map[string]string) (*http.Response, error) {
	q := req.URL.Query()

	for k, v := range query {
		q.Add(k, v)
	}

	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("[Error] %v\n", err)
		return nil, err
	}
	return res, nil
}
