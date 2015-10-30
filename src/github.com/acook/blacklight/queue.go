package main

type Queue struct {
	Items chan datatypes
}

func NewQueue() *Queue {
	q := &Queue{}
	q.Items = make(chan datatypes, 16)

	return q
}

func (q *Queue) Enqueue(item datatypes) {
	q.Items <- item
}

func (q *Queue) Dequeue() datatypes {
	return <-q.Items
}

func (q Queue) Value() interface{} {
	return q.Items
}

func (q Queue) String() string {
	str := "{"

StringLoop:
	for {
		select {
		case i := <-q.Items:
			str += i.String()
			str += " "
		default:
			break StringLoop
		}
	}
	return str[:len(str)-1] + "}"
}
