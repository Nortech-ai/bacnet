package bacnet

import (
	"github.com/Nortech-ai/bacnet/plumbing"
	"github.com/Nortech-ai/bacnet/services"
)

const (
	DEFAULT_ACCEPTED_SIZE        = 1024
	DEFAULT_SEGMENTATION_SUPPORT = 0x3 // No segmentation
)

func NewWhois() ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncBroadcast)
	npdu := plumbing.NewNPDU(false, false, false, false)
	u := services.NewUnconfirmedWhoIs(bvlc, npdu)
	return u.MarshalBinary()
}

func NewIAm(instN uint32, vendorId uint16) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncBroadcast)

	npdu := plumbing.NewNPDU(false, true, false, false)
	npdu.DNET = 0xFFFF
	npdu.DLEN = 0
	npdu.Hop = 0xFF

	u := services.NewUnconfirmedIAm(bvlc, npdu)

	u.APDU.Objects = services.IAmObjects(instN,
		DEFAULT_ACCEPTED_SIZE, DEFAULT_SEGMENTATION_SUPPORT, vendorId)
	u.SetLength()

	return u.MarshalBinary()
}

func NewCACK(service uint8, objectType uint16, instN uint32, propertyId uint16, value float32) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, false)

	c := services.NewComplexACK(bvlc, npdu)

	c.APDU.Service = service
	c.APDU.InvokeID = 1
	c.APDU.Objects = services.ComplexACKObjects(objectType, instN, propertyId, value)

	c.SetLength()

	return c.MarshalBinary()
}

func NewSACK(service uint8) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, false)

	s := services.NewSimpleACK(bvlc, npdu)

	s.APDU.Service = service
	s.APDU.InvokeID = 1

	s.SetLength()

	return s.MarshalBinary()
}

func NewError(service, errorClass, errorCode uint8) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, false)

	e := services.NewError(bvlc, npdu)

	e.APDU.Service = service
	e.APDU.InvokeID = 1
	e.APDU.Objects = services.ErrorObjects(errorClass, errorCode)

	e.SetLength()

	return e.MarshalBinary()
}

func NewReadProperty(objectType uint16, instanceNumber uint32, propertyId uint16) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, true)

	c := services.NewConfirmedReadProperty(bvlc, npdu)

	c.APDU.Service = services.ServiceConfirmedReadProperty
	c.APDU.MaxSize = 5
	c.APDU.InvokeID = 1
	c.APDU.Objects = services.ConfirmedReadPropertyObjects(objectType, instanceNumber, propertyId)

	c.SetLength()

	return c.MarshalBinary()
}

func NewReadPropertyMultiple(objectType uint16, instanceNumber uint32, propertyIds []uint16) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, true)

	c := services.NewConfirmedReadProperty(bvlc, npdu)

	c.APDU.Service = services.ServiceConfirmedReadPropMultiple
	c.APDU.MaxSeg = 7
	c.APDU.MaxSize = 5
	c.APDU.InvokeID = 1
	c.APDU.Flags = 2
	c.APDU.Objects = services.ConfirmedReadPropertyMultipleObjects(objectType, instanceNumber, propertyIds)
	c.SetLength()

	return c.MarshalBinary()
}

func NewReadRange(objectType uint16, instanceNumber uint32, propertyId uint16, rangeStart uint16, length int32) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, true)

	c, _ := services.NewConfirmedReadRange(bvlc, npdu)

	c.APDU.Service = services.ServiceConfirmedReadRange
	c.APDU.MaxSeg = 7
	c.APDU.MaxSize = 5
	c.APDU.InvokeID = 1
	c.APDU.Flags = 2
	c.APDU.Objects = services.ConfirmedReadRangeObjects(objectType, instanceNumber, propertyId, rangeStart, length)

	c.SetLength()

	return c.MarshalBinary()
}

func NewSubscribeCOV(objectType uint16, instanceNumber uint32, processId uint, lifetime uint, expect, cancel bool) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, true)
	c, _ := services.NewConfirmedSubscribeCOV(bvlc, npdu)

	c.APDU.Service = services.ServiceConfirmedSubscribeCOV
	c.APDU.MaxSeg = 7
	c.APDU.MaxSize = 5
	c.APDU.InvokeID = 1
	c.APDU.Flags = 2
	if cancel {
		c.APDU.Objects = services.CancelCOVOBjects(
			processId, objectType, instanceNumber)
	} else {
		c.APDU.Objects = services.COVObjects(
			processId, objectType, instanceNumber, expect, lifetime)
	}
	c.SetLength()

	return c.MarshalBinary()
}

func NewWriteProperty(objectType uint16, instanceNumber uint32, propertyId uint16, data interface{}) ([]byte, error) {
	bvlc := plumbing.NewBVLC(plumbing.BVLCFuncUnicast)
	npdu := plumbing.NewNPDU(false, false, false, true)

	c := services.NewConfirmedWriteProperty(bvlc, npdu)

	c.APDU.Service = services.ServiceConfirmedWriteProperty
	c.APDU.MaxSize = 5
	c.APDU.InvokeID = 1
	c.APDU.Objects = services.ConfirmedWritePropertyObjects(objectType, instanceNumber, propertyId, data)

	c.SetLength()

	return c.MarshalBinary()
}
