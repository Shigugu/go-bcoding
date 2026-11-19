package bcoding

import (
	"fmt"
	"strconv"
)

const DICT_STRUCT int = 1
const LIST_STRUCT int = 2
const INT_STRUCT int = 3

const DEFAULT_STATE int = 0
const READING_LEN int = 1
const READING_STR int = 2
const READING_INT int = 3

const DICT_BYTE byte = 'd'
const LIST_BYTE byte = 'l'
const INT_BYTE byte = 'i'
const END_BYTE byte = 'e'

func intLexxer(tokens *[]string, currentToken *[]byte, lexxerState *int, chr byte) error {
	if chr == '-' || '0' <= chr && chr <= '9' {
		*currentToken = append(*currentToken, chr)
	} else if chr == END_BYTE {
		*tokens = append(*tokens, string(*currentToken))
		*currentToken = nil
		*lexxerState = DEFAULT_STATE
	} else {
		return fmt.Errorf("expected numeric or sign character, found: %c", chr)
	}

	return nil
}

func lenLexxer(tokens *[]string, currentToken *[]byte, lexxerState *int, length *int, chr byte) error {
	if '0' <= chr && chr <= '9' {
		*currentToken = append(*currentToken, chr)
		return nil
	} else if chr == ':' {
		var err error
		*tokens = append(*tokens, string(*currentToken))
		*lexxerState = READING_STR
		*length, err = strconv.Atoi((*tokens)[len(*tokens)-1])
		if err != nil {
			return fmt.Errorf("could not convert token to int: %s", *currentToken)
		}
		*currentToken = nil
		if *length <= 0 {
			return fmt.Errorf("str length should be positive: %d", *length)
		}
		return nil
	}
	return fmt.Errorf("expected numeric or sign character, found: %c", chr)
}

func strLexxer(tokens *[]string, currentToken *[]byte, lexxerState *int, length *int, chr byte) error {
	*currentToken = append(*currentToken, chr)
	*length -= 1
	if *length == 0 {
		*tokens = append(*tokens, string(*currentToken))
		*currentToken = nil
		*lexxerState = DEFAULT_STATE
	}

	return nil
}

func defaultLexxer(tokens *[]string, currentToken *[]byte, nestedStructs *[]int, lexxerState *int, chr byte) error {
	switch {
	case chr == DICT_BYTE:
		*tokens = append(*tokens, string(chr))
		*nestedStructs = append(*nestedStructs, DICT_STRUCT)
		*lexxerState = READING_LEN
	case chr == LIST_BYTE:
		*tokens = append(*tokens, string(chr))
		*nestedStructs = append(*nestedStructs, LIST_STRUCT)
	case chr == END_BYTE:
		*tokens = append(*tokens, string(*currentToken))
		*currentToken = nil
		*tokens = append(*tokens, string(chr))
		// pop last nested struct
		*nestedStructs = (*nestedStructs)[:len(*nestedStructs)-1]
		if len(*nestedStructs) > 0 && (*nestedStructs)[len(*nestedStructs)-1] == DICT_STRUCT {
			*lexxerState = READING_LEN
		} else {
			*lexxerState = DEFAULT_STATE
		}
	case chr == INT_BYTE:
		*tokens = append(*tokens, string(chr))
		*nestedStructs = append(*nestedStructs, INT_STRUCT)
		*lexxerState = READING_INT
	case '0' <= chr && chr <= '9':
		*currentToken = append(*currentToken, chr)
		*lexxerState = READING_LEN
	case chr == 0:
		return nil
	default:
		return fmt.Errorf("unexpected token while tokenizing: %c", chr)
	}

	return nil
}

func btokenize(ch chan []byte) ([]string, error) {
	var tokens []string
	var currentToken []byte
	var nestedStructs []int
	var lexxerState int = DEFAULT_STATE
	var length int
	var err error

	for str := range ch {
		for _, chr := range str {
			switch lexxerState {
			case DEFAULT_STATE:
				err = defaultLexxer(&tokens, &currentToken, &nestedStructs, &lexxerState, chr)
			case READING_LEN:
				err = lenLexxer(&tokens, &currentToken, &lexxerState, &length, chr)
			case READING_STR:
				err = strLexxer(&tokens, &currentToken, &lexxerState, &length, chr)
			case READING_INT:
				err = intLexxer(&tokens, &currentToken, &lexxerState, chr)
			}
			if err != nil {
				// DEBUG
				fmt.Println(tokens)
				return nil, err
			}
		}
	}

	return tokens, nil
}

func Bdecode(ch chan []byte) (map[string]any, error) {
	m := make(map[string]any)
	tokens, err := btokenize(ch)
	if err != nil {
		return nil, err
	}
	fmt.Println(tokens)

	// parse tokens

	return m, nil
}
