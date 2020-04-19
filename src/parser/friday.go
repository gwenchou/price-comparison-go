package parser

import "encoding/json"

type Friday struct{}

type Data struct {
	Name       string
	Real_Price int
}

type Payload struct {
	Data []Data
}

type FridayResult struct {
	Payload Payload
}

func (Friday) Apply(body string) []Result {
	var fridayResult FridayResult
	json.Unmarshal([]byte(body), &fridayResult)

	var result []Result

	for _, prod := range fridayResult.Payload.Data {
		result = append(result, Result{Name: prod.Name, Price: prod.Real_Price})
	}

	return result
}
