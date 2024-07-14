package wallet

import (
	"blockchain-app/encoders"
	"blockchain-app/handlers"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"log"
	"math/big"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) GetAddress() string {
	pubHash := getPublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubHash...)
	checksumHash := getChecksum(versionedHash)
	fullHash := append(versionedHash, checksumHash...)
	walletAddress := encoders.Base58Encode(fullHash)

	return walletAddress
}

func ValidateAddress(address string) bool {
	pubKeyHash := encoders.Base58Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := getChecksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func CreateAndSaveNewWallet(nodeID string) string {
	wallets := GetWallets(nodeID)
	address := wallets.AddNewWallet()
	wallets.SaveFile(nodeID)

	return address
}

func createEcdsaKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func createNewWallet() *Wallet {
	private, public := createEcdsaKeyPair()
	wallet := Wallet{private, public}

	return &wallet
}

func getPublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		handlers.HandleErrors(err)
	}

	return hasher.Sum(nil)
}

func getChecksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}

func (w *Wallet) MarshalJSON() ([]byte, error) {

	curveName := func(c elliptic.Curve) string {
		switch c {
		case elliptic.P256():
			return "P-256"
		case elliptic.P384():
			return "P-384"
		case elliptic.P521():
			return "P-521"
		default:
			return "Unknown"
		}
	}

	mapStringAny := map[string]any{
		"PrivateKey": map[string]any{
			"D": w.PrivateKey.D,
			"PublicKey": map[string]any{
				"Curve": curveName(w.PrivateKey.PublicKey.Curve),
				"X":     w.PrivateKey.PublicKey.X,
				"Y":     w.PrivateKey.PublicKey.Y,
			},
		},
		"PublicKey": w.PublicKey,
	}
	return json.MarshalIndent(mapStringAny, "", "  ")
}

func (w *Wallet) UnmarshalJSON(data []byte) error {

	type ECPoint struct {
		Curve string
		X, Y  *big.Int
	}
	type PrivateKey struct {
		D         *big.Int
		PublicKey ECPoint
	}
	type WalletJSON struct {
		PrivateKey PrivateKey
		PublicKey  []byte
	}

	var wj WalletJSON
	if err := json.Unmarshal(data, &wj); err != nil {
		return err
	}

	getCurve := func(name string) elliptic.Curve {
		switch name {
		case "P-256":
			return elliptic.P256()
		case "P-384":
			return elliptic.P384()
		case "P-521":
			return elliptic.P521()
		default:
			return nil
		}
	}

	w.PublicKey = wj.PublicKey
	w.PrivateKey.D = wj.PrivateKey.D
	w.PrivateKey.PublicKey.Curve = getCurve(wj.PrivateKey.PublicKey.Curve)
	w.PrivateKey.PublicKey.X = wj.PrivateKey.PublicKey.X
	w.PrivateKey.PublicKey.Y = wj.PrivateKey.PublicKey.Y

	return nil
}
