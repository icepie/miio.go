package proto

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	_ "crypto/md5" // hash function required by crypto module
)

func md5(chunks ...[]byte) []byte {
	src := crypto.MD5.New()
	for _, chunk := range chunks {
		src.Write(chunk)
	}

	return src.Sum(nil)
}

type deviceKeys struct {
	key []byte
	iv  []byte
}

func newDeviceKeys(token []byte) deviceKeys {
	key := md5(token)

	return deviceKeys{
		key,
		md5(key, token),
	}
}

func (keys *deviceKeys) newCipher() cipher.Block {
	block, err := aes.NewCipher(keys.key)
	if err != nil {
		panic(err)
	}

	return block
}

func (keys *deviceKeys) encrypt(src []byte) []byte {
	mode := cipher.NewCBCEncrypter(keys.newCipher(), keys.iv)
	padded := padding(src, mode.BlockSize())
	dst := make([]byte, len(padded))
	mode.CryptBlocks(dst, padded)

	return dst
}

func (keys *deviceKeys) decrypt(src []byte) []byte {
	mode := cipher.NewCBCDecrypter(keys.newCipher(), keys.iv)
	dst := make([]byte, len(src))
	mode.CryptBlocks(dst, src)

	return dst
}

func padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, padtext...)
}
