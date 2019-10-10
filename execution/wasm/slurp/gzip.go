package slurp

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func MustGunzip(bs []byte) []byte {
	bs, err := Gunzip(bs)
	if err != nil {
		panic(err)
	}
	return bs
}

func Gunzip(bs []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	err = r.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(r)
}

func Gzip(bs []byte) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := gzip.NewWriter(buf)
	_, err := w.Write(bs)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
