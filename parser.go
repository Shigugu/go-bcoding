package bcoding

import (
	"fmt"
	"strconv"
)

// FIXME: use other values that can't possibly be used in the .torrent file
const (
	dictToken = string(dictByte)
	listToken = string(listByte)
	intToken  = string(intByte)
	endToken  = string(endByte)
)

type bParser struct {
	tokens []string
	index  int
}

// returns currentToken and move the cursor forward
func (p *bParser) next() (string, error) {
	if p.index < len(p.tokens) {
		token := p.tokens[p.index]
		p.index += 1

		return token, nil
	}

	return "", fmt.Errorf("Unexpectedly reached last token")
}

func (p *bParser) parseInt() (int, error) {
	var token string
	var err error
	var value int

	token, err = p.next()
	if err != nil {
		return -1, err
	}
	value, err = strconv.Atoi(token)
	if err != nil {
		return -1, err
	}

	return value, nil
}

func (p *bParser) parseList() ([]any, error) {
	var result []any
	var err error
	var token string
	var value any

	for {
		token, err = p.next()
		if err != nil {
			return nil, err
		}

		switch token {
		case dictToken:
			value, err = p.parseDict()
		case listToken:
			value, err = p.parseList()
		case intToken:
			value, err = p.parseInt()
		case endToken:
			return result, nil
		default:
			value = token
		}
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
}

func (p *bParser) parseDict() (map[string]any, error) {
	var result = make(map[string]interface{})
	var key, value string
	var err error

	for {
		key, err = p.next()
		if err != nil {
			return nil, err
		}

		if key == endToken {
			return result, nil
		}

		value, err = p.next()
		if err != nil {
			return nil, err
		}

		switch value {
		case dictToken:
			result[key], err = p.parseDict()
		case listToken:
			result[key], err = p.parseList()
		case intToken:
			result[key], err = p.parseInt()
		case endToken:
			return result, nil
		default:
			result[key] = value
		}
		if err != nil {
			return nil, err
		}
	}
}

func bparse(tokens []string) (map[string]any, error) {
	var parser = &bParser{tokens: tokens}
	var err error
	var result map[string]any

	token, err := parser.next()
	if err != nil {
		return nil, err
	}
	if token != dictToken {
		return nil, fmt.Errorf("Unexpected root token: %s", token)
	}
	result, err = parser.parseDict()
	if err != nil {
		return nil, err
	}
	if parser.index < len(parser.tokens) {
		return nil, fmt.Errorf("Unexpected tokens after root node end: %v", parser.tokens[parser.index:])
	}

	return result, nil
}
