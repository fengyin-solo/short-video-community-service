package tagparse

type Parser struct{ buffer []string }

func NewParser() *Parser { return &Parser{buffer: make([]string, 2)} }

func (p *Parser) Parse(values []string) []string {
	copy(p.buffer, values)
	return p.buffer[:len(values)]
}
