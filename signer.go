package signature

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
)

//签名接口
type signer interface {
	Sign(message string) (string, error)
}

func NewSigner(accessKeySecret string) *SHA1Signer {
	return &SHA1Signer{accessKeySecret: accessKeySecret}
}

//SHA1签名
type SHA1Signer struct {
	accessKeySecret string
}

func (s *SHA1Signer) Sign(message string) (string, error) {
	hashfun := hmac.New(sha1.New, []byte(s.accessKeySecret))
	hashfun.Write([]byte(message))
	rawSignature := hashfun.Sum(nil)
	base64signature := base64.StdEncoding.EncodeToString(rawSignature)
	return base64signature, nil
}
