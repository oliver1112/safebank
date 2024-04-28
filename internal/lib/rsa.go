package lib

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// openssl genrsa -out rsa_private_key.pem 1024
var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXwIBAAKBgQDua+b+w7zlPaZ7MoFQ5GqdYfeluRV2jC2/CraY6p6QJPsySZV1
uDhjw7aUd+tZORWuQJ3EJkvPaaqaFxzCumi5OER9HlhnFzFoXSjTOBxJAxQ0fVxl
pRDQ/m9ei6DsESO1SFYHh8AroSnjpiHcziTaB3Sjr1QxwUO3IBDgN546KQIDAQAB
AoGBAIzQl+6yRreSUOiEbNINBlcLzItJpD6PDlm+BxiLwbjazq3lvet4MX3i0swf
g5X/3Ck9qrB1+eJ3wzYdHnR4Sm6uy4mXcYSIjS9NF5Tq7n2wv4i+6RL7h3re7pMa
XhMVkILA9Qra5JU3q61DUi5TM51hAJ1JROyaP8CsZVX5v4w9AkEA+pTjKC/PdHbz
cftPkZQghJqKX/ROrXpVTCnGl5ap/07v6cTfp9XdX8tRvU5d1fDZoYULeUM3L9r3
7W6+KMYc7wJBAPOTs67Gc0/wU7cxtRQeozzg0Qi5+cvcai4qQH4uX318inhyuGgB
G/t0qsVOokV8VS7Sv/sOgz9ZeeQRa7wsymcCQQCuT2XZLbD9TkW48139YfJg6/P4
HcWhTakKS0E3b/offLTNhEMkyFOvcIsSyfHigiGSBy/dEdHQ+1xeERw81tuHAkEA
pYvXbYwXR1dxrmqsRZZlH7U0nRe5POL7j5DL8HaYE/OXMTHXP2ixmf+7KQq+ozdT
tdUrAfjlHyMzAt0MOgK/NQJBAK+aa4jQ4Y83w6J1ilN3VeNgtXQDDQQ1ZOdbQjVP
9Zwcsebydz5Gh40hNjSDWTtQdRUsufHZu9KvjtP4CBQKdYU=
-----END RSA PRIVATE KEY-----
`)

// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDua+b+w7zlPaZ7MoFQ5GqdYfel
uRV2jC2/CraY6p6QJPsySZV1uDhjw7aUd+tZORWuQJ3EJkvPaaqaFxzCumi5OER9
HlhnFzFoXSjTOBxJAxQ0fVxlpRDQ/m9ei6DsESO1SFYHh8AroSnjpiHcziTaB3Sj
r1QxwUO3IBDgN546KQIDAQAB
-----END PUBLIC KEY-----
`)

func RsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func main() {
	data, _ := RsaEncrypt([]byte("hello world"))
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	origData, _ := RsaDecrypt(data)
	fmt.Println(string(origData))
}
