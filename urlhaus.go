package gourlhaus

import (
	"io"
	"net/http"
	"strings"

	"github.com/smartystreets/scanners/csv"
)

// Internal constants
const (
	urlHausRecentURLsLink    = "https://urlhaus.abuse.ch/downloads/csv_recent/"
	urlHausAllURLsLink       = "https://urlhaus.abuse.ch/downloads/csv/"
	urlHausAllOnlineURLsLink = "https://urlhaus.abuse.ch/downloads/csv_online/"
)

// URLStatus State of URL
type URLStatus int

// URL State
const (
	URLOnline URLStatus = iota
	URLOffline
	URLUnknown
)

// URLEntries Map of url to details about that url
type URLEntries map[string]*URLDetails

// URLDetails Entry of URL in URLHaus
type URLDetails struct {
	ID          string
	DateAdded   string
	URLStatus   URLStatus
	Threat      string
	Tags        []string
	URLHausLink string
	Reporter    string

	// Default to nothing, but filled in when calling GetURLPayloadHashes
	Filetype string
	MD5      string
	SHA256   string
}

// GetRecentURLs Get all recent urls from URLHaus
func GetRecentURLs() (URLEntries, error) {
	return linkToURLEntries(urlHausRecentURLsLink)
}

// GetAllURLs Get all urlhaus urls
func GetAllURLs() (URLEntries, error) {
	return linkToURLEntries(urlHausAllURLsLink)
}

// GetAllOnlineURLs Get all urlhaus urls
func GetAllOnlineURLs() (URLEntries, error) {
	return linkToURLEntries(urlHausAllOnlineURLsLink)
}

func linkToURLEntries(url string) (URLEntries, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return csvToURLEntries(response.Body), nil
}

// csvToURLEntries Given a reader with csv content, convert to URLEntries
func csvToURLEntries(reader io.Reader) URLEntries {
	scanner := csv.NewScanner(reader, csv.Comma(','), csv.Comment('#'), csv.ContinueOnError(true))
	URLEntries := URLEntries{}

	for scanner.Scan() {
		// Current header: id,dateadded,url,url_status,threat,tags,urlhaus_link,reporter
		if err := scanner.Error(); err != nil {
			return nil
		}
		fields := scanner.Record()

		// Craft URLEntry
		URLDetails := URLDetails{}
		URLDetails.ID = fields[0]
		URLDetails.DateAdded = fields[1]
		if fields[3] == "offline" {
			URLDetails.URLStatus = URLOffline
		} else if fields[3] == "online" {
			URLDetails.URLStatus = URLOnline
		} else {
			URLDetails.URLStatus = URLUnknown
		}
		URLDetails.Threat = fields[4]
		if fields[5] != "None" {
			URLDetails.Tags = strings.Split(fields[5], ",")
		}
		URLDetails.URLHausLink = fields[6]
		URLDetails.Reporter = fields[7]

		// Add to slice
		URLEntries[fields[2]] = &URLDetails
	}

	return URLEntries
}
