package gourlhaus

import (
	"context"
)

// Internal constants
const (
	urlHausRecentURLsLink    = "https://urlhaus.abuse.ch/downloads/csv_recent/"
	urlHausAllURLsLink       = "https://urlhaus.abuse.ch/downloads/csv/"
	urlHausAllOnlineURLsLink = "https://urlhaus.abuse.ch/downloads/csv_online/"
	urlHausPayloadsURL       = "https://urlhaus.abuse.ch/downloads/payloads/"
)

// URLStatus State of URL
type URLStatus int

// URLEntry Entry of URL in URLHaus
type URLEntry struct {
	ID          string `csv:"id"`
	DateAdded   string `csv:"dateadded"`
	URL         string `csv:"url"`
	URLStatus   string `csv:"url_status"`
	Threat      string `csv:"threat"`
	Tags        string `csv:"tags"`
	URLHausLink string `csv:"urlhaus_link"`
	Reporter    string `csv:"reporter"`

	// Hash details populated after calling FillInURLHashDetails
	// URLHashDetails `csv:"-"`
}

// URLHashDetails about the hash hosted at this url
type URLHashDetails struct {
	FirstSeen string
	Filetype  string
	MD5       string
	SHA256    string
	Signature string
}

// GetRecentURLs Get all recent urls from URLHaus
func GetRecentURLs(ctx context.Context) ([]URLEntry, error) {
	ret := []URLEntry{}
	err := downloadCSV(ctx, urlHausRecentURLsLink, false, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetAllURLs Get all urlhaus urls
func GetAllURLs(ctx context.Context) ([]URLEntry, error) {
	ret := []URLEntry{}
	err := downloadCSV(ctx, urlHausAllURLsLink, true, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetAllOnlineURLs Get all urlhaus urls
func GetAllOnlineURLs(ctx context.Context) ([]URLEntry, error) {
	ret := []URLEntry{}
	err := downloadCSV(ctx, urlHausAllOnlineURLsLink, true, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// TODO:
// func FillInURLHashDetails(ctx context.Context, urlEntries []URLEntry) error {

// }
