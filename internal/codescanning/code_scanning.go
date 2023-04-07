package codescanning

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/cli/go-gh"
)

type Alert struct {
	Number             int                `json:"number"`
	CreatedAt          time.Time          `json:"created_at"`
	URL                string             `json:"url"`
	HTMLURL            string             `json:"html_url"`
	State              string             `json:"state"`
	DismissedBy        *DismissedBy       `json:"dismissed_by"`
	DismissedAt        *time.Time         `json:"dismissed_at"`
	DismissedReason    string             `json:"dismissed_reason"`
	DismissedComment   string             `json:"dismissed_comment"`
	Rule               Rule               `json:"rule"`
	Tool               Tool               `json:"tool"`
	MostRecentInstance MostRecentInstance `json:"most_recent_instance"`
	InstancesURL       string             `json:"instances_url"`
	Repository         Repository         `json:"repository"`
}

type DismissedBy struct {
	Login string `json:"login"`
}

type Rule struct {
	ID          string   `json:"id"`
	Severity    string   `json:"severity"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
}

type Tool struct {
	Name    string `json:"name"`
	GUID    string `json:"guid"`
	Version string `json:"version"`
}

type MostRecentInstance struct {
	Ref             string   `json:"ref"`
	AnalysisKey     string   `json:"analysis_key"`
	Category        string   `json:"category"`
	Environment     string   `json:"environment"`
	State           string   `json:"state"`
	CommitSHA       string   `json:"commit_sha"`
	Message         Message  `json:"message"`
	Location        Location `json:"location"`
	Classifications []string `json:"classifications"`
}

type Message struct {
	Text string `json:"text"`
}

type Location struct {
	Path        string `json:"path"`
	StartLine   int    `json:"start_line"`
	EndLine     int    `json:"end_line"`
	StartColumn int    `json:"start_column"`
	EndColumn   int    `json:"end_column"`
}

type Repository struct {
	ID       int64  `json:"id"`
	NodeID   string `json:"node_id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    Owner  `json:"owner"`
}

type Owner struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	SiteAdmin bool   `json:"site_admin"`
}

func GetCodeScanningAlerts(org string) []Alert {
	var linkRE = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)
	findNextPage := func(response *http.Response) (string, bool) {
		for _, m := range linkRE.FindAllStringSubmatch(response.Header.Get("Link"), -1) {
			if len(m) > 2 && m[2] == "next" {
				return m[1], true
			}
		}
		return "", false
	}

	client, err := gh.RESTClient(nil)
	if err != nil {
		log.Fatal(err)
	}

	requestPath := fmt.Sprintf("orgs/%s/code-scanning/alerts?per_page=100", org)
	alerts := []Alert{}
	page := 1
	for {
		response, err := client.Request(http.MethodGet, requestPath, nil)
		if err != nil {
			log.Fatal(err)
		}

		a := []Alert{}
		decoder := json.NewDecoder(response.Body)
		err = decoder.Decode(&a)
		if err != nil {
			log.Fatal(err)
		}
		if err := response.Body.Close(); err != nil {
			log.Fatal(err)
		}

		for _, alert := range a {
			alerts = append(alerts, alert)
		}

		var hasNextPage bool
		if requestPath, hasNextPage = findNextPage(response); !hasNextPage {
			break
		}

		page++

	}

	return alerts

}
