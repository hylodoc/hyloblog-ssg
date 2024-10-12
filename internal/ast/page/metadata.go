package page

import (
	"fmt"
	"net/url"
	"time"

	"gopkg.in/yaml.v3"
)

type metadata struct {
	URL       string       `yaml:"url"`
	Published parsabletime `yaml:"published"`
	Updated   parsabletime `yaml:"updated"`
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

func (m *metadata) timing() *timing {
	published, updated := time.Time(m.Published), time.Time(m.Updated)
	if published.IsZero() {
		return nil
	}
	if updated.IsZero() {
		return &timing{published, published}
	}
	return &timing{published, updated}
}

type parsabletime time.Time

func (pt *parsabletime) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}
	for _, format := range []string{
		"2006-01-02",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04:05 -0700",

		"Jan 02, 2006",
		"Jan 02, 2006 15:04",
		"Jan 02, 2006 15:04:05",

		"Jan 2, 2006",
		"Jan 2, 2006 15:04",
		"Jan 2, 2006 15:04:05",

		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z07:00",

		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
	} {
		t, err := time.Parse(format, raw)
		if err == nil {
			*pt = parsabletime(t)
			return nil
		}
	}
	return fmt.Errorf("unable to parse date: %s", raw)
}
