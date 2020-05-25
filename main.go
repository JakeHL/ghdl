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
	Owner string `arg:"positional,required" help:"Repository owner"`
	Repo  string `arg:"positional,required" help:"Repository name"`
	Terse bool   `arg:"-t" help:"Minimal output mode"`
}

type releaseAsset struct {
	DownloadCount int64 `json:"download_count"`
}

type releaseResponse struct {
	Assets []releaseAsset `json:"assets"`
}

func main() {
	var args arguments
	arg.MustParse(&args)

	repoAPIURL := fmt.Sprintf("https://api.github.com/repos/%v/%v/releases", args.Owner, args.Repo)

	if !args.Terse {
		fmt.Printf("Fetching release downloads for https://github.com/%v/%v\n", args.Owner, args.Repo)
	}

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

	if args.Terse {
		fmt.Print(totalDownloads)
	} else {
		fmt.Printf("Assets have been downloaded %v times", totalDownloads)
	}
}
