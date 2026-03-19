package openapi

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/url"
	"sort"
	"strings"
)

func SignString(stringToSign, secret string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	_, _ = mac.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func ComposeStringToSign(queries map[string]string) string {
	return percentEncode(GenSortedParamsStr(queries))
}

func GenSortedParamsStr(queries map[string]string) string {
	keys := make([]string, 0, len(queries))
	for key := range queries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, percentEncode(key)+"="+percentEncode(queries[key]))
	}
	return strings.Join(parts, "&")
}

func percentEncode(value string) string {
	encoded := url.QueryEscape(value)
	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "*", "%2A")
	encoded = strings.ReplaceAll(encoded, "%7E", "~")
	return encoded
}
