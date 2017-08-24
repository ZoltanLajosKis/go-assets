package assets

import (
	"testing"
)

func TestChecksumMD5(t *testing.T) {
	err := verifyChecksum(&Checksum{MD5,
		"9aedeaf1f77b8642abe528503b8c5de8"}, []byte("Assets"))
	assertEqual(t, err, nil)

	err = verifyChecksum(&Checksum{MD5,
		"1234567890abcdef1234567890abcdef"}, []byte("Assets"))
	assertEqual(t, err, ErrChecksumMismatch)
}

func TestChecksumSHA1(t *testing.T) {
	err := verifyChecksum(&Checksum{SHA1,
		"20e338624cee29d0effead85b0dd0e70de783b4c"}, []byte("Assets"))
	assertEqual(t, err, nil)

	err = verifyChecksum(&Checksum{SHA1,
		"1234567890abcdef1234567890abcdef12345678"}, []byte("Assets"))
	assertEqual(t, err, ErrChecksumMismatch)
}

func TestChecksumSHA256(t *testing.T) {
	err := verifyChecksum(&Checksum{SHA256,
		"bd12731d7bc9b843d8523e654ae92abe735ee95f0777e46e77ee286b17833acd"},
		[]byte("Assets"))
	assertEqual(t, err, nil)

	err = verifyChecksum(&Checksum{SHA256,
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"},
		[]byte("Assets"))
	assertEqual(t, err, ErrChecksumMismatch)
}

func TestChecksumSHA512(t *testing.T) {
	err := verifyChecksum(&Checksum{SHA512,
		"775b844893be4c703d6b67af589c71203ff4137950e5ae27946a152e9e10a5cac4b2355a682ec515e9a919e699ccf2f5255b68a7f2b6026c5a173bad84047b4a"},
		[]byte("Assets"))
	assertEqual(t, err, nil)

	err = verifyChecksum(&Checksum{SHA512,
		"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"},
		[]byte("Assets"))
	assertEqual(t, err, ErrChecksumMismatch)
}

func TestChecksumUnknown(t *testing.T) {
	err := verifyChecksum(&Checksum{-1, "12345678"}, []byte("Assets"))
	assertEqual(t, err, ErrChecksumUnknown)
}
