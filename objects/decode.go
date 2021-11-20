package objects

import (
	"bufio"
	"fmt"
	"os"

	bencode "github.com/jackpal/bencode-go"
)

func DecodeBencodeFile(filename string) (interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file for bencode decoding: %v", err)
	}

	reader := bufio.NewReader(file)

	data, err := bencode.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("could not bencode decode file: %v", err)
	}

	return data, nil
}
