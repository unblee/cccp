package main

import (
	"context"
	"log"
	"path"

	"github.com/unblee/cccp"
)

func main() {
	urls := []string{
		"https://dl.google.com/go/go1.11.linux-amd64.tar.gz",
		"https://dl.google.com/go/go1.11.windows-amd64.msi",
		"https://dl.google.com/go/go1.11.darwin-amd64.pkg",
	}

	for _, url := range urls {
		err := cccp.SetFromURLToFile(url, path.Base(url), "")
		if err != nil {
			log.Fatalln(err)
		}
	}

	cccp.SetOptions(cccp.WithConcurrentNumCPU())
	// cccp.SetOptions(cccp.WithConcurrent(1), cccp.WithEnableSequentialProgressbars())
	err := cccp.Run(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
