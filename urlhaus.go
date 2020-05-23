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
	urlHashes []URLHashDetails `csv:"-"`
}

// URLHashDetails about the hash hosted at this url
type URLHashDetails struct {
	FirstSeen string `csv:"firstseen"`
	URL       string `csv:"url"`
	Filetype  string `csv:"filetype"`
	MD5       string `csv:"md5"`
	SHA256    string `csv:"sha256"`
	Signature string `csv:"signature"`
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
	err := downloadCSV(ctx, urlHausAllOnlineURLsLink, false, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// FillInURLHashDetails Downloads the payload data from URL haus and adds in the hashes to the provided urlEntries slice
// Note that this function can quickly bring the memory used by []URLEntry to > 1GB
func FillInURLHashDetails(ctx context.Context, urlEntries []URLEntry) error {
	payloads := []URLHashDetails{}
	err := downloadCSV(ctx, urlHausPayloadsURL, true, &payloads)
	if err != nil {
		return err
	}

	// Create map of urls to the index in the list
	// This helps in the speed of finding urlentries looping through the payloads array
	urlEntriesMap := map[string]int{}
	for i, entry := range urlEntries {
		urlEntriesMap[entry.URL] = i
	}

	for _, payload := range payloads {
		// Find the url entry with this link
		if i, ok := urlEntriesMap[payload.URL]; ok {
			// Remove url from payload to save space
			payload.URL = ""
			urlEntries[i].urlHashes = append(urlEntries[i].urlHashes, payload)
		}
	}

	return nil
}
