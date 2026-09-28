package completer

type completerState struct {
	lastPrefix []rune
	tabPressed bool
}

type completerAction int

const (
	actBell completerAction = iota
	actInsertOne
	actInsertLCP
	actList
)

func decide(numOfMatches int, hasLCP bool, state completerState) (action completerAction, newState completerState) {
	switch {
	case numOfMatches == 0:
		return actBell, state

	case numOfMatches == 1:
		state.tabPressed = false
		return actInsertOne, state

	case numOfMatches > 1 && !hasLCP:
		if state.tabPressed {
			state.tabPressed = false
			return actList, state
		} else {
			state.tabPressed = true
			return actBell, state
		}

	case hasLCP:
		return actInsertLCP, state

	default:
		return actBell, state
	}
}
