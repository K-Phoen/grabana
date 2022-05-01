package stackdriver

type FilterOption func(filter *filter)

type filter struct {
	operator string

	leftOperand  string
	rightOperand string
}

func Eq(leftOperand string, rightOperand string) FilterOption {
	return func(filter *filter) {
		filter.operator = "="
		filter.leftOperand = leftOperand
		filter.rightOperand = rightOperand
	}
}

func Neq(leftOperand string, rightOperand string) FilterOption {
	return func(filter *filter) {
		filter.operator = "!="
		filter.leftOperand = leftOperand
		filter.rightOperand = rightOperand
	}
}

func Matches(leftOperand string, rightOperand string) FilterOption {
	return func(filter *filter) {
		filter.operator = "=~"
		filter.leftOperand = leftOperand
		filter.rightOperand = rightOperand
	}
}

func NotMatches(leftOperand string, rightOperand string) FilterOption {
	return func(filter *filter) {
		filter.operator = "!=~"
		filter.leftOperand = leftOperand
		filter.rightOperand = rightOperand
	}
}
