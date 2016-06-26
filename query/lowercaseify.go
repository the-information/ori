package query

import (
	"bytes"
	"net/url"
)

// Lowercaseify efficiently lowercases the first character of every key in v.
// This is useful if you want to maintain consistency between JSON and Go
// field naming conventions in your queries.
//
// If buf is nil, Lowercaseify will create its own buffer.
func Lowercaseify(v url.Values, buf *bytes.Buffer) url.Values {

	if buf == nil {
		buf = bytes.NewBuffer(nil)
	}

	result := make(map[string][]string, len(v))
	for k, p := range v {

		buf.Reset()

		if k[0] >= 'A' && k[0] <= 'Z' {
			buf.WriteByte(k[0] + 32)
			buf.WriteString(k[1:])
		} else {
			buf.WriteString(k)
		}

		result[buf.String()] = p

	}
	return result

}
