package auth

import (
	"io/ioutil"
)

var privateKeyPath string
var publicKeyPath string

func InitKeys(privateKey, publicKey string) {
	privateKeyPath = privateKey
	publicKeyPath = publicKey
}

func GetPrivateKey() []byte {
	if len(privateKeyPath) <= 0 {
		return []byte("secret")
	}
	privateBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}

	return privateBytes
}

func GetPublicKey() []byte {
	publicBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		panic(err)
	}

	return publicBytes
}
