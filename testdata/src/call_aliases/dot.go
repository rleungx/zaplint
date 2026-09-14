package call_aliases

import . "go.uber.org/zap"

func dot() {
	String("BadKey", "v") // want "key 'BadKey' should be in snake_case"
	Any("value", true)    // want "replace zap.Any with zap.Bool"
}

func shadowed() {
	Bool := func(string, bool) int { return 0 }
	_ = Bool
	Any("value", true) // Replacing Any by Bool would resolve to the local function.
}
