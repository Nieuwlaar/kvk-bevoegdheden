package bevoegdheden

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var client = &http.Client{Transport: tr}

func SearchCompanies(searchTerm string, apiKey string, useCache bool) ([]byte, error) {
	searchTermRegexp := regexp.MustCompile(`^[\w\-\s]+$`)
	if searchTerm == "" || !searchTermRegexp.MatchString(searchTerm) {
		fmt.Println("invalid searchterm")
		return nil, errors.New("invalid searchterm")
	}

	cachePath := "cache-search"
	if useCache {
		cachedBody, err := os.ReadFile(cachePath + "/" + searchTerm + ".json")
		if err == nil {
			fmt.Println("using cache")
			return cachedBody, nil
		}
	}

	fmt.Println("not using cache")

	termEscaped := url.PathEscape(searchTerm)
	fmt.Println("term: ", termEscaped)

	match, _ := regexp.MatchString("^[0-9]{8}$", termEscaped)

	category := "handelsnaam"
	if match {
		category = "kvkNummer"
	}

	url := fmt.Sprintf("https://api.kvk.nl/api/v1/zoeken?%s=%s", category, termEscaped)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// if useCache {
	// 	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
	// 		os.MkdirAll(cachePath, 0700)
	// 	}
	// 	_ = os.WriteFile(cachePath+"/"+searchTerm+".json", body, 0644)
	// }

	return body, nil
}
