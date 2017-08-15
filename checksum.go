package assets

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
)

// ChecksumAlgo enumerates checksum algorihms.
type ChecksumAlgo int

const (
	// ChecksumMD5 is the MD5 algorithm.
	ChecksumMD5 = iota
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
