package parser

import "encoding/json"

type Pchome struct{}

type Prod struct {
	Name  string
	Price int
}

type PchomeResult struct {
	Prods []Prod
}

func (Pchome) Apply(body string) []Result {
	var pchomeResult PchomeResult
	json.Unmarshal([]byte(body), &pchomeResult)

	var result []Result

	for _, prod := range pchomeResult.Prods {
		result = append(result, Result{Name: prod.Name, Price: prod.Price})
	}

	return result
}
