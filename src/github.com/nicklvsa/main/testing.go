package main

import (
	"fmt"

	syncgo "../syncgo"
)

//testing
func main() {
	s := new(syncgo.Sync)
	s.Init("https://example.com/some_backend.php", "file", nil)

	response, err := s.SyncDir("../syncgo/test")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(response))
}
