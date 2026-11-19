package bcoding

import (
	"fmt"
	"io"
	"os"
	"testing"
)

const BUF_SZ = 1024

func TestBdecode(t *testing.T) {
	// read ubuntu-26.04.1-desktop-amd64.iso.torrent into stream
	// file, err := os.Open("ubuntu-26.04.1-live-server-amd64.iso.torrent")
	file, err := os.Open("small.torrent")
	if err != nil {
		fmt.Println("Error opening file")
	}
	defer file.Close()

	// decode stream into map
	buf := make([]byte, BUF_SZ)
	ch := make(chan []byte)
	go func() {
		for {
			_, err := file.Read(buf)
			if err != nil {
				if err != io.EOF {
					fmt.Println(err)
				}
				break
			}
			ch <- buf
		}
		close(ch)
	}()
	result, err := Bdecode(ch)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)

	// assert map values
}
