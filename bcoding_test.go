package bcoding

import (
	"fmt"
	"testing"
)

func TestBdecode(t *testing.T) {
	// var filename = "small.torrent"
	var filename = "ubuntu-26.04.1-live-server-amd64.iso.torrent"

	result, err := BdecodeFile(filename)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)

	// assert map values
}
