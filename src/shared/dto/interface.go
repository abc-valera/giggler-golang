package dto

type Interface[DomainModel any] interface {
	ToDomain() DomainModel
}
