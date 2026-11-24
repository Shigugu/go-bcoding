package bcoding

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type nestedStruct int
type lexerState int

const (
	dictStruct nestedStruct = iota
	listStruct
	intStruct
)

const (
	defaultState lexerState = iota
	readingLen
	readingStr
	readingInt
)

const (
	dictByte byte = 'd'
	listByte byte = 'l'
	intByte  byte = 'i'
	endByte  byte = 'e'
)

type lexer struct {
	tokens        []string
	nestedStructs []nestedStruct
	currentToken  []byte
	state         lexerState
	length        int // next str token length
	position      int
}

func (l *lexer) intLexxer(chr byte) error {
	if chr == '-' || '0' <= chr && chr <= '9' {
		l.currentToken = append(l.currentToken, chr)
	} else if chr == endByte {
		l.tokens = append(l.tokens, string(l.currentToken))
		l.currentToken = nil
		l.tokens = append(l.tokens, string(chr))
		l.state = defaultState
	} else {
		return fmt.Errorf("expected numeric or sign character, found: %c at position %d", chr, l.position)
	}

	return nil
}

func (l *lexer) lenLexxer(chr byte) error {
	if '0' <= chr && chr <= '9' {
		l.currentToken = append(l.currentToken, chr)

		return nil
	} else if chr == ':' {
		var err error
		l.tokens = append(l.tokens, string(l.currentToken))
		l.state = readingStr
		l.length, err = strconv.Atoi((l.tokens)[len(l.tokens)-1])
		if err != nil {
			return fmt.Errorf("could not convert token to int: %s at position %d", l.currentToken, l.position)
		}
		l.currentToken = nil
		if l.length <= 0 {
			return fmt.Errorf("str length should be positive: %d at position %d", l.length, l.position)
		}

		return nil
	}
	return fmt.Errorf("expected numeric or sign character, found: %c at position %d", chr, l.position)
}

func (l *lexer) strLexxer(chr byte) error {
	l.currentToken = append(l.currentToken, chr)
	l.length -= 1
	if l.length == 0 {
		l.tokens = append(l.tokens, string(l.currentToken))
		l.currentToken = nil
		l.state = defaultState
	}

	return nil
}

func (l *lexer) defaultLexxer(chr byte) error {
	switch {
	case chr == dictByte:
		l.tokens = append(l.tokens, string(chr))
		l.nestedStructs = append(l.nestedStructs, dictStruct)
		l.state = readingLen
	case chr == listByte:
		l.tokens = append(l.tokens, string(chr))
		l.nestedStructs = append(l.nestedStructs, listStruct)
	case chr == endByte:
		l.tokens = append(l.tokens, string(chr))
		// pop last nested struct
		l.nestedStructs = (l.nestedStructs)[:len(l.nestedStructs)-1]
	case chr == intByte:
		l.tokens = append(l.tokens, string(chr))
		l.nestedStructs = append(l.nestedStructs, intStruct)
		l.state = readingInt
	case '0' <= chr && chr <= '9':
		l.currentToken = append(l.currentToken, chr)
		l.state = readingLen
	case chr == 0:
		return nil
	default:
		return fmt.Errorf("unexpected token while tokenizing: %c at position %d", chr, l.position)
	}

	return nil
}

func btokenize(chBuf chan []byte, chAck chan bool) ([]string, error) {
	var lexer = &lexer{}
	var err error

	for str := range chBuf {
		for _, chr := range str {
			switch lexer.state {
			case defaultState:
				err = lexer.defaultLexxer(chr)
			case readingLen:
				err = lexer.lenLexxer(chr)
			case readingStr:
				err = lexer.strLexxer(chr)
			case readingInt:
				err = lexer.intLexxer(chr)
			}
			if err != nil {

				return nil, err
			}
			lexer.position += 1
		}
		chAck <- true
	}

	return lexer.tokens, nil
}

func Bdecode(chBuf chan []byte, chAck chan bool) (map[string]any, error) {
	m := make(map[string]any)
	_, err := btokenize(chBuf, chAck)
	if err != nil {
		return nil, err
	}
	// fmt.Println(tokens)

	// parse tokens

	return m, nil
}

func BdecodeFile(filename string) (map[string]any, error) {
	const BUF_SZ = 128

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	buf := make([]byte, BUF_SZ)
	chBuf := make(chan []byte)
	var chAck = make(chan bool, 1)
	go func() {
		for {
			n, err := file.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			chBuf <- buf[:n]
			// wait for ack
		outer:
			for {
				select {
				case <-chAck:
					break outer
				default:
					time.Sleep(100 * time.Nanosecond)
				}
			}
		}
		close(chBuf)
	}()

	return Bdecode(chBuf, chAck)
}
