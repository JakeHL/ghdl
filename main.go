package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alexflint/go-arg"
)

func ezPanic(err error) {
	if err != nil {
		fmt.Print("There was an error fetching data from github\n")
		os.Exit(1)
	}
}

type arguments struct {
	Owner    string `arg:"positional,required" help:"Repository owner"`
	Repo     string `arg:"positional,required" help:"Repository name"`
	Username string `arg:"-u" help:"Username for authenticated requests"`
	Password string `arg:"-p" help:"Password for authenticated requests, if you have 2fa, use a Personal Access Token"`
	OAuth    string `arg:"-a" help:"OAuth token for authenticated requests, you can also use a Personal Access Token"`
	Terse    bool   `arg:"-t" help:"Minimal output mode"`
}

func (arguments) Version() string {
	return "ghdl 1.0.0"
}

type releaseAsset struct {
	DownloadCount int64 `json:"download_count"`
}

type releaseResponse struct {
	Assets []releaseAsset `json:"assets"`
}

func main() {
	var args arguments
	parsed := arg.MustParse(&args)

	if (args.OAuth != "" && args.Username != "") ||
		(args.OAuth != "" && args.Password != "") {
		parsed.Fail("Use OAauth OR basic Authentication, not both.")
	}

	if (args.Username != "" && args.Password == "") ||
		(args.Password != "" && args.Username == "") {
		parsed.Fail("Both username and password must be provided when making authenticated requests")
	}

	repoAPIURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", args.Owner, args.Repo)

	if !args.Terse {
		fmt.Printf("Fetching release downloads for https://github.com/%v/%v\n", args.Owner, args.Repo)
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", repoAPIURL, nil)
	if args.Username != "" && args.Password != "" {
		req.SetBasicAuth(args.Username, args.Password)
	} else if args.OAuth != "" {
		req.Header.Set("Authorization", "token "+args.OAuth)
	}
	response, err := client.Do(req)
	ezPanic(err)

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		fmt.Printf("Failed to fetch release information. Statuscode: %v\n", response.StatusCode)
		os.Exit(1)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	ezPanic(err)

	var releases []releaseResponse
	json.Unmarshal(body, &releases)

	var totalDownloads int64 = 0
	for _, release := range releases {
		for _, asset := range release.Assets {
			totalDownloads += asset.DownloadCount
		}
	}

	if args.Terse {
		fmt.Print(totalDownloads)
	} else {
		fmt.Printf("Assets have been downloaded %v times", totalDownloads)
	}
}
