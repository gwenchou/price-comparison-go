package parser

import (
	"encoding/json"
	"strconv"
)

type Pcone struct{}

type Product struct {
	Best_Price string
	Name       string
}

type PconeResult struct {
	Products []Product
}

func (Pcone) Apply(body string) []Result {
	var PconeResult PconeResult
	json.Unmarshal([]byte(body), &PconeResult)

	var result []Result

	for _, product := range PconeResult.Products {
		price, err := strconv.Atoi(product.Best_Price)
		if err != nil {
			continue
		}

		result = append(result, Result{Name: product.Name, Price: price})
	}

	return result
}
