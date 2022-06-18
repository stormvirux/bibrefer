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

type respCrossRefJSON struct {
	Query        map[string]interface{} `json:"query"`
	Facets       map[string]interface{} `json:"facets"`
	TotalResults int                    `json:"total-results"`
	Items        []*work                `json:"items"`
	ItemsPerPage int                    `json:"items-per-page"`
}

type queryJSON struct {
	Status  string            `json:"status"`
	Type    string            `json:"message-type"`
	Version string            `json:"message-version"`
	Message *respCrossRefJSON `json:"message"`
}

type work struct {
	DOI string `json:"DOI"`
}

type respDataCite struct {
	DOI string `json:"DOI"`
}

type Doi []respDataCite
