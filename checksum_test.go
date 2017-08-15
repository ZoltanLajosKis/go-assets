package assets

import (
	"testing"
)

func TestChecksumNil(t *testing.T) {
	err := verifyChecksum(nil, []byte("Assets"))
	assertEqual(t, err, nil)
}

func TestChecksumMD5(t *testing.T) {
	err := verifyChecksum(&Checksum{ChecksumMD5, "9aedeaf1f77b8642abe528503b8c5de8"}, []byte("Assets"))
	assertEqual(t, err, nil)

	err = verifyChecksum(&Checksum{ChecksumMD5, "1234567890abcdefghijklmnopqrstuv"}, []byte("Assets"))
	assertEqual(t, err, ErrChecksumMismatch)
}

func TestChecksumUnknown(t *testing.T) {
	err := verifyChecksum(&Checksum{-1, "12345678"}, []byte("Assets"))
	assertEqual(t, err, ErrChecksumUnknown)
}
