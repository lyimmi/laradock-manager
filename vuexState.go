package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var sKey = "2V6dUxoPqZFC7pK6779R07S605icY6CS"

//VuexStore store's content
type VuexStore struct {
	Settings map[string]string `json:"settings"`
}

//VuexState struct
type VuexState struct{}

//NewVuexState Create a new VuexState struct
func NewVuexState() *VuexState {
	result := &VuexState{}
	return result
}

//encode Encodes vuex store's data
func (t *VuexState) encode(data string) ([]byte, error) {
	plainText := []byte(data)

	block, err := aes.NewCipher([]byte(sKey))
	if err != nil {
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

	return cipherText, nil
}

//decode Decodes vuex store's data
func (t *VuexState) decode(data []byte) ([]byte, error) {
	cipherText := data

	block, err := aes.NewCipher([]byte(sKey))
	if err != nil {
		return nil, err
	}

	if len(cipherText) < aes.BlockSize {
		err = errors.New("ciphertext block size is too short")
		return nil, err
	}

	//IV needs to be unique, but doesn't have to be secure.
	//It's common to put it at the beginning of the ciphertext.
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)
	return cipherText, nil
}

//Write Writes vuex store data to file
func (t *VuexState) Write(data string) bool {

	encodedData, err := t.encode(data)
	if err != nil {
		return false
	}

	f, err := os.Create("vuex-store.data")
	if err != nil {
		fmt.Println(err)
		return false
	}
	l, err := f.Write(encodedData)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return false
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

//Read Reads vuex store data from file
func (t *VuexState) Read() string {
	fmt.Println(string("called"))
	data, err := ioutil.ReadFile("vuex-store.data")
	if err != nil {
		fmt.Println(err)
		return "{}"
	}

	decodedData, err := t.decode(data)
	if err != nil {
		fmt.Println(err)
		return "{}"
	}
	return string(decodedData)
}
