package rsa

import (
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"im_go/config"
	"errors"
	"io/ioutil"
	"fmt"
)

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	b := config.Conf().RSAConfig.Public
	if b == nil {
		b, _ = ioutil.ReadFile("./public.pem")
	}
	fmt.Println(string(b))
	block, _ := pem.Decode(b) //将密钥解析成公钥实例
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData) //RSA算法加密
}


// 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	b := config.Conf().RSAConfig.Public
	if b == nil {
		b, _ = ioutil.ReadFile("./private.pem")
	}
	fmt.Println(string(b))
	block, _ := pem.Decode(b) //将密钥解析成私钥实例
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext) //RSA算法解密
}

