package gourlhaus

import (
	"context"
	"net/http"

	"github.com/smartystreets/scanners/csv"
)

const (
	urlHausHashesLink = "https://urlhaus.abuse.ch/downloads/payloads/"
)

// GetURLsWithHashes Get URLDetails with hash information
func GetURLsWithHashes(ctx context.Context) (chan URLDetails, error) {
	response, err := http.Get(urlHausHashesLink)
	if err != nil {
		return nil, err
	}

	URLEntries := make(chan URLDetails)

	// Initialize scanner
	scanner := csv.NewScanner(response.Body, csv.Comma(','), csv.Comment('#'), csv.ContinueOnError(true))

	go func() {
		defer close(URLEntries)

		// Parse each row
		for scanner.Scan() {
			if err := scanner.Error(); err != nil {
				return
			}
			// Header: firstseen,url,filetype,md5,sha256,signature
			fields := scanner.Record()
			if len(fields) < 6 {
				continue
			}

			entry := URLDetails{}
			entry.FirstSeen = fields[0]
			entry.URL = fields[1]
			entry.Filetype = fields[2]
			entry.MD5 = fields[2]
			entry.SHA256 = fields[3]
			entry.Signature = fields[4]

			// Add to channel or cancel
			select {
			case URLEntries <- entry:
			case <-ctx.Done():
				return
			}
		}
	}()

	return URLEntries, nil
}
