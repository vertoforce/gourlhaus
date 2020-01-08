# URLHaus Go Wrapper

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/gourlhaus)](https://goreportcard.com/report/github.com/vertoforce/gourlhaus)
[![Documentation](https://godoc.org/github.com/vertoforce/gourlhaus?status.svg)](https://godoc.org/github.com/vertoforce/gourlhaus)

Simple library to fetch URLHaus URLs in go

## Usage

### Getting URLs

```go
urlhaus.GetRecentURLs()
urlhaus.GetAllURLs()
urlhaus.GetAllOnlineURLs()
```

Each returns `URLEntries` which is a map of the url to details about it `map[string]URLDetails`

```go
type URLDetails struct {
    ID          string
    DateAdded   string
    URL         string
    URLStatus   URLStatus
    Threat      string
    Tags        []string
    URLHausLink string
    Reporter    string

    // Default to nothing, but filled in when calling PopulateURLEntriesWithHashes
    Filetype string
    MD5      string
    SHA256   string
}
```

### Getting hashes of the content hosted at those URLs

Calling `PopulateURLEntriesWithHashes(entries URLEntries)` populates the `URLEntries` with hashes of the content hosted there if it's found in URLHaus (separate endpoint).
