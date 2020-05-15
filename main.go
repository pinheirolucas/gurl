package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var out string

func main() {
	flag.StringVar(
		&out,
		"o",
		"./gurl.gif",
		"the path where the gif will be stored (defaults to ./out.gif)",
	)
	flag.Parse()

	out = strings.TrimSpace(out)
	if out == "" {
		fmt.Println("empty destination path")
		os.Exit(-1)
	}

	args := flag.Args()
	if len(args) == 0 {
		printHelp()
	}

	gifID, err := getGifID(args[0])
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	downloadURL := fmt.Sprintf("https://i.giphy.com/media/%s/giphy.gif", gifID)
	c := http.Client{
		Timeout: 20 * time.Second,
	}
	resp, err := c.Get(downloadURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	file, err := os.Create(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer file.Close()

	absOut, err := filepath.Abs(out)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("writing gif from %s to %s \n", downloadURL, absOut)
	if err := resp.Write(file); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	fmt.Printf("gif downloaded to %s\n", absOut)
}

func getGifID(u string) (string, error) {
	s := strings.Split(u, "/")

	if isGifPage(u) {
		if len(s) < 5 {
			return "", errors.New("the provided URL does not have an id")
		}

		code := s[4]
		cs := strings.Split(code, "-")
		return cs[len(cs)-1], nil
	}

	if isSharePage(u) {
		if len(s) < 6 {
			return "", errors.New("the provided URL does not have an id")
		}

		return s[4], nil
	}

	return "", errors.New("the provided URL is invalid")
}

func isGifPage(u string) bool {
	return strings.HasPrefix(u, "https://giphy.com/gifs/")
}

func isSharePage(u string) bool {
	return strings.HasPrefix(u, "https://media.giphy.com/media/") &&
		strings.HasSuffix(u, "/giphy.gif")
}

func printHelp() {
	fmt.Println("You should provide a giphy URL. Example command: gurl https://giphy.com/gifs/dog-world-doggy-dBRaPog8yxFWU")
	os.Exit(-1)
}
