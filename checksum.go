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
	// MD5 is the MD5 algorithm.
	MD5 = iota
	// SHA1 is the SHA1 algorithm.
	SHA1
	// SHA256 is the SHA256 algorithm.
	SHA256
	// SHA512 is the SHA512 algorithm.
	SHA512
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
	switch chksum.Algo {
	case MD5:
		return compare(chksum.Value, calcMD5(data))
	case SHA1:
		return compare(chksum.Value, calcSHA1(data))
	case SHA256:
		return compare(chksum.Value, calcSHA256(data))
	case SHA512:
		return compare(chksum.Value, calcSHA512(data))
	default:
		return ErrChecksumUnknown
	}
}

func calcMD5(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func calcSHA1(data []byte) []byte {
	checksum := sha1.Sum(data)
	return checksum[:]
}

func calcSHA256(data []byte) []byte {
	checksum := sha256.Sum256(data)
	return checksum[:]
}

func calcSHA512(data []byte) []byte {
	checksum := sha512.Sum512(data)
	return checksum[:]
}

func compare(value string, checksum []byte) error {
	if value != hex.EncodeToString(checksum[:]) {
		return ErrChecksumMismatch
	}

	return nil
}
