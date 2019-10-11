package urlhaus

import "testing"

func TestGetURLs(t *testing.T) {
	// Recent
	recentURLs, err := GetRecentURLs()
	if err != nil {
		t.Errorf("Error fetching recent URLs: " + err.Error())
	}
	if len(recentURLs) == 0 {
		t.Errorf("GetRecentURLs returned no URLs")
	}

	// Online
	onlineURLs, err := GetAllOnlineURLs()
	if err != nil {
		t.Errorf("Error fetching online URLs: " + err.Error())
	}
	if len(onlineURLs) == 0 {
		t.Errorf("GetAllOnlineURLs returned no URLs")
	}

	// All
	allURLs, err := GetAllURLs()
	if err != nil {
		t.Errorf("GetAllURLs returned error: " + err.Error())
	}
	if len(allURLs) == 0 {
		t.Errorf("GetAllURLs returned no URLs")
	}
}
