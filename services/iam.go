package services

import (
	"fmt"

	"github.com/jonalfarlinga/bacnet/common"
	"github.com/jonalfarlinga/bacnet/objects"
	"github.com/jonalfarlinga/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type UnconfirmedIAm struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type UnconfirmedIAmDec struct {
	InstanceNum           uint32
	DeviceType            uint16
	MaxAPDULength         uint16
	SegmentationSupported uint8
	VendorId              uint16
}

// IAmObjects creates an instance of UnconfirmedIAm objects.
func IAmObjects(instN uint32, acceptedSize uint16, supportedSeg uint8, vendorID uint16) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 4)

	objs[0] = objects.EncObjectIdentifier(false, objects.TagBACnetObjectIdentifier, 8, instN)
	objs[1] = objects.EncUnsignedInteger(uint(acceptedSize))
	objs[2] = objects.EncEnumerated(supportedSeg)
	objs[3] = objects.EncUnsignedInteger(uint(vendorID))

	return objs
}

// NewUnconfirmedIAm creates a UnconfirmedIam.
func NewUnconfirmedIAm(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *UnconfirmedIAm {
	u := &UnconfirmedIAm{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.UnConfirmedReq, ServiceUnconfirmedIAm, IAmObjects(1, 1024, 0, 1)),
	}
	u.SetLength()

	return u
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnconfirmedIAm frame.
func (u *UnconfirmedIAm) UnmarshalBinary(b []byte) error {
	if l := len(b); l < u.MarshalLen() {
		return fmt.Errorf(
			"failed to unmarshal UnconfirmedIAm - marshal length %d binary length %d: %v",
			u.MarshalLen(), l,
			common.ErrTooShortToParse,
		)
	}

	var offset int = 0
	if err := u.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling UnconfirmedIAm %+v: %v",
			u, common.ErrTooShortToParse,
		)
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling UnconfirmedIAm %+v: %v",
			u, common.ErrTooShortToParse,
		)
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling UnconfirmedIAm %+v: %v",
			u, common.ErrTooShortToParse,
		)
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a UnconfirmedIAm instance.
func (u *UnconfirmedIAm) MarshalBinary() ([]byte, error) {
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, fmt.Errorf("failed to marshal binary: %v", err)
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (u *UnconfirmedIAm) MarshalTo(b []byte) error {
	if len(b) < u.MarshalLen() {
		return fmt.Errorf(
			"failed to marshal UnconfirmedIAm - marshal length %d binary length %d: %v",
			u.MarshalLen(), len(b),
			common.ErrTooShortToMarshalBinary,
		)
	}
	var offset = 0
	if err := u.BVLC.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("marshalling UnconfirmedIAm: %v", err)
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("marshalling UnconfirmedIAm: %v", err)
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("marshalling UnconfirmedIAm: %v", err)
	}

	return nil
}

// MarshalLen returns the serial length of UnconfirmedIAm.
func (u *UnconfirmedIAm) MarshalLen() int {
	l := u.BVLC.MarshalLen()
	l += u.NPDU.MarshalLen()
	l += u.APDU.MarshalLen()
	return l
}

// SetLength sets the length in Length field.
func (u *UnconfirmedIAm) SetLength() {
	u.BVLC.Length = uint16(u.MarshalLen())
}

func (u *UnconfirmedIAm) Decode() (UnconfirmedIAmDec, error) {
	decIAm := UnconfirmedIAmDec{}

	if len(u.APDU.Objects) != 4 {
		return decIAm, fmt.Errorf(
			"failed to decode UnconfirmedIAm - number of objects %d: %v",
			len(u.APDU.Objects),
			common.ErrWrongObjectCount,
		)
	}

	for i, obj := range u.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decIAm, fmt.Errorf("decoding UnconfirmedIAm: %v", err)
			}
			decIAm.DeviceType = objId.ObjectType
			decIAm.InstanceNum = objId.InstanceNumber
		case 1:
			maxLen, err := objects.DecUnsignedInteger(obj)
			if err != nil {
				return decIAm, fmt.Errorf("decoding UnconfirmedIAm: %v", err)
			}
			decIAm.MaxAPDULength = uint16(maxLen)
		case 2:
			segSupport, err := objects.DecEnumerated(obj)
			if err != nil {
				return decIAm, fmt.Errorf("decoding UnconfirmedIAm: %v", err)
			}
			decIAm.SegmentationSupported = uint8(segSupport)
		case 3:
			vendorId, err := objects.DecUnsignedInteger(obj)
			if err != nil {
				return decIAm, fmt.Errorf("decoding UnconfirmedIAm: %v", err)
			}
			decIAm.VendorId = uint16(vendorId)
		}
	}

	return decIAm, nil
}

func (u *UnconfirmedIAm) GetService() uint8 {
	return u.APDU.Service
}

func (u *UnconfirmedIAm) GetType() uint8 {
	return u.APDU.Type
}

//--------------------------------------------------------------------
//---------------------Unicast implementation-------------------------
//--------------------------------------------------------------------

// UnicastIAm is a BACnet message for unicast IAm responses.
type UnicastIAm struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

// NewUnicastIAm creates a UnicastIAm.  The BVLC and NPDU should be pre-populated
// with the correct values (including the unicast BVLC function).
func NewUnicastIAm(bvlc *plumbing.BVLC, npdu *plumbing.NPDU, apdu *plumbing.APDU) *UnicastIAm {
	u := &UnicastIAm{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: apdu, //We pass apdu now
	}
	u.SetLength() // Set the BVLC length based on the entire message.

	return u
}

// UnmarshalBinary sets the values retrieved from byte sequence in a UnicastIAm frame.
func (u *UnicastIAm) UnmarshalBinary(b []byte) error {
	//This is the same as the other
	if l := len(b); l < u.MarshalLen() {
		return common.ErrTooShortToParse
	}

	var offset int = 0
	if err := u.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return common.ErrTooShortToParse
	}

	return nil
}

// MarshalBinary returns the byte sequence generated from a UnicastIAm instance.
func (u *UnicastIAm) MarshalBinary() ([]byte, error) {
	// Same as other
	b := make([]byte, u.MarshalLen())
	if err := u.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (u *UnicastIAm) MarshalTo(b []byte) error {
	//same as other
	if len(b) < u.MarshalLen() {
		return common.ErrTooShortToMarshalBinary
	}
	var offset = 0
	if err := u.BVLC.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += u.BVLC.MarshalLen()

	if err := u.NPDU.MarshalTo(b[offset:]); err != nil {
		return err
	}
	offset += u.NPDU.MarshalLen()

	if err := u.APDU.MarshalTo(b[offset:]); err != nil {
		return err
	}

	return nil
}

// MarshalLen returns the serial length of UnicastIAm.
func (u *UnicastIAm) MarshalLen() int {
	//same as other
	l := u.BVLC.MarshalLen()
	l += u.NPDU.MarshalLen()
	l += u.APDU.MarshalLen()

	return l
}

// SetLength sets the length in the BVLC Length field.
func (u *UnicastIAm) SetLength() {
	//same as other
	u.BVLC.Length = uint16(u.MarshalLen())
}

// Decode extracts the relevant fields from the UnicastIAm message.
func (u *UnicastIAm) Decode() (UnconfirmedIAmDec, error) {
	//same as other
	decIAm := UnconfirmedIAmDec{} // Use same decoding struct

	if len(u.APDU.Objects) != 4 {
		return decIAm, common.ErrWrongObjectCount
	}

	for i, obj := range u.APDU.Objects {
		switch i {
		case 0:
			objId, err := objects.DecObjectIdentifier(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.InstanceNum = objId.InstanceNumber
		case 1:
			maxLen, err := objects.DecUnsignedInteger(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.MaxAPDULength = uint16(maxLen)
		case 2:
			segSupport, err := objects.DecEnumerated(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.SegmentationSupported = uint8(segSupport)
		case 3:
			vendorId, err := objects.DecUnsignedInteger(obj)
			if err != nil {
				return decIAm, err
			}
			decIAm.VendorId = uint16(vendorId)
		}
	}

	return decIAm, nil
}

func (u *UnicastIAm) GetService() uint8 {
	return u.APDU.Service
}

func (u *UnicastIAm) GetType() uint8 {
	return u.APDU.Type
}
