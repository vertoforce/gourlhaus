package gourlhaus

import (
	"net/http"

	"github.com/smartystreets/scanners/csv"
)

const (
	urlHausHashesLink = "https://urlhaus.abuse.ch/downloads/payloads/"
)

// PopulateURLEntriesWithHashes Given a URLEntry list, populate the entries with hash information from urlHaus
func PopulateURLEntriesWithHashes(entries URLEntries) error {
	response, err := http.Get(urlHausHashesLink)
	if err != nil {
		return err
	}

	// Get rows
	scanner := csv.NewScanner(response.Body, csv.Comma(','), csv.Comment('#'), csv.ContinueOnError(true))
	for scanner.Scan() {
		if err := scanner.Error(); err != nil {
			return nil
		}
		fields := scanner.Record()
		// Header: firstseen,url,filetype,md5,sha256,signature

		// Add hashes to the found entry based on url
		if _, ok := entries[fields[1]]; ok {
			entries[fields[1]].Filetype = fields[2]
			entries[fields[1]].MD5 = fields[3]
			entries[fields[1]].SHA256 = fields[4]
			// TODO: Add signature?
		}
	}

	return nil
}
