package service

type GenerationService interface {
	// Generation service methods will be defined here
}

type generationService struct {
	// dependencies will be injected
}

func NewGenerationService(repository interface{}, db interface{}, validator interface{}, aiService interface{}) GenerationService {
	return &generationService{}
}
