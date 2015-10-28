package main

import (
	"fmt"
	"strconv"
)

type operation interface {
	Eval(stack) stack
	Value() datatypes
	String() string
}

type Op struct {
	Name string
	Data datatypes
}

func (o Op) Eval(current stack) stack {
	switch o.Name {
	// @stack (current Stack)
	case "decap":
		current.Decap()
	case "depth":
		current.Push(NewInt(current.Depth()))
	case "drop":
		current.Drop()
	case "dup":
		current.Dup()
	case "over":
		current.Over()
	case "purge":
		current.Purge()
	case "rot":
		current.Rot()
	case "swap":
		current.Swap()

	// NativeIntegers (Int)
	case "add":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		sum := n2 + n1
		current.Push(NewInt(sum))
	case "sub":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 - n1
		current.Push(NewInt(result))
	case "mul":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		product := n2 * n1
		current.Push(NewInt(product))
	case "div":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 / n1
		current.Push(NewInt(result))
	case "mod":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		remainder := n2 % n1
		current.Push(NewInt(remainder))
	case "n-to-s":
		i := current.Pop()
		n := i.Value().(int)
		str := strconv.Itoa(n)
		current.Push(NewCharVector(str))

	// Debug
	case "print":
		i := current.Pop()
		switch i.(type) {
		case *Int:
			v := i.(*Int).Value().(int)
			fmt.Printf("%v", v)
		case *CharVector:
			v := i.(*CharVector).Value().(string)
			print(v)
		case *MetaStack:
			v := i.(*MetaStack).String()
			print(v)
		case *SystemStack:
			v := i.(*SystemStack).String()
			print(v)
		default:
			fmt.Printf("%#v", i)
		}
		print("\n")

	// Vectors
	case "cat":
		i1 := current.Pop().(*CharVector)
		i2 := current.Pop().(*CharVector)

		result := i2.Cat(i1)
		current.Push(result)

	// Queues
	case "newq":
		q := NewQueue()
		current.Push(q)
	case "enq":
		i := current.Pop()
		q := (*current.Peek()).(*Queue)
		q.Enqueue(i)
	case "deq":
		q := (*current.Peek()).(*Queue)
		i := q.Dequeue()
		current.Push(i)
	case "proq":
		current = processQueue(current)

	default:
		warn("UNIMPLEMENTED operation: " + o.String())
	}
	return current
}

func (o Op) Value() datatypes {
	return o.Data
}

func (o Op) String() string {
	return o.Name
}

func newOp(t string) *Op {
	op := new(Op)
	op.Name = t
	return op
}

type metaOp struct {
	Op
}

func (m metaOp) Eval(meta stack) stack {
	switch m.Name {
	case "@":
		s := *meta.Peek()
		current := s.(*SystemStack)
		current.Push(current)
	case "^":
		s := *meta.Peek()
		current := s.(*SystemStack)
		meta.Swap()
		s = *meta.Peek()
		prev := s.(*SystemStack)
		meta.Swap()
		current.Push(prev)
	case "$decap":
		meta.Decap()
	case "$drop":
		meta.Drop()
	case "$new":
		if meta.Depth() > 0 {
			s := *meta.Peek()
			os := s.(*SystemStack)
			ns := NewSystemStack()
			ns.Push(os)
			meta.Push(ns)
		} else {
			meta.Push(NewSystemStack())
			meta = newMetaOp("$new").Eval(meta)
		}
	case "$swap":
		meta.Swap()
	default:
		warn("UNIMPLEMENTED $operation: " + m.String())
	}

	return meta
}

func newMetaOp(t string) *metaOp {
	op := new(metaOp)
	op.Name = t
	return op
}

type pushLiteral struct {
	Op
}

func (o pushLiteral) Eval(s stack) stack {
	s.Push(o.Value())
	return s
}

type pushInteger struct {
	pushLiteral
}

func newPushInteger(t string) *pushInteger {
	pi := new(pushInteger)
	pi.Name = t
	i, _ := strconv.Atoi(t)
	pi.Data = NewInt(i)
	return pi
}

type pushWord struct {
	pushLiteral
}

func newPushWord(t string) *pushWord {
	pw := new(pushWord)
	pw.Name = t
	w := NewWord(t)
	pw.Data = w
	return pw
}

type pushWordVector struct {
	pushLiteral
	Contents []operation
}

func newPushWordVector(t string) *pushWordVector {
	pwv := new(pushWordVector)
	pwv.Name = t
	return pwv
}

func (pwv *pushWordVector) Eval(s stack) stack {
	wv := NewWordVector(pwv.Contents)
	s.Push(wv)
	return s
}

type pushCharVector struct {
	pushLiteral
}

func newPushCharVector(t string) *pushCharVector {
	ps := new(pushCharVector)
	ps.Name = t
	ps.Data = NewCharVector(t)
	return ps
}

type pushChar struct {
	pushLiteral
}

func newPushChar(t string) *pushChar {
	pc := new(pushChar)
	pc.Name = t
	return pc
}

type pushQueue struct {
	pushLiteral
	Contents []operation
}

func processQueue(s stack) stack {
	wv := s.Pop().(WordVector)
	q := s.Pop().(*Queue)
	var tokens []string

	for _, w := range wv.Data {
		tokens = append(tokens, w.Name)
	}

ProcLoop:
	for {
		select {
		case item := <-q.Items:
			s.Push(item)
			meta := NewMetaStack() // FIXME: this should be the actual $meta stack
			meta.Push(s)
			ops := lex(tokens)
			doEval(meta, ops)
		default:
			break ProcLoop
		}
	}

	return s
}
