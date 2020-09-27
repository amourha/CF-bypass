package providers

import (
	b64 "encoding/base64"

	mmh3 "github.com/reusee/mmh3"
)

func b64Split(s string, size int) string {
	ss := ""
	for len(s) > 0 {
		if len(s) < size {
			size = len(s)
		}
		ss += s[:size] + "\n"
		s = s[size:]

	}
	return ss
}

func calcFaviconHash(favicon []byte) int32 {
	var b64Favicon = b64Split(b64.StdEncoding.EncodeToString(favicon), 76)
	return int32(mmh3.Hash32([]byte(b64Favicon)))
}
