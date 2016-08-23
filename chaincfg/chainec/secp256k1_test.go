// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainec

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestGeneralSecp256k1(t *testing.T) {
	// Sample expanded pubkey (Satoshi from Genesis block)
	samplePubkey, _ := hex.DecodeString("04" +
		"678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb6" +
		"49f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5f")
	_, err := Secp256k1.ParsePubKey(samplePubkey)
	if err != nil {
		t.Errorf("failure parsing pubkey: %v", err)
	}

	// Sample compressed pubkey
	samplePubkey, _ = hex.DecodeString("02" +
		"4627032575180c2773b3eedd3a163dc2f3c6c84f9d0a1fc561a9578a15e6d0e3")
	_, err = Secp256k1.ParsePubKey(samplePubkey)
	if err != nil {
		t.Errorf("failure parsing pubkey: %v", err)
	}

	// Sample signature from https://en.bitcoin.it/wiki/Transaction
	sampleSig, _ := hex.DecodeString("30" +
		"45" +
		"02" +
		"20" +
		"6e21798a42fae0e854281abd38bacd1aeed3ee3738d9e1446618c4571d1090db" +
		"02" +
		"21" +
		"00e2ac980643b0b82c0e88ffdfec6b64e3e6ba35e7ba5fdd7d5d6cc8d25c6b2415")
	_, err = Secp256k1.ParseDERSignature(sampleSig)
	if err != nil {
		t.Errorf("failure parsing DER signature: %v", err)
	}
}

type signatureTest struct {
	name    string
	sig     []byte
	der     bool
	isValid bool
}

// decodeHex decodes the passed hex string and returns the resulting bytes.  It
// panics if an error occurs.  This is only used in the tests as a helper since
// the only way it can fail is if there is an error in the test source code.
func decodeHex(hexStr string) []byte {
	b, err := hex.DecodeString(hexStr)
	if err != nil {
		panic("invalid hex string in test source: err " + err.Error() +
			", hex: " + hexStr)
	}

	return b
}

type pubKeyTest struct {
	name    string
	key     []byte
	format  byte
	isValid bool
}

const (
	TstPubkeyUncompressed byte = 0x4 // x coord + y coord
	TstPubkeyCompressed   byte = 0x2 // y_bit + x coord
	TstPubkeyHybrid       byte = 0x6 // y_bit + x coord + y coord
)

var pubKeyTests = []pubKeyTest{
	// pubkey from bitcoin blockchain tx
	// 0437cd7f8525ceed2324359c2d0ba26006d92d85
	{
		name: "uncompressed ok",
		key: []byte{0x04, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: true,
		format:  TstPubkeyUncompressed,
	},
	{
		name: "uncompressed as hybrid ok",
		key: []byte{0x07, 0x11, 0xdb, 0x93, 0xe1, 0xdc, 0xdb, 0x8a,
			0x01, 0x6b, 0x49, 0x84, 0x0f, 0x8c, 0x53, 0xbc, 0x1e,
			0xb6, 0x8a, 0x38, 0x2e, 0x97, 0xb1, 0x48, 0x2e, 0xca,
			0xd7, 0xb1, 0x48, 0xa6, 0x90, 0x9a, 0x5c, 0xb2, 0xe0,
			0xea, 0xdd, 0xfb, 0x84, 0xcc, 0xf9, 0x74, 0x44, 0x64,
			0xf8, 0x2e, 0x16, 0x0b, 0xfa, 0x9b, 0x8b, 0x64, 0xf9,
			0xd4, 0xc0, 0x3f, 0x99, 0x9b, 0x86, 0x43, 0xf6, 0x56,
			0xb4, 0x12, 0xa3,
		},
		isValid: true,
		format:  TstPubkeyHybrid,
	},
	// from tx 0b09c51c51ff762f00fb26217269d2a18e77a4fa87d69b3c363ab4df16543f20
	{
		name: "compressed ok (ybit = 0)",
		key: []byte{0x02, 0xce, 0x0b, 0x14, 0xfb, 0x84, 0x2b, 0x1b,
			0xa5, 0x49, 0xfd, 0xd6, 0x75, 0xc9, 0x80, 0x75, 0xf1,
			0x2e, 0x9c, 0x51, 0x0f, 0x8e, 0xf5, 0x2b, 0xd0, 0x21,
			0xa9, 0xa1, 0xf4, 0x80, 0x9d, 0x3b, 0x4d,
		},
		isValid: true,
		format:  TstPubkeyCompressed,
	},
	// from tx fdeb8e72524e8dab0da507ddbaf5f88fe4a933eb10a66bc4745bb0aa11ea393c
	{
		name: "compressed ok (ybit = 1)",
		key: []byte{0x03, 0x26, 0x89, 0xc7, 0xc2, 0xda, 0xb1, 0x33,
			0x09, 0xfb, 0x14, 0x3e, 0x0e, 0x8f, 0xe3, 0x96, 0x34,
			0x25, 0x21, 0x88, 0x7e, 0x97, 0x66, 0x90, 0xb6, 0xb4,
			0x7f, 0x5b, 0x2a, 0x4b, 0x7d, 0x44, 0x8e,
		},
		isValid: true,
		format:  TstPubkeyCompressed,
	},
	{
		name: "hybrid",
		key: []byte{0x06, 0x79, 0xbe, 0x66, 0x7e, 0xf9, 0xdc, 0xbb,
			0xac, 0x55, 0xa0, 0x62, 0x95, 0xce, 0x87, 0x0b, 0x07,
			0x02, 0x9b, 0xfc, 0xdb, 0x2d, 0xce, 0x28, 0xd9, 0x59,
			0xf2, 0x81, 0x5b, 0x16, 0xf8, 0x17, 0x98, 0x48, 0x3a,
			0xda, 0x77, 0x26, 0xa3, 0xc4, 0x65, 0x5d, 0xa4, 0xfb,
			0xfc, 0x0e, 0x11, 0x08, 0xa8, 0xfd, 0x17, 0xb4, 0x48,
			0xa6, 0x85, 0x54, 0x19, 0x9c, 0x47, 0xd0, 0x8f, 0xfb,
			0x10, 0xd4, 0xb8,
		},
		format:  TstPubkeyHybrid,
		isValid: true,
	},
}

func TestPubKeys(t *testing.T) {
	for _, test := range pubKeyTests {
		pk, err := Secp256k1.ParsePubKey(test.key)
		if err != nil {
			if test.isValid {
				t.Errorf("%s pubkey failed when shouldn't %v",
					test.name, err)
			}
			continue
		}
		if !test.isValid {
			t.Errorf("%s counted as valid when it should fail",
				test.name)
			continue
		}
		var pkStr []byte
		switch test.format {
		case TstPubkeyUncompressed:
			pkStr = (PublicKey)(pk).SerializeUncompressed()
		case TstPubkeyCompressed:
			pkStr = (PublicKey)(pk).SerializeCompressed()
		case TstPubkeyHybrid:
			pkStr = (PublicKey)(pk).SerializeHybrid()
		}
		if !bytes.Equal(test.key, pkStr) {
			t.Errorf("%s pubkey: serialized keys do not match.",
				test.name)
			spew.Dump(test.key)
			spew.Dump(pkStr)
		}
	}
}

func TestPrivKeys(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
	}{
		{
			name: "check curve",
			key: []byte{
				0xea, 0xf0, 0x2c, 0xa3, 0x48, 0xc5, 0x24, 0xe6,
				0x39, 0x26, 0x55, 0xba, 0x4d, 0x29, 0x60, 0x3c,
				0xd1, 0xa7, 0x34, 0x7d, 0x9d, 0x65, 0xcf, 0xe9,
				0x3c, 0xe1, 0xeb, 0xff, 0xdc, 0xa2, 0x26, 0x94,
			},
		},
	}

	for _, test := range tests {
		priv, pub := Secp256k1.PrivKeyFromBytes(test.key)
		_, err := Secp256k1.ParsePubKey(pub.SerializeUncompressed())
		if err != nil {
			t.Errorf("%s privkey: %v", test.name, err)
			continue
		}

		hash := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9}
		r, s, err := Secp256k1.Sign(priv, hash)
		if err != nil {
			t.Errorf("%s could not sign: %v", test.name, err)
			continue
		}
		sig := Secp256k1.NewSignature(r, s)

		if !Secp256k1.Verify(pub, hash, sig.GetR(), sig.GetS()) {
			t.Errorf("%s could not verify: %v", test.name, err)
			continue
		}

		serializedKey := priv.Serialize()
		if !bytes.Equal(serializedKey, test.key) {
			t.Errorf("%s unexpected serialized bytes - got: %x, "+
				"want: %x", test.name, serializedKey, test.key)
		}
	}
}

var signatureTests = []signatureTest{
	// signatures from bitcoin blockchain tx
	// 0437cd7f8525ceed2324359c2d0ba26006d92d85
	{
		name: "valid signature.",
		sig: []byte{0x30, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: true,
	},
	{
		name:    "empty.",
		sig:     []byte{},
		isValid: false,
	},
	{
		name: "bad magic.",
		sig: []byte{0x31, 0x44, 0x02, 0x20, 0x4e, 0x45, 0xe1, 0x69,
			0x32, 0xb8, 0xaf, 0x51, 0x49, 0x61, 0xa1, 0xd3, 0xa1,
			0xa2, 0x5f, 0xdf, 0x3f, 0x4f, 0x77, 0x32, 0xe9, 0xd6,
			0x24, 0xc6, 0xc6, 0x15, 0x48, 0xab, 0x5f, 0xb8, 0xcd,
			0x41, 0x02, 0x20, 0x18, 0x15, 0x22, 0xec, 0x8e, 0xca,
			0x07, 0xde, 0x48, 0x60, 0xa4, 0xac, 0xdd, 0x12, 0x90,
			0x9d, 0x83, 0x1c, 0xc5, 0x6c, 0xbb, 0xac, 0x46, 0x22,
			0x08, 0x22, 0x21, 0xa8, 0x76, 0x8d, 0x1d, 0x09,
		},
		der:     true,
		isValid: false,
	},
}

func TestSignatures(t *testing.T) {
	for _, test := range signatureTests {
		var err error
		if test.der {
			_, err = Secp256k1.ParseDERSignature(test.sig)
		} else {
			_, err = Secp256k1.ParseSignature(test.sig)
		}
		if err != nil {
			if test.isValid {
				t.Errorf("%s signature failed when shouldn't %v",
					test.name, err)
			}
			continue
		}
		if !test.isValid {
			t.Errorf("%s counted as valid when it should fail",
				test.name)
		}
	}
}
