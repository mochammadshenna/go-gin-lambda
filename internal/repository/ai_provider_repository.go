package repository

type AIProviderRepository interface {
	// AI provider repository methods will be defined here
}

type aiProviderRepository struct {
	// dependencies will be injected
}

func NewAIProviderRepository(redis interface{}) AIProviderRepository {
	return &aiProviderRepository{}
}
