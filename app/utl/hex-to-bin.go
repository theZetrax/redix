package utl

import "encoding/hex"

func DecodeHexToBinary(raw_hex []byte) ([]byte, error) {
	src := raw_hex
	dest := make([]byte, hex.DecodedLen(len(src)))
	n, err := hex.Decode(dest, src)
	if err != nil {
		return []byte{}, err
	}

	// Convert hex to binary
	return dest[:n], nil
}
