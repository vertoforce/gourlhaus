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

### Getting hashes of the content hosted at those URLs

URLHaus also provides the hashes that were found on the urls.  However it's a separate endpoint, so to populate the URLEntries with the hash data, call `FillInURLHashDetails()` on the URLEntry list.
