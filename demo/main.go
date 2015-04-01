package main

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/zhangpeihao/tupu"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	SECRET_ID    = "<-- Your SECRET_ID -->"
	MODEL_ID     = "<-- Your MODEL_ID -->"
	TUPU_API_URL = `<-- Your TUPU_API_URL -->`
	PRIVATE_KEY  = `-----BEGIN PRIVATE KEY-----
<-- Your PRIVATE_KEY -->
-----END PRIVATE KEY-----
`
)

var (
	test_image_url = "<-- Your image url -->"
)

func main() {
	private_key, err := load_key()
	if err != nil {
		fmt.Printf("load_key err: %s\n", err.Error())
		os.Exit(-1)
	}
	img_buf, err := get_image(test_image_url)
	if err != nil {
		fmt.Printf("get_image err: %s\n", err.Error())
		os.Exit(-1)
	}
	req := tupu.NewRequest(TUPU_API_URL, SECRET_ID, MODEL_ID, private_key)
	splits := strings.Split(test_image_url, "/")
	resp, err := req.CheckSingleImage(img_buf, splits[len(splits)-1])
	if err != nil {
		fmt.Printf("CheckSingleImage err: %s\n", err.Error())
		os.Exit(-1)
	}
	fmt.Printf("response: %+v\n", resp)
}

func load_key() (privateKey *rsa.PrivateKey, err error) {
	parent_private_pem_block, _ := pem.Decode([]byte(PRIVATE_KEY))

	x509_private_key, err := x509.ParsePKCS8PrivateKey(parent_private_pem_block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err:", err)
		return
	}
	var ok bool
	privateKey, ok = x509_private_key.(*rsa.PrivateKey)
	if !ok {
		fmt.Println("Not RSA key")
		return nil, errors.New("Load key error")
	}
	return
}

func get_image(image_url string) (img_buf *bytes.Buffer, err error) {
	var resp *http.Response

	resp, err = http.Get(image_url)
	if err != nil {
		fmt.Printf("http.Get(%s), error: %s\n", image_url, err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Get %s, response %s", image_url, resp.Status))
	}
	img_buf = new(bytes.Buffer)
	_, err = io.Copy(img_buf, resp.Body)
	if err != nil {
		fmt.Printf("io.Copy, error: %s\n", err.Error())
		return
	}
	return
}
