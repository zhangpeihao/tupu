package tupu

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Request struct {
	url        string
	secretId   string
	modelId    string
	privateKey *rsa.PrivateKey
}

func NewRequest(url, secretId, modelId string, privateKey *rsa.PrivateKey) (req *Request) {
	req = &Request{
		url:        url,
		secretId:   secretId,
		modelId:    modelId,
		privateKey: privateKey,
	}
	return
}

func (req *Request) CheckSingleImage(imgBuf *bytes.Buffer, imgName string) (resp *Response, err error) {
	imgHashStr := fmt.Sprintf("%x", sha1.Sum(imgBuf.Bytes()))
	timestamp := time.Now().Format(time.RFC1123)

	nonce_int, err := rand.Int(rand.Reader, big.NewInt(9999999999))
	if err != nil {
		fmt.Printf("rand.Int, error: %s\n", err.Error())
		return nil, err
	}
	nonce := nonce_int.String()
	sign_params := []string{
		imgHashStr,
		req.secretId,
		req.modelId,
		timestamp,
		nonce,
	}
	sort.Strings(sign_params)
	signStr := strings.Join(sign_params, ",")

	sign_hash_bytes := sha256.Sum256([]byte(signStr))
	sign_hash := sign_hash_bytes[:]

	sign, err := rsa.SignPKCS1v15(rand.Reader, req.privateKey, crypto.SHA256, sign_hash)
	if err != nil {
		fmt.Println("rsa.SignPKCS1v15 err:", err)
		return nil, err
	}

	sign_base64 := base64.StdEncoding.EncodeToString(sign)

	tupu_req_buffer := new(bytes.Buffer)
	multipart_writer := multipart.NewWriter(tupu_req_buffer)
	multipart_writer.WriteField("modelId", req.modelId)
	multipart_writer.WriteField("secretId", req.secretId)
	multipart_writer.WriteField("timestamp", timestamp)
	multipart_writer.WriteField("nonce", nonce)
	multipart_writer.WriteField("signature", sign_base64)
	image_multipart_writer, err := multipart_writer.CreateFormFile("image", imgName)
	if err != nil {
		fmt.Println("multipart_writer.CreateFormFile err:", err)
		return nil, err
	}
	io.Copy(image_multipart_writer, imgBuf)
	multipart_writer.Close()
	tupu_req, err := http.NewRequest("POST", req.url, tupu_req_buffer)
	if err != nil {
		fmt.Println("http.NewRequest err:", err)
		return nil, err
	}
	tupu_req.Header.Set("Content-Type", multipart_writer.FormDataContentType())
	var client http.Client
	tupu_resp, err := client.Do(tupu_req)
	if err != nil {
		fmt.Println("client.Do err:", err)
		return nil, err
	}

	defer tupu_resp.Body.Close()
	if tupu_resp.StatusCode != http.StatusOK {
		fmt.Println("tupu_resp.StatusCode:", tupu_resp.Status)
		return nil, err
	}
	tupu_resp_bytes, err := ioutil.ReadAll(tupu_resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return nil, err
	}
	var tupuResponseCapsule ResponseCapsule
	err = json.Unmarshal(tupu_resp_bytes, &tupuResponseCapsule)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return nil, err
	}

	var tupuResponse Response
	resp = &tupuResponse
	err = json.Unmarshal([]byte(tupuResponseCapsule.Json), resp)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return nil, err
	}
	return
}
