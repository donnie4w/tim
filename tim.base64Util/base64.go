/**
 * donnie4w@gmail.com  tim server
 */
package base64Util

import (
	"encoding/base64"
)

func Base64Encode(src []byte) string {
	//	return []byte(coder.EncodeToString(src))
	return base64.StdEncoding.EncodeToString(src)
}

func Base64Decode(src string) ([]byte, error) {
	//	return coder.DecodeString(string(src)	)
	return base64.StdEncoding.DecodeString(src)
}
