package plumbing_test

import (
	"bytes"
	"testing"

	"github.com/Nortech-ai/bacnet/plumbing"
	. "github.com/Nortech-ai/bacnet/test_utils"
)

func TestNPDUUnmarshall_Simple(t *testing.T) {
	npdu := plumbing.NewNPDU(false, false, false, false)
	// This NPDU is version 1, no control flags
	err := npdu.UnmarshalBinary([]byte{0x01, 0x00})
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual(t, uint8(1), npdu.Version)
	AssertEqual(t, uint8(0), npdu.Control)
	// AssertEqual(t, uint16(0),  npdu.SNET)
	AssertEqual(t, uint16(0), npdu.DNET)
	AssertEqual(t, uint8(0), npdu.DLEN)
	AssertEqual(t, uint8(0), npdu.Hop)
}

func TestNPDUUnmarshall_WithDNET(t *testing.T) {
	npdu := plumbing.NewNPDU(false, false, false, false)
	// This NPDU is version 1, dnet only (control flag 0x20)
	err := npdu.UnmarshalBinary([]byte{0x1, 0x20, 0xff, 0xff, 0x0, 0xff})
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual(t, uint8(1), npdu.Version)
	AssertEqual(t, uint8(0x20), npdu.Control)
	// AssertEqual(t, uint16(0),  npdu.SNET)
	AssertEqual(t, uint16(0xffff), npdu.DNET)
	AssertEqual(t, uint8(0), npdu.DLEN)
	AssertEqual(t, uint8(0xff), npdu.Hop)
}

func TestNPDUUnmarshall_WithSNET(t *testing.T) {
	npdu := plumbing.NewNPDU(false, false, false, false)
	// This NPDU is version 1, snet only (control flag 0x08)
	err := npdu.UnmarshalBinary([]byte{0x1, 0x8, 0x0, 0x8, 0x1, 0x8})
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual(t, uint8(1), npdu.Version)
	AssertEqual(t, uint8(0x08), npdu.Control)
	AssertEqual(t, uint16(0x0008), npdu.SNET)
	AssertEqual(t, uint8(1), npdu.SLEN)
	AssertEqual(t, uint8(8), npdu.SADR)
	AssertEqual(t, uint16(0), npdu.DNET)
	AssertEqual(t, uint8(0), npdu.DLEN)
	AssertEqual(t, uint8(0), npdu.Hop)
}

func TestNPDUUnmarshall_WithSNETAndDNET(t *testing.T) {
	npdu := plumbing.NewNPDU(false, false, false, false)
	// This NPDU is version 1, snet and dnet (control flag 0x28)
	err := npdu.UnmarshalBinary([]byte{0x1, 0x28, 0xff, 0xff, 0x0, 0x0, 0x8, 0x1, 0x18, 0xfe})
	if err != nil {
		t.Fatal(err)
	}
	AssertEqual(t, uint8(1), npdu.Version)
	AssertEqual(t, uint8(0x28), npdu.Control)
	AssertEqual(t, uint16(0x0008), npdu.SNET)
	AssertEqual(t, uint8(1), npdu.SLEN)
	AssertEqual(t, uint8(24), npdu.SADR)
	AssertEqual(t, uint16(0xffff), npdu.DNET)
	AssertEqual(t, uint8(0), npdu.DLEN)
	AssertEqual(t, uint8(0xfe), npdu.Hop)
}

func TestNPDUMarshall_WithSNETAndDNET(t *testing.T) {
	npdu := plumbing.NewNPDU(false, false, false, false)

	npdu.Control = 0x28
	npdu.SNET = 0x0008
	npdu.SLEN = 1
	npdu.SADR = 24
	npdu.DNET = 0xffff
	npdu.DLEN = 0
	npdu.Hop = 0xfe

	// This NPDU is version 1, snet and dnet (control flag 0x28)
	b := make([]byte, npdu.MarshalLen())
	err := npdu.MarshalTo(b)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x1, 0x28, 0xff, 0xff, 0x0, 0x0, 0x8, 0x1, 0x18, 0xfe}
	if !bytes.Equal(expected, b) {
		t.Errorf("Expected %v, got %v", expected, b)
	}
}
