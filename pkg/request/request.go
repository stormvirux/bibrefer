package request

import (
	"fmt"
	"github.com/goccy/go-json"
	"io"
	"net/http"
	"strings"
)

func DoiDataCite(host string, query string) (string, error) {
	// host = https://api.datacite.org/dois
	urlBuilder := &strings.Builder{}
	urlBuilder.WriteString(host)
	header := make(map[string]string)
	header["Accept"] = "application/citeproc+json"
	header["Content-Type"] = "application/json"
	req, err := bibRequest("GET", urlBuilder.String(), nil, header)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	client := &http.Client{}
	// verbosePrint(verbose, fmt.Sprintf("Getting reference for DOI: %s from host: %s", strippedDoi, host))
	urlQuery := make(map[string]string)
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
		return "", fmt.Errorf("could not find any article with that name")
	}
	return data[0].DOI, err
}

func DoiCrossRef(host string, query string) (string, error) {
	urlBuilder := &strings.Builder{}
	urlBuilder.WriteString(host)
	header := make(map[string]string)
	header["Accept"] = "application/json"
	header["Content-Type"] = "application/json"

	req, err := bibRequest("GET", urlBuilder.String(), nil, header)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	client := &http.Client{}
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

func bibRequest(method, path string, body io.Reader, query map[string]string) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	for k, v := range query {
		req.Header.Add(k, v)
	}
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
		return nil, err
	}
	// Any common handling of response
	return res, nil
}
