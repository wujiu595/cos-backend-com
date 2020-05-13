package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"hash"
	"strings"
)

// use pbkdf2 encode password
func EncodePBKDF2Password(rawPwd string, salt string) string {
	pbkdf2 := PBKDF2([]byte(rawPwd), []byte(salt), 10000, 50, sha256.New)
	pbkdf2b64 := base64.URLEncoding.EncodeToString(pbkdf2)
	return fmt.Sprintf("PBKDF2:%s:%s", salt, pbkdf2b64)
}

// use pbkdf2 decode valid password
func ValidPBKDF2Password(rawPwd, encodedPwd string) bool {
	parts := strings.Split(encodedPwd, ":")
	if len(parts) != 3 {
		return false
	}
	if parts[0] != "PBKDF2" {
		return false
	}
	value := EncodePBKDF2Password(rawPwd, parts[1])
	if value != encodedPwd {
		return false
	}
	return true
}

// https://github.com/golang/crypto/blob/master/pbkdf2/pbkdf2.go
func PBKDF2(password, salt []byte, iter, keyLen int, h func() hash.Hash) []byte {
	prf := hmac.New(h, password)
	hashLen := prf.Size()
	numBlocks := (keyLen + hashLen - 1) / hashLen

	var buf [4]byte
	dk := make([]byte, 0, numBlocks*hashLen)
	U := make([]byte, hashLen)
	for block := 1; block <= numBlocks; block++ {
		// N.B.: || means concatenation, ^ means XOR
		// for each block T_i = U_1 ^ U_2 ^ ... ^ U_iter
		// U_1 = PRF(password, salt || uint(i))
		prf.Reset()
		prf.Write(salt)
		buf[0] = byte(block >> 24)
		buf[1] = byte(block >> 16)
		buf[2] = byte(block >> 8)
		buf[3] = byte(block)
		prf.Write(buf[:4])
		dk = prf.Sum(dk)
		T := dk[len(dk)-hashLen:]
		copy(U, T)

		// U_n = PRF(password, U_(n-1))
		for n := 2; n <= iter; n++ {
			prf.Reset()
			prf.Write(U)
			U = U[:0]
			U = prf.Sum(U)
			for x := range U {
				T[x] ^= U[x]
			}
		}
	}
	return dk[:keyLen]
}
