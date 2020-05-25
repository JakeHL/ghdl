package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ezPanic(err error) {
	if err != nil {
		fmt.Print("There was an error fetching data from github\n")
		os.Exit(1)
	}
}

type releaseAsset struct {
	DownloadCount int64 `json:"download_count"`
}

type releaseResponse struct {
	Assets []releaseAsset `json:"assets"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("User must be specified\n")
		os.Exit(1)
	}
	user := os.Args[1]

	if len(os.Args) < 3 {
		fmt.Print("Repository must be specified\n")
		os.Exit(1)
	}
	repo := os.Args[2]

	repoAPIURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", user, repo)

	fmt.Printf("Fetching release downloads for https://github.com/%v/%v\n", user, repo)

	response, err := http.Get(repoAPIURL)
	ezPanic(err)
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

	fmt.Printf("Assets have been downloaded %v times", totalDownloads)
}
