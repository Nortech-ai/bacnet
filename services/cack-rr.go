package services

import (
	"fmt"
	"log"

	"github.com/Nortech-ai/bacnet/common"
	"github.com/Nortech-ai/bacnet/objects"
	"github.com/Nortech-ai/bacnet/plumbing"
)

// LogBufferCACK is a BACnet message.
type LogBufferCACK struct {
	*plumbing.BVLC
	*plumbing.NPDU
	*plumbing.APDU
}

type LogBufferCACKDec struct {
	ObjectType uint16
	InstanceId uint32
	PropertyId uint16
	FirstItem  bool
	LastItem   bool
	MoreItems  bool
	ItemCount  uint32
	Tags       []*objects.Object
}

type StatusFlags struct {
	InAlarm      bool
	Fault        bool
	Overridden   bool
	OutOfService bool
}

func (c *ComplexACK) DecodeRR() (LogBufferCACKDec, error) {
	decCACK := LogBufferCACKDec{}

	if len(c.APDU.Objects) < 3 {
		return decCACK, fmt.Errorf(
			"failed to decode CACK - objects count %d: %v",
			len(c.APDU.Objects),
			common.ErrWrongObjectCount,
		)
	}

	context := []uint8{8}
	objs := make([]*objects.Object, 0)
	for i, obj := range c.APDU.Objects {
		enc_obj, ok := obj.(*objects.Object)
		if !ok {
			return decCACK, fmt.Errorf(
				"LogBufferCACK object at index %d is not Object type: %s",
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
				return decCACK, fmt.Errorf(
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
					return decCACK, fmt.Errorf("decode Context object case 0: %v", err)
				}
				decCACK.ObjectType = objId.ObjectType
				decCACK.InstanceId = objId.InstanceNumber
			case combine(8, 1):
				value, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 1: %v", err)
				}
				propId := uint16(value)
				decCACK.PropertyId = propId
			case combine(8, 3):
				first, last, more, err := decResultsFlag(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 3: %v", err)
				}
				decCACK.FirstItem = first
				decCACK.LastItem = last
				decCACK.MoreItems = more
			case combine(8, 4):
				data, err := objects.DecUnsignedInteger(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 4: %v", err)
				}
				decCACK.ItemCount = data
			case combine(1, 2):
				value, err := objects.DecReal(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 2: %v", err)
				}
				objs = append(objs, &objects.Object{
					TagNumber: 2,
					TagClass:  true,
					Length:    uint8(obj.MarshalLen()),
					Value:     value,
				})
			case combine(5, 2):
				value, err := decStatusFlags(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 2: %v", err)
				}
				objs = append(objs, &objects.Object{
					TagNumber: 2,
					TagClass:  true,
					Length:    uint8(obj.MarshalLen()),
					Value:     value,
				})
			case combine(1, 0):
				value, err := objects.DecLogStatus(obj)
				if err != nil {
					return decCACK, fmt.Errorf("decode Context object case 0: %v", err)
				}
				objs = append(objs, &objects.Object{
					TagNumber: 0,
					TagClass:  true,
					Length:    uint8(obj.MarshalLen()),
					Value:     value,
				})
			default:
				log.Printf("Unknown Context object: context %v tag class %t tag number %d\n", context, enc_obj.TagClass, enc_obj.TagNumber)
			}
		} else {
			tag, err := decodeAppTags(enc_obj, &obj)
			if err != nil {
				return decCACK, fmt.Errorf("decode Application Tag: %v", err)
			}
			objs = append(objs, tag)
		}
	}
	decCACK.Tags = objs

	return decCACK, nil
}

func decResultsFlag(obj objects.APDUPayload) (bool, bool, bool, error) {
	var first, last, more bool
	enc_obj, ok := obj.(*objects.Object)
	if !ok {
		return false, false, false, common.ErrInvalidObjectType
	}
	first = enc_obj.Data[1]&0x80 == 0x80
	last = enc_obj.Data[1]&0x40 == 0x40
	more = enc_obj.Data[1]&0x20 == 0x20
	return first, last, more, nil
}

func decStatusFlags(obj objects.APDUPayload) (StatusFlags, error) {
	var status StatusFlags
	enc_obj, ok := obj.(*objects.Object)
	if !ok {
		return status, common.ErrInvalidObjectType
	}
	status.InAlarm = enc_obj.Data[1]&0x80 == 0x80
	status.Fault = enc_obj.Data[1]&0x40 == 0x40
	status.Overridden = enc_obj.Data[1]&0x20 == 0x20
	status.OutOfService = enc_obj.Data[1]&0x10 == 0x10
	return status, nil
}
