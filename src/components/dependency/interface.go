package dependency

type Interface[T any] interface {
	Dependency() T
}
