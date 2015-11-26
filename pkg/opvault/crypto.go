package opvault

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"hash"
)

var (
	opdata01 = []byte("opdata01")
)

func decrypt(dst, src []byte, encKey, macKey []byte) ([]byte, error) {
	var (
		macBuf  [32]byte
		macSrc  []byte
		mac     hash.Hash
		dataLen uint64
		dataIV  [16]byte
	)

	{ // verify mac
		if len(src) < 32 {
			return nil, errors.New("invalid opdata signature length")
		}

		macSrc = src[len(src)-32:]
		src = src[:len(src)-32]
		mac = hmac.New(sha256.New, macKey)
		mac.Write(src)
		mac.Sum(macBuf[:0])
		if subtle.ConstantTimeCompare(macBuf[:], macSrc) != 1 {
			return nil, errors.New("invalid opdata signature")
		}
	}

	{ // get header
		if len(src) < 32 {
			return nil, errors.New("invalid opdata header length")
		}
		if !bytes.HasPrefix(src, opdata01) {
			return nil, errors.New("invalid opdata header")
		}

		dataLen = binary.LittleEndian.Uint64(src[8:])
		copy(dataIV[:], src[16:])
		src = src[32:]
	}

	{
		block, err := aes.NewCipher(encKey)
		if err != nil {
			return nil, err
		}

		if dst == nil || cap(dst) < len(src) {
			dst = make([]byte, len(src))
		}
		if len(dst) != len(src) {
			dst = dst[:len(src)]
		}

		mode := cipher.NewCBCDecrypter(block, dataIV[:])
		mode.CryptBlocks(dst, src)

		dst = dst[len(dst)-int(dataLen):]
	}

	return dst, nil
}

func decryptKey(dst, src []byte, encKey, macKey []byte) ([]byte, error) {
	var (
		macBuf [32]byte
		macSrc []byte
		mac    hash.Hash
		dataIV [16]byte
	)

	{ // verify mac
		if len(src) < 32 {
			return nil, errors.New("invalid opdata signature length")
		}

		macSrc = src[len(src)-32:]
		src = src[:len(src)-32]
		mac = hmac.New(sha256.New, macKey)
		mac.Write(src)
		mac.Sum(macBuf[:0])
		if subtle.ConstantTimeCompare(macBuf[:], macSrc) != 1 {
			return nil, errors.New("invalid opdata signature")
		}
	}

	{ // get header
		if len(src) < 16 {
			return nil, errors.New("invalid opdata header length")
		}

		copy(dataIV[:], src)
		src = src[16:]
	}

	{
		block, err := aes.NewCipher(encKey)
		if err != nil {
			return nil, err
		}

		if dst == nil || cap(dst) < len(src) {
			dst = make([]byte, len(src))
		}
		if len(dst) != len(src) {
			dst = dst[:len(src)]
		}

		mode := cipher.NewCBCDecrypter(block, dataIV[:])
		mode.CryptBlocks(dst, src)
	}

	return dst, nil
}
