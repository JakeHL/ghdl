package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func ezPanic(err error) {
	if err != nil {
		log.Fatal("There was an error fetching data from github")
	}
}

type releaseAsset struct {
	DownloadCount int64 `json:"download_count"`
}

type releaseResponse struct {
	Assets []releaseAsset `json:"assets"`
}

func main() {
	user := os.Args[1]
	if user == "" {
		log.Fatal("User must be specified")
		return
	}

	repo := os.Args[2]
	if repo == "" {
		log.Fatal("Repository must be specified")
		return
	}

	repoAPIURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", user, repo)

	fmt.Printf("Fethcing release downloads for https://github.com/%v/%v\n", user, repo)

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
