package main

func assert4true(cond bool) {
	if !cond {
		panic(cond)
	}
}

func assert4false(cond bool) {
	if cond {
		panic(cond)
	}
}

func check4true(cond bool) bool {
	return cond
}

func check4false(cond bool) bool {
	return !cond
}
