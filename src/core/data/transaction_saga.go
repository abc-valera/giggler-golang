package data

// TODO: implement saga pattern

type SagaCommit struct {
	commitFunc   func() error
	rollbackFunc func() error
}

func CommitWithSaga(sagaCommits ...SagaCommit) {
	panic("Not implemented")
}
