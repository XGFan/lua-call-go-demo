package main

import "C"

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

// go build -o sign.so -buildmode=c-shared sign.go
var privateKey *rsa.PrivateKey

func init() {
	fmt.Println("Loading Key")
	bytes, err := ioutil.ReadFile("private.pem")
	if err != nil {
		fmt.Println("Load Key Fail")
	}
	decodeBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(bytes)))
	if err != nil {
		fmt.Println("Load Key Fail")
	}
	key, err := x509.ParsePKCS1PrivateKey(decodeBytes)
	if err != nil {
		fmt.Println("Load Key Fail", err)
	}
	privateKey = key
}

//export Sign
func Sign(str string) *C.char {
	h := sha256.New()
	h.Write([]byte(str))
	d := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, d)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Signature in byte: %v\n\n", signature)
	encodedSig := base64.StdEncoding.EncodeToString(signature)
	//fmt.Printf("Encoded signature: %v\n\n", encodedSig)
	return C.CString(encodedSig)
}

//export Hello
func Hello() *C.char {
	return C.CString(time.Now().Format(time.RFC3339Nano))
}

func main() {
}
