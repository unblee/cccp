# cccp

[![Build Status](https://img.shields.io/travis-ci/unblee/cccp.svg?style=flat-square)][travis]
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)][godoc]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[travis]: https://travis-ci.org/unblee/cccp
[godoc]: https://godoc.org/github.com/unblee/cccp
[license]: https://github.com/unblee/cccp/blob/master/LICENSE

## Description

A library that provides a concurrent copy function with progress bars.

See [godoc](https://godoc.org/github.com/unblee/cccp) for complete documentation.

## Screenshot

![Screenshot 01](assets/screenshot-01.gif)

![Screenshot 02](assets/screenshot-02.gif)

## Usage

```go
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
```

## Installation

```console
$ go get github.com/unblee/cccp
```

## TODO

- Add tests

## Contribution

1.  Fork(https://github.com/unblee/cccp/fork)
2.  Create a branch (`git checkout -b my-fix`)
3.  Commit your changes (`git commit -am "fix something"`)
4.  Push to the branch (`git push origin my-fix`)
5.  Create a new [Pull Request](https://github.com/unblee/cccp/pulls)
6.  Have a coffee break and wait

## Author

[unblee](https://github.com/unblee)

## License

[![license](https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/unblee/cccp/blob/master/LICENSE)
