
package utilCipher

import (
	"io"
	"encoding/hex"
	"encoding/pem"
	"encoding/base64"

	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
)

func RsaEncode(data string, key string) string {
	if(len(key) == 0) {
		return "";
	}
	bData := []byte(data);
	// bEncodeData, _ := base64.StdEncoding.DecodeString(data);
	// if bEncodeData == nil {
	// 	return "";
	// }
	block, _ := pem.Decode([]byte(key));
	tmp, _ := x509.ParsePKIXPublicKey([]byte(block.Bytes));
	if(tmp == nil) {
		return "";
	}
	rsaPrivateKey := tmp.(*rsa.PublicKey);
	brst, _ := rsa.EncryptPKCS1v15(rand.Reader, rsaPrivateKey, bData);
	if(brst == nil) {
		return "";
	}
	return base64.StdEncoding.EncodeToString(brst);
}

func RsaDecode(data string, key string) string {
	if(len(key) == 0) {
		return "";
	}
	bEncodeData, _ := base64.StdEncoding.DecodeString(data);
	if bEncodeData == nil {
		return "";
	}
	block, _ := pem.Decode([]byte(key));
	rsaPrivateKey, _ := x509.ParsePKCS1PrivateKey([]byte(block.Bytes));
	if(rsaPrivateKey == nil) {
		return "";
	}
	// rsaPrivateKey := tmp.(*rsa.PrivateKey);
	bdata, _ := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, bEncodeData);
	if(bdata == nil) {
		return "";
	}
	return string(bdata);
}

func GetRsaKey(rsaBase64 string) string {
	rst, _ := base64.StdEncoding.DecodeString(rsaBase64);
	if(rst == nil) {
		return "";
	}
	return string(rst);
}

func aesPaddingPKCS7(arrText []byte) []byte {
	blockSize := 16;
	count := blockSize;
	mod := len(arrText)%blockSize;
	if(mod != 0) {
		count = blockSize - mod;
	}
	rst := arrText;
	for i:=0; i < count; i++ {
		rst = append(rst, byte(count));
	}
	return rst;
}

func aesUnPaddingPKCS7(arrText []byte) []byte {
	if(len(arrText) == 0) {
		return arrText;
	}
	count := int(arrText[len(arrText) - 1]);
	if(len(arrText) < count) {
		return arrText;
	}
	rst := arrText[0: len(arrText) - count];
	// rst = rst;
	return rst;
}

func AesEncode(text string, key string) string {
	if(len(key) == 0) {
		return "";
	}
	arrKey, err1 := hex.DecodeString(key);
	if(err1 != nil) {
		return "";
	}
	keyObj, err2 := aes.NewCipher(arrKey);
	if(err2 != nil) {
		return "";
	}
	textByte := []byte(text);
	textByte = aesPaddingPKCS7(textByte);
	// ciphertext := make([]byte, aes.BlockSize + len(plaintext));
	textByteEncode := make([]byte, len(textByte));
	iv := make([]byte, aes.BlockSize);
	// iv := ciphertext[:aes.BlockSize];
	if _, err2 := io.ReadFull(rand.Reader, iv); err2 != nil {
		return "";
	}
	strIV := hex.EncodeToString(iv);
	// fmt.Println(arrKey);
	// fmt.Println(iv);

	stream := cipher.NewCFBEncrypter(keyObj, iv);
	stream.XORKeyStream(textByteEncode, textByte);

	textEncode := hex.EncodeToString(textByteEncode);
	// textEncode := base64.URLEncoding.EncodeToString(textByteEncode);

	return strIV + textEncode;

	// return base64.URLEncoding.EncodeToString(ciphertext);
	// return hex.EncodeToString(textByteEncode);
}

func AesDecode(text string, key string) string {
	if(len(key) == 0) {
		return "";
	}
	arrKey, err1 := hex.DecodeString(key);
	if(err1 != nil) {
		return "";
	}
	keyObj, err2 := aes.NewCipher(arrKey);
	if(err2 != nil) {
		return "";
	}
	
	if len(text) < len(key) {
		return "";
	}

	strIV := text[0:len(key)];
	iv, err3 := hex.DecodeString(strIV);
	if(err3 != nil) {
		return "";
	}

	text = text[len(key):];

	// textByte, err4 := base64.URLEncoding.DecodeString(text);
	textByte, err4 := hex.DecodeString(text);
	
	if err4 != nil {
		return "";
	}
	// iv := textByte[:aes.BlockSize];
	// textByte = textByte[aes.BlockSize:];

	// fmt.Println("iv ", iv);
	// fmt.Println("str", textByte);

	stream := cipher.NewCFBDecrypter(keyObj, iv);
	stream.XORKeyStream(textByte, textByte);
	
	textByte = aesUnPaddingPKCS7(textByte);

	return string(textByte);
}

func Hash(text string) string {
	hash := sha256.New();
	hash.Write([]byte(text));
	bytes := hash.Sum(nil);
	hashCode := hex.EncodeToString(bytes);
	return hashCode;
}
