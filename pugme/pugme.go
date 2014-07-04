package pugme

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Pugs is a bunch of jucy dogs
type Pugs struct {
	Pugs []string
}

// RandomPugs return a bunch of pugs
func RandomPugs(count int) []string {
	resp, err := http.Get(fmt.Sprintf("http://pugme.herokuapp.com/bomb?count=%d", count))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	r := new(Pugs)
	json.Unmarshal(body, r)
	return r.Pugs
}

// Get filename from URL
func getFileName(URL string) (string, error) {
	r, err := url.Parse(URL)
	if err != nil {
		return "", errors.New("failed to parse url")
	}

	urlPath := strings.Split(r.Path, "/")
	return urlPath[len(urlPath)-1], nil
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// DownloadPugs to filepath
func DownloadPugs(count int, path string) {

	if !fileExist(path) {
		fmt.Println(path + " does not exist...")
		os.Exit(1)
	}
	result := make(chan string)
	pugs := RandomPugs(count)
	var wg sync.WaitGroup
	wg.Add(len(pugs))

	for _, url := range pugs {
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			filename, err := getFileName(url)
			if err != nil {
				return
			}

			filePath := filepath.Join(path, filename)

			// Skip pug if it already exist
			if fileExist(filePath) {
				return
			}
			// Store pug
			ioutil.WriteFile(filePath, body, 0654)
			result <- filePath
		}(url)
	}

	go func() {
		for filePath := range result {
			fmt.Println(filePath)
		}
	}()

	wg.Wait()
}
