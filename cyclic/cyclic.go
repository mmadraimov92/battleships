package cyclic

type Number struct {
	maxValue     int8
	currentValue int8
}

func NewNumber(max int8) *Number {
	return &Number{maxValue: max}
}

func (n *Number) Current() int8 {
	return n.currentValue
}

func (n *Number) Reset() {
	n.currentValue = 0
}

func (n *Number) Set(value int8) {
	n.currentValue = value
}

func (n *Number) Increment() {
	if n.currentValue+1 > n.maxValue {
		n.currentValue = 0
	} else {
		n.currentValue++
	}
}

func (n *Number) Decrement() {
	if n.currentValue-1 < 0 {
		n.currentValue = n.maxValue
	} else {
		n.currentValue--
	}
}
