package main

import (
	"fmt"
)

type datatypes interface {
	Value() interface{}
	String() string
}

type Datatype struct {
	Data interface{}
}

func (d Datatype) Value() interface{} {
	return d.Data
}

func (d Datatype) String() string {
	v := d.Value()
	s := fmt.Sprintf("%#v", v)
	return s
}

type Tag struct {
	Datatype
	Kind string
}

func (t Tag) Value() interface{} {
	return t.Kind + ":" + t.Data.(string)
}

func NewTag(kind string, msg string) *Tag {
	t := new(Tag)
	t.Kind = kind
	t.Data = msg
	return t
}

func NewErr(msg string) *Tag {
	return NewTag("err", msg)
}

func NewTrue(msg string) *Tag {
	return NewTag("true", msg)
}

func NewNil(msg string) *Tag {
	return NewTag("nil", msg)
}

func (t Tag) String() string {
	return t.Kind + "(" + t.Data.(string) + ")"
}

type Int struct {
	Datatype
}

func NewInt(v int) *Int {
	i := new(Int)
	i.Data = v
	return i
}

func (i Int) Value() interface{} {
	return i.Data
}

type Vector struct {
	Datatype
}

func NewVector(items []datatypes) Vector {
	v := *new(Vector)
	v.Data = items
	return v
}

func (v *Vector) App(i datatypes) Vector {
	return NewVector(append(v.Data.([]datatypes), i))
}

func (v *Vector) Ato(n int) datatypes {
	return v.Data.([]datatypes)[n]
}

func (v Vector) String() string {
	str := "("
	for _, i := range v.Data.([]datatypes) {
		str += i.String() + " "
	}
	if len(str) > 1 {
		str = str[:len(str)-1]
	}
	return str + ")"
}

type CharVector struct {
	Vector
}

func NewCharVector(str string) *CharVector {
	cv := new(CharVector)
	if str[0] == "'"[0] && str[len(str)-1] == "'"[0] {
		cv.Data = str[1 : len(str)-1]
	} else {
		cv.Data = str
	}
	return cv
}

func (cv *CharVector) Cat(cv2 *CharVector) *CharVector {
	return NewCharVector(cv.Value().(string) + cv2.Value().(string))
}

func (cv CharVector) String() string {
	return cv.Value().(string)
}
