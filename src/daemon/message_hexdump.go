package daemon

import (
	"strconv"

	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/daemon/gnet"

	"reflect"

	"github.com/skycoin/skycoin/src/util"
)

var registered = false

// MessagesAnnotationsGenerator : Implementation of IAnnotationsGenerator for type gnet.Message
type MessagesAnnotationsGenerator struct {
	Message gnet.Message
}

// GenerateAnnotations : Implementation of IAnnotationsGenerator.GenerateAnnotations for type gnet.Message
func (mag *MessagesAnnotationsGenerator) GenerateAnnotations() []util.Annotation {
	var annotations = make([]util.Annotation, 2)
	annotations[0] = util.Annotation{Size: 4, Name: "Length"}
	annotations[1] = util.Annotation{Size: 4, Name: "Prefix"}
	var v = reflect.Indirect(reflect.ValueOf(mag.Message))

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		vF := v.Field(i)
		f := t.Field(i)
		if f.Tag.Get("enc") != "-" {
			if vF.CanSet() || f.Name != "_" {
				if v.Field(i).Kind() == reflect.Slice {
					annotations = append(annotations, util.Annotation{Size: 4, Name: f.Name + " length"})
					for j := 0; j < v.Field(i).Len(); j++ {
						annotations = append(annotations, util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "#" + strconv.Itoa(j)})
					}
				} else {
					annotations = append(annotations, util.Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name})
				}
			} else {
				//don't write anything
			}
		}
	}

	return annotations
}

// MessagesAnnotationsIterator : Implementation of IAnnotationsIterator for type gnet.Message
type MessagesAnnotationsIterator struct {
	Message      gnet.Message
	LengthCalled bool
	PrefixCalled bool
	CurrentField int
	MaxField     int
	CurrentIndex int
}

// NewMessagesAnnotationsIterator : Initializes struct MessagesAnnotationsIterator
func NewMessagesAnnotationsIterator(message gnet.Message) MessagesAnnotationsIterator {
	var mai = MessagesAnnotationsIterator{}
	mai.Message = message
	mai.LengthCalled = false
	mai.PrefixCalled = false
	mai.CurrentField = 0
	mai.CurrentIndex = -1

	mai.MaxField = reflect.Indirect(reflect.ValueOf(mai.Message)).NumField()

	return mai
}

// Next : Yields next element of MessagesAnnotationsIterator
func (mai *MessagesAnnotationsIterator) Next() (util.Annotation, bool) {
	if !mai.LengthCalled {
		mai.LengthCalled = true
		return util.Annotation{Size: 4, Name: "Length"}, true

	}
	if !mai.PrefixCalled {
		mai.PrefixCalled = true
		return util.Annotation{Size: 4, Name: "Prefix"}, true

	}
	if mai.CurrentField == mai.MaxField {
		return util.Annotation{}, false
	}

	var i = mai.CurrentField
	var j = mai.CurrentIndex

	var v = reflect.Indirect(reflect.ValueOf(mai.Message))
	t := v.Type()
	vF := v.Field(i)
	f := t.Field(i)
	if f.Tag.Get("enc") != "-" {
		if vF.CanSet() || f.Name != "_" {
			if v.Field(i).Kind() == reflect.Slice {
				if mai.CurrentIndex == -1 {
					mai.CurrentIndex = 0
					return util.Annotation{Size: 4, Name: f.Name + " length"}, true
				}
				mai.CurrentIndex++
				if mai.CurrentIndex != v.Field(i).Len() {

					//mai.CurrentField++
					return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "#" + strconv.Itoa(j)}, true
				}
				mai.CurrentIndex = -1
				mai.CurrentField++
				return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Slice(j, j+1).Interface())[4:]), Name: f.Name + "#" + strconv.Itoa(j)}, true

			}

			mai.CurrentField++
			return util.Annotation{Size: len(encoder.Serialize(v.Field(i).Interface())), Name: f.Name}, true

		}
	}

	return util.Annotation{}, false
}
