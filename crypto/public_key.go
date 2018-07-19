package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	abci "github.com/tendermint/abci/types"
	tmCrypto "github.com/tendermint/go-crypto"
	"github.com/tmthrgd/go-hex"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ripemd160"
)

type PublicKeyJSON struct {
	CurveType string
	PublicKey string
}

// Returns the length in bytes of the public key
func PublicKeyLength(curveType CurveType) int {
	switch curveType {
	case CurveTypeEd25519:
		return ed25519.PublicKeySize
	case CurveTypeSecp256k1:
		return btcec.PubKeyBytesLenCompressed
	default:
		// Other functions rely on this
		return 0
	}
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	jStruct := PublicKeyJSON{
		CurveType: p.CurveType.String(),
		PublicKey: hex.EncodeUpperToString(p.Key),
	}
	txt, err := json.Marshal(jStruct)
	return txt, err
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return p.MarshalJSON()
}

func (p *PublicKey) UnmarshalJSON(text []byte) error {
	var jStruct PublicKeyJSON
	err := json.Unmarshal(text, &jStruct)
	if err != nil {
		return err
	}
	CurveType, err := CurveTypeFromString(jStruct.CurveType)
	if err != nil {
		return err
	}
	bs, err := hex.DecodeString(jStruct.PublicKey)
	if err != nil {
		return err
	}
	p.CurveType = CurveType
	p.Key = bs
	return nil
}

func (p *PublicKey) UnmarshalText(text []byte) error {
	return p.UnmarshalJSON(text)
}

func (p PublicKey) IsValid() bool {
	publicKeyLength := PublicKeyLength(p.CurveType)
	return publicKeyLength != 0 && publicKeyLength == len(p.Key)
}

func (p PublicKey) Verify(msg []byte, signature Signature) bool {
	switch p.CurveType {
	case CurveTypeEd25519:
		return ed25519.Verify(p.Key, msg, signature)
	case CurveTypeSecp256k1:
		pub, err := btcec.ParsePubKey(p.Key, btcec.S256())
		if err != nil {
			return false
		}
		sig, err := btcec.ParseDERSignature(signature, btcec.S256())
		if err != nil {
			return false
		}
		return sig.Verify(msg, pub)
	default:
		panic(fmt.Sprintf("invalid curve type"))
	}
}

func (p PublicKey) PublicKey() PublicKey {
	return p
}

func (p PublicKey) Address() Address {
	switch p.CurveType {
	case CurveTypeEd25519:
		// FIMXE: tendermint go-crypto-0.5.0 uses weird scheme, this is fixed in 0.6.0
		tmPubKey := new(tmCrypto.PubKeyEd25519)
		copy(tmPubKey[:], p.Key)
		addr, _ := AddressFromBytes(tmPubKey.Address())
		return addr
	case CurveTypeSecp256k1:
		sha := sha256.New()
		sha.Write(p.Key[:])

		hash := ripemd160.New()
		hash.Write(sha.Sum(nil))
		addr, _ := AddressFromBytes(hash.Sum(nil))
		return addr
	default:
		panic(fmt.Sprintf("unknown CurveType %d", p.CurveType))
	}
}

func (p PublicKey) AddressHashType() string {
	switch p.CurveType {
	case CurveTypeEd25519:
		return "go-crypto-0.5.0"
	case CurveTypeSecp256k1:
		return "btc"
	default:
		return ""
	}
}

func (p PublicKey) RawBytes() []byte {
	return p.Key[:]
}

// Return the ABCI PubKey. See Tendermint protobuf.go for the go-crypto conversion this is based on
func (p PublicKey) ABCIPubKey() abci.PubKey {
	switch p.CurveType {
	case CurveTypeEd25519:
		return abci.PubKey{
			Type: "ed25519",
			Data: p.RawBytes(),
		}
	case CurveTypeSecp256k1:
		return abci.PubKey{
			Type: "secp256k1",
			Data: p.RawBytes(),
		}
	default:
		return abci.PubKey{}
	}
}

func (p PublicKey) String() string {
	return hex.EncodeUpperToString(p.Key)
}

// Produces a binary encoding of the CurveType byte plus
// the public key for padded to a fixed width on the right
func (p PublicKey) Encode() []byte {
	encoded := make([]byte, PublicKeyLength(p.CurveType)+1)
	encoded[0] = p.CurveType.Byte()
	copy(encoded[1:], p.Key)
	return encoded
}

// Decodes an encoded public key returning the number of bytes consumed
func DecodePublicKeyFixedWidth(buf []byte, publicKey *PublicKey) (int, error) {
	if len(buf) < 1 {
		return 0, fmt.Errorf("encoded bytes buffer must not be empty")
	}
	curveType := CurveType(buf[0])
	publicKeyEnd := PublicKeyLength(curveType) + 1
	if publicKeyEnd <= 0 {
		return 0, fmt.Errorf("CurveType with identifier %v is unknown", curveType.Byte())
	}
	if len(buf) < publicKeyEnd {
		return 0, fmt.Errorf("encoded bytes buffer has length %v but public key encoding for %v needs %v bytes",
			len(buf), curveType, publicKeyEnd)
	}

	publicKey.CurveType = curveType
	publicKey.Key = buf[1:publicKeyEnd]
	if !publicKey.IsValid() {
		return publicKeyEnd, fmt.Errorf("decoded public key %v is not valid", publicKey)
	}
	return publicKeyEnd, nil
}
