package dataModel

type SagaCommit struct {
	commitFunc   func() error
	rollbackFunc func() error
}

func CommitWithSaga(sagaCommits ...SagaCommit) {
	// TODO: implement
	panic("Not implemented")
}
