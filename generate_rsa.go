package main
//
//import (
//	"crypto/rand"
//	"crypto/rsa"
//	"crypto/x509"
//	"encoding/pem"
//	"os"
//)
//
//func main() {
//	//rsa 密钥文件产生
//	GenRsaKey(1024)
//}
//
////RSA公钥私钥产生
//func GenRsaKey(bits int) error {
//	// 生成私钥文件
//	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
//	if err != nil {
//		return err
//	}
//	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
//	block := &pem.Block{
//		Type:  "RSA PRIVATE KEY",
//		Bytes: derStream,
//	}
//	file, err := os.Create("./private.pem")
//	if err != nil {
//		return err
//	}
//	err = pem.Encode(file, block)
//	if err != nil {
//		return err
//	}
//	// 生成公钥文件
//	publicKey := &privateKey.PublicKey
//	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
//	if err != nil {
//		return err
//	}
//	block = &pem.Block{
//		Type:  "PUBLIC KEY",
//		Bytes: derPkix,
//	}
//	file, err = os.Create("./public.pem")
//	if err != nil {
//		return err
//	}
//	err = pem.Encode(file, block)
//	if err != nil {
//		return err
//	}
//	return nil
//}