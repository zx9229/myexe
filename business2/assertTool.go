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
