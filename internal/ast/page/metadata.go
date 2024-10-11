package page

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v3"
)

type metadata struct {
	URL string `yaml:"url"`
}

func parsemetadata(raw string) (*metadata, error) {
	var m metadata
	if err := yaml.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("cannot unmarshal: %w", err)
	}
	if err := confirmurlvalid(m.URL); err != nil {
		return nil, fmt.Errorf("url error: %w", err)
	}
	return &m, nil
}

func confirmurlvalid(u string) error {
	if len(u) > 0 && u[0] != '/' {
		return fmt.Errorf("must begin with '/'")
	}
	if _, err := url.Parse(u); err != nil {
		return fmt.Errorf("cannot parse: %w", err)
	}
	return nil
}
