package objects

import "github.com/Nortech-ai/bacnet/common"

func ContextTag(tagNumber uint8, o *Object) *Object {
	o.TagClass = true
	o.TagNumber = tagNumber
	return o
}

func EncContextBool(tagNumber uint8, value bool) *Object {
	obj := &Object{
		TagClass:  true,
		TagNumber: tagNumber,
		Length:    1,
		Data:      []byte{byte(common.BoolToInt(value))},
	}
	return obj
}

func DecContextBool(rawPayload APDUPayload) (bool, error) {
	encObj, ok := rawPayload.(*Object)
	if !ok {
		return false, common.ErrInvalidObjectType
	}
	return common.IntToBool(int(encObj.Data[0] & 0x01)), nil
}
