// Copyright 2017 Monax Industries Limited
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v0

import (
	"io"
	"io/ioutil"

	wire "github.com/tendermint/go-wire"

	rpc "github.com/hyperledger/burrow/rpc"
)

// Codec that uses tendermints 'binary' package for JSON.
type TCodec struct {
}

// Get a new codec.
func NewTCodec() rpc.Codec {
	return &TCodec{}
}

// Encode to an io.Writer.
func (tCodec *TCodec) Encode(v interface{}, w io.Writer) error {
	var err error
	var n int
	wire.WriteJSON(v, w, &n, &err)
	return err
}

// Encode to a byte array.
func (tCodec *TCodec) EncodeBytes(v interface{}) ([]byte, error) {
	return wire.JSONBytes(v), nil
}

// TODO: [ben] implement EncodeBytesPtr ?

// Decode from an io.Reader.
func (tCodec *TCodec) Decode(v interface{}, r io.Reader) error {
	bts, errR := ioutil.ReadAll(r)
	if errR != nil {
		return errR
	}
	var err error
	wire.ReadJSON(v, bts, &err)
	return err
}

// Decode from a byte array.
func (tCodec *TCodec) DecodeBytes(v interface{}, bts []byte) error {
	var err error
	wire.ReadJSON(v, bts, &err)
	return err
}

// Decode from a byte array pointer.
func (tCodec *TCodec) DecodeBytesPtr(v interface{}, bts []byte) error {
	var err error
	wire.ReadJSONPtr(v, bts, &err)
	return err
}
