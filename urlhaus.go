package urlhaus

import (
	"io"
	"net/http"
	"strings"

	"github.com/smartystreets/scanners/csv"
)

// Internal constants
const (
	urlHausRecentURLsLink = "https://urlhaus.abuse.ch/downloads/csv_recent/"
	urlHausAllURLsLink    = "https://urlhaus.abuse.ch/downloads/csv/"
)

// URLStatus State of URL
type URLStatus int

// URL State
const (
	URLOnline URLStatus = iota
	URLOffline
	URLUnknown
)

// URLEntry Entry of URL in URLHaus
type URLEntry struct {
	ID          string
	DateAdded   string
	URL         string
	URLStatus   URLStatus
	Threat      string
	Tags        []string
	URLHausLink string
	Reporter    string
}

// GetRecentURLs Get all recent urls from URLHaus
func GetRecentURLs() ([]URLEntry, error) {
	return linkToURLEntries(urlHausRecentURLsLink)
}

// GetAllURLs Get all urlhaus urls
func GetAllURLs() ([]URLEntry, error) {
	return linkToURLEntries(urlHausAllURLsLink)
}

func linkToURLEntries(url string) ([]URLEntry, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return csvToURLEntries(response.Body), nil
}

func csvToURLEntries(reader io.Reader) []URLEntry {
	scanner := csv.NewScanner(reader, csv.Comma(','), csv.Comment('#'), csv.ContinueOnError(true))
	URLEntries := []URLEntry{}

	for scanner.Scan() {
		// Current header: id,dateadded,url,url_status,threat,tags,urlhaus_link,reporter
		if err := scanner.Error(); err != nil {
			return nil
		}
		fields := scanner.Record()

		// Craft URLEntry
		URLEntry := URLEntry{}
		URLEntry.ID = fields[0]
		URLEntry.DateAdded = fields[1]
		URLEntry.URL = fields[2]
		if fields[3] == "offline" {
			URLEntry.URLStatus = URLOffline
		} else if fields[3] == "online" {
			URLEntry.URLStatus = URLOnline
		} else {
			URLEntry.URLStatus = URLUnknown
		}
		URLEntry.Threat = fields[4]
		if fields[5] != "None" {
			URLEntry.Tags = strings.Split(fields[5], ",")
		}
		URLEntry.URLHausLink = fields[6]
		URLEntry.Reporter = fields[7]

		// Add to slice
		URLEntries = append(URLEntries, URLEntry)
	}

	return URLEntries
}
