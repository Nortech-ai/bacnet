package plumbing

import (
	"encoding/binary"
	"fmt"

	"github.com/Nortech-ai/bacnet/common"
)

// NPDU is a Network Protocol Data Units.
type NPDU struct {
	Version uint8
	Control uint8
	DNET    uint16
	DLEN    uint8
	SNET    uint16
	SLEN    uint8
	SADR    uint8
	Hop     uint8
}

// NewNPDU creates a NPDU.
func NewNPDU(nsduContain bool, dstSpecifier bool, srcSpecifier bool, expectingReply bool) *NPDU {
	n := &NPDU{
		Version: 1,
	}
	n.SetControlFlags(nsduContain, dstSpecifier, srcSpecifier, expectingReply)
	return n
}

// SetControlFlags sets control flags to NPDU.
func (n *NPDU) SetControlFlags(nsduContain bool, dstSpecifier bool, srcSpecifier bool, expectingReply bool) {
	n.Control = uint8(
		common.BoolToInt(nsduContain)<<7 | common.BoolToInt(dstSpecifier)<<5 |
			common.BoolToInt(srcSpecifier)<<3 | common.BoolToInt(expectingReply)<<2,
	)
}

// UnmarshalBinary sets the values retrieved from byte sequence in a NPDU frame.
func (n *NPDU) UnmarshalBinary(b []byte) error {
	if l := len(b); l < n.MarshalLen() {
		return fmt.Errorf(
			"failed to unmarshal NPDU - marshal length %d binary length %d: %v",
			n.MarshalLen(), l,
			common.ErrTooShortToParse,
		)
	}

	n.Version = b[0]
	n.Control = b[1]

	offset := 2
	if n.flagDNET() {
		n.DNET = binary.BigEndian.Uint16(b[offset : offset+2])
		n.DLEN = b[offset+2]
		offset += 3
	}

	if n.flagSNET() {
		n.SNET = binary.BigEndian.Uint16(b[offset : offset+2])
		n.SLEN = b[offset+2]
		n.SADR = b[offset+3]
		offset += 4
	}

	if n.flagDNET() {
		n.Hop = b[offset]
		offset++
	}

	return nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (n *NPDU) MarshalTo(b []byte) error {
	if len(b) < n.MarshalLen() {
		return fmt.Errorf(
			"failed to marshall NPDU - marshal length %d binary length %d: %v",
			n.MarshalLen(),
			len(b),
			common.ErrTooShortToMarshalBinary,
		)
	}

	b[0] = n.Version
	b[1] = n.Control

	offset := 2

	if n.flagDNET() {
		binary.BigEndian.PutUint16(b[offset:offset+2], n.DNET)
		b[offset+2] = n.DLEN
		offset += 3
	}

	if n.flagSNET() {
		binary.BigEndian.PutUint16(b[offset:offset+2], n.SNET)
		b[offset+2] = n.SLEN
		b[offset+3] = n.SADR
		offset += 4
	}

	if n.flagDNET() {
		b[offset] = n.Hop
		offset++
	}

	return nil
}

const npduLenMin = 2

// MarshalLen returns the serial length of NPDU.
func (n *NPDU) MarshalLen() int {
	len := npduLenMin

	if n.flagDNET() {
		len += 4
	}

	if n.flagSNET() {
		len += 4
	}

	return len
}

func (n *NPDU) flagDNET() bool {
	return (n.Control & 0x20) != 0
}

func (n *NPDU) flagSNET() bool {
	return (n.Control & 0x08) != 0
}
