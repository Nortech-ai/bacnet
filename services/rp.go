package services

import (
	"fmt"

	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/objects"
	"github.com/Nortech-ai/bacnet/plumbing"
)

// UnconfirmedReadProperty is a BACnet message.
type ConfirmedReadProperty struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedReadPropertyDec struct {
	ObjectType  uint16
	InstanceNum uint32
	PropertyId  uint16
}

func ConfirmedReadPropertyObjects(objectType uint16, instN uint32, propId uint16) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 2)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.ContextTag(1, objects.EncUnsignedInteger(uint(propId)))

	return objs
}

func ConfirmedReadPropertyMultipleObjects(objectType uint16, instN uint32, propIds []uint16) []objects.APDUPayload {
	length := 3 + (1 * len(propIds))
	objs := make([]objects.APDUPayload, length)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.EncOpeningTag(1)
	for i, p := range propIds {
		objs[i+2] = objects.ContextTag(
			0, objects.EncUnsignedInteger(uint(p)))
	}
	objs[len(objs)-1] = objects.EncClosingTag(1)

	return objs
}

func NewConfirmedReadProperty(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedReadProperty {
	c := &ConfirmedReadProperty{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadProperty, ConfirmedReadPropertyObjects(
			objects.ObjectTypeAnalogOutput, 1, objects.PropertyIdPresentValue)),
	}
	c.SetLength()

	return c
}

func NewConfirmedReadPropertyMultiple(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedReadProperty {
	c := &ConfirmedReadProperty{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedReadProperty, ConfirmedReadPropertyMultipleObjects(
			objects.ObjectTypeAnalogOutput, 1, []uint16{})),
	}
	c.SetLength()

	return c
}

func (c *ConfirmedReadProperty) UnmarshalBinary(b []byte) error {
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

func (c *ConfirmedReadProperty) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, fmt.Errorf("failed to marshal binary: %v", err)
	}
	return b, nil
}

func (c *ConfirmedReadProperty) MarshalTo(b []byte) error {
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

func (c *ConfirmedReadProperty) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedReadProperty) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedReadProperty) Decode() (ConfirmedReadPropertyDec, error) {
	decCRP := ConfirmedReadPropertyDec{}

	if len(c.APDU.Objects) != 2 {
		return decCRP, fmt.Errorf(
			"failed to decode ConfirmedRP - object count %d: %v",
			len(c.APDU.Objects),
			common.ErrWrongObjectCount,
		)
	}

	context := []uint8{8}
	for i, obj := range c.APDU.Objects {
		enc_obj, ok := obj.(*objects.Object)
		if !ok {
			return decCRP, fmt.Errorf(
				"ComplexACK object at index %d is not Object type: %v",
				i, common.ErrInvalidObjectType,
			)
		}

		// add or remove context based on opening and closing tags
		if enc_obj.Length == 6 {
			context = append(context, enc_obj.TagNumber)
			continue
		}
		if enc_obj.Length == 7 {
			if len(context) == 0 {
				return decCRP, fmt.Errorf(
					"LogBufferCACK object at index %d has mismatched closing tag: %v",
					i, common.ErrInvalidObjectType,
				)
			}
			context = context[:len(context)-1]
			continue
		}
		if enc_obj.TagClass {
			c := combine(context[len(context)-1], enc_obj.TagNumber)
			switch c {
			case combine(8, 0):
				objId, err := objects.DecObjectIdentifier(obj)
				if err != nil {
					return decCRP, fmt.Errorf("decoding ConfirmedRP: %v", err)
				}
				decCRP.ObjectType = objId.ObjectType
				decCRP.InstanceNum = objId.InstanceNumber
			case combine(8, 2):
				value, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCRP, fmt.Errorf("decoding ConfirmedRP: %v", err)
				}
				propId := uint16(value)
				decCRP.PropertyId = propId
			}
		}
	}
	return decCRP, nil
}

func (u *ConfirmedReadProperty) GetService() uint8 {
	return u.APDU.Service
}

func (u *ConfirmedReadProperty) GetType() uint8 {
	return u.APDU.Type
}
