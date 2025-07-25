package services

import (
	"fmt"

	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/objects"
	"github.com/Nortech-ai/bacnet/plumbing"
)

// UnconfirmedIAm is a BACnet message.
type ConfirmedWriteProperty struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type ConfirmedWritePropertyDec struct {
	ObjectType  uint16
	InstanceNum uint32
	PropertyId  uint16
	Value       float32
	Priority    uint8
	Tags        []*objects.Object
}

func ConfirmedWritePropertyObjects(objectType uint16, instN uint32, propertyId uint16, data interface{}) []objects.APDUPayload {
	objs := make([]objects.APDUPayload, 6)

	objs[0] = objects.EncObjectIdentifier(true, 0, objectType, instN)
	objs[1] = objects.ContextTag(1, objects.EncUnsignedInteger(uint(propertyId)))
	objs[2] = objects.EncOpeningTag(3)
	switch data := data.(type) {
	case float32:
		objs[3] = objects.EncReal(data)
	case uint:
		objs[3] = objects.EncUnsignedInteger(data)
	case string:
		objs[3] = objects.EncString(data)
	}
	objs[4] = objects.EncClosingTag(3)
	objs[5] = objects.ContextTag(4, objects.EncUnsignedInteger(16))

	return objs
}

func NewConfirmedWriteProperty(bvlc *plumbing.BVLC, npdu *plumbing.NPDU) *ConfirmedWriteProperty {
	c := &ConfirmedWriteProperty{
		BVLC: bvlc,
		NPDU: npdu,
		APDU: plumbing.NewAPDU(plumbing.ConfirmedReq, ServiceConfirmedWriteProperty, nil),
	}
	c.SetLength()

	return c
}

func (c *ConfirmedWriteProperty) UnmarshalBinary(b []byte) error {
	if l := len(b); l < c.MarshalLen() {
		return fmt.Errorf(
			"failed to unmarshal ConfirmedWP - marshal length %d binary length %d: %v",
			c.MarshalLen(), l,
			common.ErrTooShortToParse,
		)
	}

	var offset int = 0
	if err := c.BVLC.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedWP %+v: %v", c, common.ErrTooShortToParse,
		)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedWP %+v: %v", c, common.ErrTooShortToParse,
		)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.UnmarshalBinary(b[offset:]); err != nil {
		return fmt.Errorf(
			"unmarshalling ConfirmedWP %+v: %v", c, common.ErrTooShortToParse,
		)
	}

	return nil
}

func (c *ConfirmedWriteProperty) MarshalBinary() ([]byte, error) {
	b := make([]byte, c.MarshalLen())
	if err := c.MarshalTo(b); err != nil {
		return nil, fmt.Errorf("failed to marshal binary: %v", err)
	}
	return b, nil
}

func (c *ConfirmedWriteProperty) MarshalTo(b []byte) error {
	if len(b) < c.MarshalLen() {
		return fmt.Errorf(
			"failed to marshal ConfirmedWP - marshal length %d binary length %d: %v",
			c.MarshalLen(), len(b),
			common.ErrTooShortToMarshalBinary,
		)
	}
	var offset = 0
	if err := c.BVLC.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedWP: %v", err)
	}
	offset += c.BVLC.MarshalLen()

	if err := c.NPDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedWP: %v", err)
	}
	offset += c.NPDU.MarshalLen()

	if err := c.APDU.MarshalTo(b[offset:]); err != nil {
		return fmt.Errorf("failed to marshal ConfirmedWP: %v", err)
	}

	return nil
}

func (c *ConfirmedWriteProperty) MarshalLen() int {
	l := c.BVLC.MarshalLen()
	l += c.NPDU.MarshalLen()
	l += c.APDU.MarshalLen()

	return l
}

func (c *ConfirmedWriteProperty) SetLength() {
	c.BVLC.Length = uint16(c.MarshalLen())
}

func (c *ConfirmedWriteProperty) Decode() (ConfirmedWritePropertyDec, error) {
	decCWP := ConfirmedWritePropertyDec{}

	if len(c.APDU.Objects) != 5 {
		return decCWP, fmt.Errorf(
			"failed to decode ConfirmedWP - object count %d: %v",
			len(c.APDU.Objects),
			common.ErrWrongObjectCount,
		)
	}

	context := []uint8{8}
	objs := make([]*objects.Object, 0)
	for i, obj := range c.APDU.Objects {
		enc_obj, ok := obj.(*objects.Object)
		if !ok {
			return decCWP, fmt.Errorf(
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
				return decCWP, fmt.Errorf(
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
					return decCWP, fmt.Errorf("decoding ConfirmedWP: %v", err)
				}
				decCWP.ObjectType = objId.ObjectType
				decCWP.InstanceNum = objId.InstanceNumber
			case combine(8, 1):
				propId, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCWP, fmt.Errorf("decoding ConfirmedWP: %v", err)
				}
				decCWP.PropertyId = uint16(propId)
			case combine(8, 4):
				priority, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCWP, fmt.Errorf("decoding ConfirmedWP: %v", err)
				}
				decCWP.Priority = uint8(priority)
			}
		} else {
			tag, err := decodeAppTags(enc_obj, &obj)
			if err != nil {
				return decCWP, fmt.Errorf("decode Application Tag: %v", err)
			}
			objs = append(objs, tag)
		}
	}
	decCWP.Tags = objs

	return decCWP, nil
}

func (u *ConfirmedWriteProperty) GetService() uint8 {
	return u.APDU.Service
}

func (u *ConfirmedWriteProperty) GetType() uint8 {
	return u.APDU.Type
}
