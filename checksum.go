package assets

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
)

// ChecksumAlgo enumerates checksum algorihms.
type ChecksumAlgo int

const (
	// ChecksumMD5 is the MD5 algorithm.
	ChecksumMD5 = iota
	// ChecksumSHA1 is the SHA1 algorithm.
	ChecksumSHA1
	// ChecksumSHA256 is the SHA256 algorithm.
	ChecksumSHA256
	// ChecksumSHA512 is the SHA512 algorithm.
	ChecksumSHA512
)

var (
	// ErrChecksumMismatch is returned when the expected checksum differs from the calculated one
	ErrChecksumMismatch = errors.New("checksum mismatch")
	// ErrChecksumUnknown is returned when an invalid checksum algorithm is specified
	ErrChecksumUnknown = errors.New("unknown checksum algorithm")
)

// Checksum describes a checksum verification for an asset source.
type Checksum struct {
	Algo  ChecksumAlgo
	Value string
}

func verifyChecksum(chksum *Checksum, data []byte) error {
	if chksum == nil {
		return nil
	}

	switch chksum.Algo {
	case ChecksumMD5:
		return verifyChecksumMD5(chksum.Value, data)
	case ChecksumSHA1:
		return verifyChecksumSHA1(chksum.Value, data)
	case ChecksumSHA256:
		return verifyChecksumSHA256(chksum.Value, data)
	case ChecksumSHA512:
		return verifyChecksumSHA512(chksum.Value, data)
	default:
		return ErrChecksumUnknown
	}
}

func verifyChecksumMD5(value string, data []byte) error {
	h := md5.New()
	h.Write(data)
	if value != hex.EncodeToString(h.Sum(nil)) {
		return ErrChecksumMismatch
	}

	return nil
}

func verifyChecksumSHA1(value string, data []byte) error {
	checksum := sha1.Sum(data)
	if value != hex.EncodeToString(checksum[:]) {
		return ErrChecksumMismatch
	}

	return nil
}

func verifyChecksumSHA256(value string, data []byte) error {
	checksum := sha256.Sum256(data)
	if value != hex.EncodeToString(checksum[:]) {
		return ErrChecksumMismatch
	}

	return nil
}

func verifyChecksumSHA512(value string, data []byte) error {
	checksum := sha512.Sum512(data)
	if value != hex.EncodeToString(checksum[:]) {
		return ErrChecksumMismatch
	}

	return nil
}
