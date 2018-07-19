package keys

import (
	"context"
	"crypto/sha256"
	"fmt"
	"hash"
	"net"
	"strings"

	"github.com/hyperledger/burrow/crypto"
	"github.com/hyperledger/burrow/logging"
	"github.com/tmthrgd/go-hex"
	"golang.org/x/crypto/ripemd160"
	"google.golang.org/grpc"
)

//------------------------------------------------------------------------
// all cli commands pass through the http KeyStore
// the KeyStore process also maintains the unlocked accounts

func StartStandAloneServer(keysDir, host, port string, AllowBadFilePermissions bool, logger *logging.Logger) error {
	listen, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	RegisterKeysServer(grpcServer, NewKeyStore(keysDir, AllowBadFilePermissions, logger))
	return grpcServer.Serve(listen)
}

//------------------------------------------------------------------------
// handlers

func (k *KeyStore) GenerateKey(ctx context.Context, in *GenRequest) (*GenResponse, error) {
	curveT, err := crypto.CurveTypeFromString(in.CurveType)
	if err != nil {
		return nil, err
	}

	key, err := k.Gen(in.Passphrase, curveT)
	if err != nil {
		return nil, fmt.Errorf("error generating key %s %s", curveT, err)
	}

	addrH := key.Address.String()
	if in.KeyName != "" {
		err = coreNameAdd(k.keysDirPath, in.KeyName, addrH)
		if err != nil {
			return nil, err
		}
	}

	return &GenResponse{Address: addrH}, nil
}

func (k *KeyStore) Export(ctx context.Context, in *ExportRequest) (*ExportResponse, error) {
	addr, err := getNameAddr(k.keysDirPath, in.GetName(), in.GetAddress())

	if err != nil {
		return nil, err
	}

	addrB, err := crypto.AddressFromHexString(addr)
	if err != nil {
		return nil, err
	}

	// No phrase needed for public key. I hope.
	key, err := k.GetKey(in.GetPassphrase(), addrB.Bytes())
	if err != nil {
		return nil, err
	}

	return &ExportResponse{
		Address:    addrB[:],
		CurveType:  key.CurveType.String(),
		Publickey:  key.PublicKey.Key[:],
		Privatekey: key.PrivateKey.Key[:],
	}, nil
}

func (k *KeyStore) PublicKey(ctx context.Context, in *PubRequest) (*PubResponse, error) {
	addr, err := getNameAddr(k.keysDirPath, in.GetName(), in.GetAddress())
	if err != nil {
		return nil, err
	}

	addrB, err := crypto.AddressFromHexString(addr)
	if err != nil {
		return nil, err
	}

	// No phrase needed for public key. I hope.
	key, err := k.GetKey("", addrB.Bytes())
	if key == nil {
		return nil, err
	}

	return &PubResponse{CurveType: key.CurveType.String(), PublicKey: key.Pubkey()}, nil
}

func (k *KeyStore) Sign(ctx context.Context, in *SignRequest) (*SignResponse, error) {
	addr, err := getNameAddr(k.keysDirPath, in.GetName(), in.GetAddress())
	if err != nil {
		return nil, err
	}

	addrB, err := crypto.AddressFromHexString(addr)
	if err != nil {
		return nil, err
	}

	key, err := k.GetKey(in.GetPassphrase(), addrB[:])
	if err != nil {
		return nil, err
	}

	sig, err := key.Sign(in.GetMessage())

	return &SignResponse{Signature: sig, CurveType: key.CurveType.String()}, nil
}

func (k *KeyStore) Verify(ctx context.Context, in *VerifyRequest) (*VerifyResponse, error) {
	if in.GetPublicKey() == nil {
		return nil, fmt.Errorf("must provide a pubkey")
	}
	if in.GetMessage() == nil {
		return nil, fmt.Errorf("must provide a message")
	}
	if in.GetSignature() == nil {
		return nil, fmt.Errorf("must provide a signature")
	}

	curveT, err := crypto.CurveTypeFromString(in.GetCurveType())
	if err != nil {
		return nil, err
	}
	sig, err := crypto.SignatureFromBytes(in.GetSignature(), curveT)
	if err != nil {
		return nil, err
	}
	pubkey, err := crypto.PublicKeyFromBytes(in.GetPublicKey(), curveT)
	if err != nil {
		return nil, err
	}
	match := pubkey.Verify(in.GetMessage(), sig)
	if !match {
		return nil, fmt.Errorf("Signature does not match")
	}

	return &VerifyResponse{}, nil
}

func (k *KeyStore) Hash(ctx context.Context, in *HashRequest) (*HashResponse, error) {
	var hasher hash.Hash
	switch in.GetHashtype() {
	case "ripemd160":
		hasher = ripemd160.New()
	case "sha256":
		hasher = sha256.New()
	// case "sha3":
	default:
		return nil, fmt.Errorf("Unknown hash type %v", in.GetHashtype())
	}

	hasher.Write(in.GetMessage())

	return &HashResponse{Hash: hex.EncodeUpperToString(hasher.Sum(nil))}, nil
}

func (k *KeyStore) ImportJSON(ctx context.Context, in *ImportJSONRequest) (*ImportResponse, error) {
	keyJSON := []byte(in.GetJSON())
	var err error
	addr := IsValidKeyJson(keyJSON)
	if addr != nil {
		_, err = writeKey(k.keysDirPath, addr, keyJSON)
	} else {
		err = fmt.Errorf("invalid json key passed on command line")
	}
	if err != nil {
		return nil, err
	}
	return &ImportResponse{Address: hex.EncodeUpperToString(addr)}, nil
}

func (k *KeyStore) Import(ctx context.Context, in *ImportRequest) (*ImportResponse, error) {
	curveT, err := crypto.CurveTypeFromString(in.GetCurveType())
	if err != nil {
		return nil, err
	}
	key, err := NewKeyFromPriv(curveT, in.GetKeyBytes())
	if err != nil {
		return nil, err
	}

	// store the new key
	if err = k.StoreKey(in.GetPassphrase(), key); err != nil {
		return nil, err
	}

	if in.GetName() != "" {
		if err := coreNameAdd(k.keysDirPath, in.GetName(), key.Address.String()); err != nil {
			return nil, err
		}
	}
	return &ImportResponse{Address: hex.EncodeUpperToString(key.Address[:])}, nil
}

func (k *KeyStore) List(ctx context.Context, in *ListRequest) (*ListResponse, error) {
	names, err := coreNameList(k.keysDirPath)
	if err != nil {
		return nil, err
	}

	var list []*KeyID

	for name, addr := range names {
		list = append(list, &KeyID{KeyName: name, Address: addr})
	}

	return &ListResponse{Key: list}, nil
}

func (k *KeyStore) RemoveName(ctx context.Context, in *RemoveNameRequest) (*RemoveNameResponse, error) {
	if in.GetKeyName() == "" {
		return nil, fmt.Errorf("please specify a name")
	}

	return &RemoveNameResponse{}, coreNameRm(k.keysDirPath, in.GetKeyName())
}

func (k *KeyStore) AddName(ctx context.Context, in *AddNameRequest) (*AddNameResponse, error) {
	if in.GetKeyname() == "" {
		return nil, fmt.Errorf("please specify a name")
	}

	if in.GetAddress() == "" {
		return nil, fmt.Errorf("please specify an address")
	}

	return &AddNameResponse{}, coreNameAdd(k.keysDirPath, in.GetKeyname(), strings.ToUpper(in.GetAddress()))
}
