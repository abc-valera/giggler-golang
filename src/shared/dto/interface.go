package dto

// Interface marks a type as a DTO, and allows it to be casted to a domain model.
type Interface[DomainModel any] interface {
	ToDomain() DomainModel
}
