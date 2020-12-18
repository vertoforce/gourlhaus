package gourlhaus

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"golang.org/x/net/context/ctxhttp"
)

// Internal constants
const (
	urlHausRecentURLsLink    = "https://urlhaus.abuse.ch/downloads/csv_recent/"
	urlHausAllURLsLink       = "https://urlhaus.abuse.ch/downloads/csv/"
	urlHausAllOnlineURLsLink = "https://urlhaus.abuse.ch/downloads/csv_online/"
	urlHausPayloadsURL       = "https://urlhaus.abuse.ch/downloads/payloads/"
	urlHausURLSubmit         = "https://urlhaus.abuse.ch/api/"
)

// URLStatus State of URL
type URLStatus int

// Submission struct of the JSON required to submit a URL
type Submission struct {
	Token      string          `json:"token"`
	Anonymous  string          `json:"anonymous"`
	Submission []submissionURL `json:"submission"`
}

type submissionURL struct {
	URL    string   `json:"url"`
	Threat string   `json:"threat"`
	Tags   []string `json:"tags,omitempty"`
}

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
	URLHashes []URLHashDetails `csv:"-"`
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

// SubmitURLs takes a list of URLs and attempsts to submit them. The list returned wil contain the URLs that were successfully submitted
func SubmitURLs(ctx context.Context, urls []string, apiKey string, tags []string, threat string) (io.Reader, error) {
	submission := &Submission{}
	submission.Token = apiKey
	submission.Anonymous = "0"
	submission.Submission = []submissionURL{}

	for _, url := range urls {
		if url == "" {
			continue
		}

		urlEntry := &submissionURL{}
		urlEntry.URL = url
		urlEntry.Tags = tags
		urlEntry.Threat = threat
		submission.Submission = append(submission.Submission, *urlEntry)
	}
	jsonEntries, err := json.Marshal(submission)
	if err != nil {
		return nil, err
	}

	httpBody := bytes.NewBuffer(jsonEntries)
	resp, err := ctxhttp.Post(ctx, nil, urlHausURLSubmit, "application/json", httpBody)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

//CheckForUnseenURLs takes a list of URLs and returns the ones that havent been submitted to the platform
func CheckForUnseenURLs(ctx context.Context, urls []string) ([]string, error) {
	entries, err := GetAllURLs(ctx)
	if err != nil {
		return nil, err
	}
	unseenURLs := []string{}
outerLoop:
	for _, url := range urls {
		seen := false
		for _, entry := range entries {
			if url == entry.URL {
				unseenURLs = append(unseenURLs, url)
				continue outerLoop
			}
		}
	}

	return unseenURLs, nil
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
			urlEntries[i].URLHashes = append(urlEntries[i].URLHashes, payload)
		}
	}

	return nil
}
