package host

import (
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"
	ecies "github.com/ecies/go/v2"
	seal "github.com/rubixchain/rubixgoplatform/crypto"
	secp256k1 "github.com/decred/dcrd/dcrec/secp256k1/v4"
)

func ConvertpubkeySecp256k1ToEcies(pubKey *secp256k1.PublicKey) (*ecies.PublicKey, error) {
	// Extract the X and Y coordinates by calling the functions
	x := pubKey.X()
	y := pubKey.Y()

	// Create an ECIES public key from the X and Y coordinates
	eciesPubKey := &ecies.PublicKey{
		X:     x,
		Y:     y,
		Curve: secp256k1.S256(),
	}

	return eciesPubKey, nil
}

// ConvertSecp256k1ToEcies converts a secp256k1 private key to an ECIES private key.
func ConvertSecp256k1privkeyToEcies(privKey *secp256k1.PrivateKey) (*ecies.PrivateKey, error) {
	// Serialize the private key to get the private scalar bytes
	privKeyBytes := privKey.Serialize()

	// Convert the private scalar bytes to a big.Int
	d := new(big.Int).SetBytes(privKeyBytes)
	// Create an ECIES public key from the secp256k1 public key
	pubKey := privKey.PubKey()
	eciesPubKey := &ecies.PublicKey{
		X:     pubKey.X(),
		Y:     pubKey.Y(),
		Curve: secp256k1.S256(),
	}

	// Create an ECIES private key from the D value and the ECIES public key
	eciesPrivKey := &ecies.PrivateKey{
		PublicKey: eciesPubKey,
		D:         d,
	}

	return eciesPrivKey, nil
}

func EciesEncryption(pubkey_path string, data []byte) (ciphertext []byte) {
	read_pubKey, err := os.ReadFile(pubkey_path)
	if err != nil {
		log.Fatal(err)
	}
	
	pemdecoded_pubkey, _ := pem.Decode(read_pubKey)
	pubkeyback, _ := secp256k1.ParsePubKey(pemdecoded_pubkey.Bytes)
	
	eciesPubKey, err := ConvertpubkeySecp256k1ToEcies(pubkeyback)
	if err != nil {
		fmt.Println("Error converting public key:", err)
		return
	}

	ciphertext, err = ecies.Encrypt(eciesPubKey, data)
	if err != nil {
		panic(err)
	}
	fmt.Println("ciphertext is  ", ciphertext)
	return ciphertext
}
func EciesDecryption(privkey_path string, encrypted_data []byte) (plaintext string, err error) {
	read_encodedprivkey, err := os.ReadFile(privkey_path)
	if err != nil {
		return "", err
	}
	
	pemdecoded_privkey, _ := pem.Decode(read_encodedprivkey)
	password := "mypassword"
	unsealedprivkey, err := seal.UnSeal(password, (pemdecoded_privkey).Bytes)
	if err != nil {
		return "", fmt.Errorf("unable to decrypt the private key: %v", err)
	}
	fmt.Println("Decrypted Private key is ", unsealedprivkey)
	parsedprivkey := secp256k1.PrivKeyFromBytes(unsealedprivkey)

	ecies_privkey, err := ConvertSecp256k1privkeyToEcies(parsedprivkey)
	if err != nil {
		return "", fmt.Errorf("ConvertSecp256k1privkeyToEcies func failed: %v", err)
	}

	plaintext_bytes, err := ecies.Decrypt(ecies_privkey, encrypted_data)
	if err != nil {
		return "", fmt.Errorf("decrypt func failed: %v", err)
	}

	plaintext_string := string(plaintext_bytes)
	return plaintext_string, nil
}

