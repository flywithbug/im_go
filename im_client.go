package main

//import (
//	"flag"
//	"im_go/imc"
//
//)
//
//func main() {
//	flag.Parse()
//	imc.StartClient(62310)
//}



//func main() {
//	b, _ := ioutil.ReadFile("./public_test_key.pem")
//
//	data, err := rsa.RsaEncrypt([]byte("polaris@studygolang.com啊啊啊"),b) //RSA加密
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println("RSA加密", string(data))
//	b, _ = ioutil.ReadFile("./private_test_key.pem")
//
//	origData, err := rsa.RsaDecrypt(data,b) //RSA解密
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("RSA解密", string(origData))
//}