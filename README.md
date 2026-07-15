# URLHaus Go Wrapper

[![Go Report Card](https://goreportcard.com/badge/github.com/vertoforce/gourlhaus)](https://goreportcard.com/report/github.com/vertoforce/gourlhaus)
[![Documentation](https://pkg.go.dev/badge/github.com/vertoforce/gourlhaus.svg)](https://pkg.go.dev/github.com/vertoforce/gourlhaus)

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

These bulk CSV download functions (`GetRecentURLs`, `GetAllURLs`, `GetAllOnlineURLs`, `FillInURLHashDetails`) currently work without authentication.

### Submitting URLs

`SubmitURLs` requires an abuse.ch Auth-Key (free from <https://auth.abuse.ch/>). Pass it as the `apiKey` argument; it is sent as the `Auth-Key` HTTP header.

```go
urlhaus.SubmitURLs(ctx, urls, authKey, tags, "malware_download")
```
