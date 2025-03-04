package internal

func LoopAction(action func(), times int) {
	for i := 0; i < times; i++ {
		action()
	}
}
