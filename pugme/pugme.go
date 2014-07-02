package pugme

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// Pug is just one dog
type Pug struct {
	Pug string
}

// Pugs is a bunch of jucy dogs
type Pugs struct {
	Pugs []string
}

// RandomPug returns a random pug URL
func RandomPug() string {
	resp, err := http.Get("http://pugme.herokuapp.com/random")
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
	r := new(Pug)
	json.Unmarshal(body, r)
	return r.Pug
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
func getFileName(URL string) string {
	r, err := url.Parse(URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	urlPath := strings.Split(r.Path, "/")
	return urlPath[len(urlPath)-1]
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

	pugs := RandomPugs(count)

	for i := 0; i < count; i++ {
		url := pugs[i]
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		filePath := filepath.Join(path, getFileName(url))

		// Skip pug if it already exist
		if fileExist(filePath) {
			continue
		}
		// Store pug
		ioutil.WriteFile(filePath, body, 0654)
	}
}
