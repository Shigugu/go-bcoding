package bcoding

import (
	"fmt"
	"io"
	"os"
	"time"
)

func bdecode(chBuf chan []byte, chAck chan bool) (map[string]any, error) {
	tokens, err := btokenize(chBuf, chAck)
	if err != nil {
		return nil, err
	}

	return bparse(tokens)
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

	return bdecode(chBuf, chAck)
}
