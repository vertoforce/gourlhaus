# URLHaus Go Wrapper

Simple library to fetch URLHaus URLs in go

## Usage

```go
urlhaus.GetRecentURLs()
urlhaus.GetAllURLs()
urlhaus.GetAllOnlineURLs()
```

Each returns `[]URLEntry` with `URLEntry` defined as:

```go
type URLEntry struct {
    ID          string
    DateAdded   string
    URL         string
    URLStatus   URLStatus
    Threat      string
    Tags        []string
    URLHausLink string
    Reporter    string
}
```
