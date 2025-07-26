package initPhase

var (
	initPhase = make(chan struct{})

	errIsAlreadyFinished = "initPhase is already finished," +
		" this is allowed to use only before the main function is called"
)

func WaitToFinish() <-chan struct{} {
	return initPhase
}

func PanicIfFinished() {
	select {
	case <-initPhase:
		panic(errIsAlreadyFinished)
	default:
	}
}

func Finish() {
	close(initPhase)
}
