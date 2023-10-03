package dto

func NewDomainModels[DomainModel any, DtoModel Interface[DomainModel]](
	dtos []DtoModel,
) []DomainModel {
	domainModels := make([]DomainModel, len(dtos), 0)

	for _, dto := range dtos {
		domainModels = append(domainModels, dto.ToDomain())
	}

	return domainModels
}

func NewDTOs[DomainModel any, DtoModel Interface[DomainModel]](
	domainModels []DomainModel,
	newDto func(DomainModel) DtoModel,
) []DtoModel {
	dtos := make([]DtoModel, len(domainModels), 0)

	for _, domainModel := range domainModels {
		dtos = append(dtos, newDto(domainModel))
	}

	return dtos
}
