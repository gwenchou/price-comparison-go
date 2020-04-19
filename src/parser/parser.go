package parser

type Parser interface {
	Apply(string) []Result
}

type Strategy struct {
	Parser Parser
}

type Result struct {
	Name  string
	Price int
}

func (s *Strategy) Parse(body string) []Result {
	return s.Parser.Apply(body)
}
