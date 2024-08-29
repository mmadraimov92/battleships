package cyclic

type Number struct {
	minValue     int8
	maxValue     int8
	currentValue int8
}

func NewNumber(min, max int8) *Number {
	return &Number{minValue: min, maxValue: max}
}

func (n *Number) Current() int8 {
	return n.currentValue
}

func (n *Number) Increment() {
	if n.currentValue+1 > n.maxValue {
		n.currentValue = n.minValue
	} else {
		n.currentValue++
	}
}

func (n *Number) Decrement() {
	if n.currentValue-1 < n.minValue {
		n.currentValue = n.maxValue
	} else {
		n.currentValue--
	}
}
