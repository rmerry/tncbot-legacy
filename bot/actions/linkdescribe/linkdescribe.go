package linkdescribe

import (
	"net/http"

	m "github.com/keighl/metabolize"
	"mvdan.cc/xurls"
)

// MetaData ...
type MetaData struct {
	Title       string `meta:"og:title,title"`
	Description string `meta:"og:description,description"`
}

// Execute ...
func Execute(input string) (string, bool) {
	url := xurls.Relaxed().FindString(input)
	if url == "" {
		return "", false
	}

	res, _ := http.Get(url)

	data := new(MetaData)

	err := m.Metabolize(res.Body, data)
	if err != nil || data == nil {
		return "", false
	}

	var linkDescriptor string

	if data.Title != "" {
		linkDescriptor = data.Title
	}
	if data.Description != "" {
		linkDescriptor += " - " + data.Description
	}

	if linkDescriptor != "" {
		linkDescriptor = "* " + linkDescriptor
	}

	return linkDescriptor, true
}
