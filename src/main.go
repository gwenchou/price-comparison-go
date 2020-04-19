package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	parser "price-comparison/src/parser"
)

type Website struct {
	name, url, method string
	parser            parser.Strategy
}

type ParseResult struct {
	result []parser.Result
	error  error
}

func main() {
	http.HandleFunc("/query_product", query_product)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

func query_product(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}

	ch := make(chan ParseResult)

	keyword := r.Form["keyword"][0]
	var websites = [3]Website{
		Website{
			name:   "pchome",
			url:    "https://ecshweb.pchome.com.tw/search/v3.3/all/results?page=1&sort=sale/dc&q=" + url.QueryEscape(keyword),
			method: "GET",
			parser: parser.Strategy{Parser: parser.Pchome{}},
		},
		Website{
			name:   "friday",
			url:    "https://mservice-event.shopping.friday.tw/api/v2/search?&device=desktop&rows=20&page=1&keyword=" + url.QueryEscape(keyword),
			method: "POST",
			parser: parser.Strategy{Parser: parser.Friday{}},
		},
		Website{
			name:   "pcone",
			url:    "https://www.pcone.com.tw/api/filterSearchTP?from=pc&items_per_page=20&page=1&sortBy=default&sortDir=asc&q=" + url.QueryEscape(keyword),
			method: "GET",
			parser: parser.Strategy{Parser: parser.Pcone{}},
		},
	}

	for _, website := range websites {
		go fetch(w, website, ch)
	}

	results := map[string][]parser.Result{}
	for _, website := range websites {
		parseResult := <-ch
		if parseResult.error != nil {
			continue
		}

		results[website.name] = parseResult.result
	}

	jsonResults, _ := json.Marshal(results)
	fmt.Fprint(w, string(jsonResults))
}

func fetch(w http.ResponseWriter, website Website, ch chan<- ParseResult) {
	var parseResult ParseResult
	client := &http.Client{}

	req, err := http.NewRequest(website.method, website.url, nil)

	if err != nil {
		parseResult.error = err
		ch <- parseResult
		return
	}

	res, err := client.Do(req)

	if err != nil {
		parseResult.error = err
		ch <- parseResult
		return
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		parseResult.error = err
		ch <- parseResult
		return
	}

	parseResult.result = website.parser.Parse(string(body))

	ch <- parseResult

	return
}
