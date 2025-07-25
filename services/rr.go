package services

import (
	"fmt"

	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/objects"
	"github.com/Nortech-ai/bacnet/plumbing"
)

// UnconfirmedReadRange is a BACnet message.
type ConfirmedReadRange struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedReadRangeDec struct {
	ObjectType  uint16
	InstanceNum uint32
	PropertyId  uint16
	Tags        []*objects.Object
}

func ConfirmedReadRangeObjects(objectType uint16, instN uint32, property uint16, index uint16, count int32) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 6)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	switch property {
	case objects.PropertyIdPresentValue:
		objs[1] = objects.ContextTag(1, objects.EncUnsignedInteger(uint(property)))
	case objects.PropertyIdLogBuffer:
		objs[1] = objects.ContextTag(1, objects.EncUnsignedInteger(uint(property)))
		objs[2] = objects.EncOpeningTag(3)
		objs[3] = objects.EncUnsignedInteger(uint(index))
		objs[4] = objects.EncSignedInteger(int(count))
		objs[5] = objects.EncClosingTag(3)
	default:
		panic("Not Implemented")
	}

	return objs
}

func NewConfirmedReadRange(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) (*ConfirmedReadRange, uint8) {
	c := &ConfirmedReadRange{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadRange, ConfirmedReadRangeObjects(
			0, 0, 131, 0, 0)),
	}
	c.SetLength()

	return c, c.APDU.Type
}

func (c *ConfirmedReadRange) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedReadRange) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedReadRange) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return fmt.Errorf(
			"failed to unmarshal ConfirmedRP - marshal length %d binary length %d: %v",
			c.MarshalLen(), l,
			common.ErrTooShortToParse,
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedRP %+v: %v",
			c, common.ErrTooShortToParse,
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedRP %+v: %v",
			c, common.ErrTooShortToParse,
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedRP %+v: %v",
			c, common.ErrTooShortToParse,
		)
	}

	return nil
}

func (c *ConfirmedReadRange) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, fmt.Errorf("failed to marshal binary: %v", err)
	}
	return b, nil
}

func (c *ConfirmedReadRange) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return fmt.Errorf(
			"failed to marshal ConfirmedRP - marshal length %d binary length %d: %v",
			c.MarshalLen(), len(b),
			common.ErrTooShortToMarshalBinary,
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedRP: %v", err)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedRP: %v", err)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedRP: %v", err)
	}

	return nil
}

func (c *ConfirmedReadRange) Decode() (ConfirmedReadRangeDec, error) {
	decCRP := ConfirmedReadRangeDec{}

	if len(c.APDU.Objects) < 2 {
		return decCRP, fmt.Errorf(
			"failed to decode ConfirmedRP - object count %d: %v",
			len(c.APDU.Objects),
			common.ErrWrongObjectCount,
		)
	}

	objs := make([]*objects.Object, 0)
	for i, obj := range c.APDU.Objects {
		enc_obj, ok := obj.(*objects.Object)
		if !ok {
			return decCRP, fmt.Errorf(
				"ComplexACK object at index %d is not Object type: %v",
				i, common.ErrInvalidObjectType,
			)
		}

		if enc_obj.TagClass {
			switch enc_obj.TagNumber {
			case 0:
				objId, err := objects.DecObjectIdentifier(obj)
				if err != nil {
					return decCRP, fmt.Errorf("decode Context object case 0: %v", err)
				}
				decCRP.ObjectType = objId.ObjectType
				decCRP.InstanceNum = objId.InstanceNumber
			case 1:
				value, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCRP, fmt.Errorf("decode Context object case 1: %v", err)
				}
				propId := uint16(value)
				decCRP.PropertyId = propId
			}
		} else {
			tag, err := decodeAppTags(enc_obj, &obj)
			if err != nil {
				return decCRP, fmt.Errorf("decode Application Tag: %v", err)
			}
			objs = append(objs, tag)
		}
		decCRP.Tags = objs
	}
	return decCRP, nil
}
