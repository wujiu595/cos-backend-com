package users

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestPrivateKey(t *testing.T) {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	data := []byte("hello")
	hash := crypto.Keccak256Hash(data)

	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	signatureString := hexutil.Encode(signature)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v/n%v", signatureString, publicKeyBytes)
}
