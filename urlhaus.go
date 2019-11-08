package gourlhaus

import (
	"context"
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

// URLDetails Entry of URL in URLHaus
type URLDetails struct {
	ID          string
	DateAdded   string
	URL         string
	URLStatus   URLStatus
	Threat      string
	Tags        []string
	URLHausLink string
	Reporter    string

	// Hash details from GetURLsWithHashes
	URLHashDetails
}

// URLHashDetails Details only filled by calling GetURLsWithHashes
type URLHashDetails struct {
	FirstSeen string
	Filetype  string
	MD5       string
	SHA256    string
	Signature string
}

// GetRecentURLs Get all recent urls from URLHaus
func GetRecentURLs(ctx context.Context) (chan URLDetails, error) {
	return linkToURLEntries(ctx, urlHausRecentURLsLink)
}

// GetAllURLs Get all urlhaus urls
func GetAllURLs(ctx context.Context) (chan URLDetails, error) {
	return linkToURLEntries(ctx, urlHausAllURLsLink)
}

// GetAllOnlineURLs Get all urlhaus urls
func GetAllOnlineURLs(ctx context.Context) (chan URLDetails, error) {
	return linkToURLEntries(ctx, urlHausAllOnlineURLsLink)
}

func linkToURLEntries(ctx context.Context, url string) (chan URLDetails, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return csvToURLEntries(ctx, response.Body), nil
}

// csvToURLEntries Given a reader with csv content, convert to URLEntries
func csvToURLEntries(ctx context.Context, reader io.Reader) chan URLDetails {
	URLEntries := make(chan URLDetails)

	// Initialize scanner
	scanner := csv.NewScanner(reader, csv.Comma(','), csv.Comment('#'), csv.ContinueOnError(true))

	go func() {
		defer close(URLEntries)

		// Parse each row
		for scanner.Scan() {
			// Current header: id,dateadded,url,url_status,threat,tags,urlhaus_link,reporter
			if err := scanner.Error(); err != nil {
				return
			}
			fields := scanner.Record()

			// Craft URLEntry
			URLDetails := URLDetails{}
			URLDetails.ID = fields[0]
			URLDetails.DateAdded = fields[1]
			URLDetails.URL = fields[2]
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

			// Send to channel or cancel
			select {
			case URLEntries <- URLDetails:
			case <-ctx.Done():
				return
			}
		}
	}()

	return URLEntries
}
