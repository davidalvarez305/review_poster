package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"strings"
	"time"
)

func FormatShortDate(dt time.Time) string {
	return dt.UTC().Format("20060102")
}

func FormatDate(dt time.Time) string {
	return dt.UTC().Format("20060102T150405Z")
}

func MakeHash(hash hash.Hash, b []byte) []byte {
	hash.Reset()
	hash.Write(b)
	return hash.Sum(nil)
}

func BuildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders, payloadHash string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders + "\n",
		signedHeaders,
		payloadHash,
	}, "\n")
}

func BuildStringToSign(date, credentialScope, canonicalRequest string) string {
	return strings.Join([]string{
		"AWS4-HMAC-SHA256",
		date,
		credentialScope,
		hex.EncodeToString(MakeHash(sha256.New(), []byte(canonicalRequest))),
	}, "\n")
}

func BuildCanonicalHeaders(contentType, contentEncoding, host, xAmazonDate, amazonTarget string) string {
	return strings.Join([]string{"content-encoding:" + contentEncoding, "content-type:" + contentType, "host:" + host, "x-amz-date:" + xAmazonDate, "x-amz-target:" + amazonTarget}, "\n")
}

func BuildSignature(strToSign string, sig []byte) (string, error) {
	return hex.EncodeToString(HMACSHA256(sig, []byte(strToSign))), nil
}

func HMACSHA256(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
