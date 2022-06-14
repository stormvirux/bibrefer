package getdoi

type respCrossRefJSON struct {
	Query        map[string]interface{} `json:"query"`
	Facets       map[string]interface{} `json:"facets"`
	TotalResults int                    `json:"total-results"`
	Items        []*Work                `json:"items"`
	ItemsPerPage int                    `json:"items-per-page"`
}

type queryJSON struct {
	Status  string            `json:"status"`
	Type    string            `json:"message-type"`
	Version string            `json:"message-version"`
	Message *respCrossRefJSON `json:"message"`
}

type Work struct {
	DOI string `json:"DOI"`
}

type respDataCite struct {
	Doi string `json:"DOI"`
}

type Doi []respDataCite
