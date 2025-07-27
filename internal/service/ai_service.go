package service

type AIService interface {
	// AI service methods will be defined here
}

type aiService struct {
	// dependencies will be injected
}

func NewAIService(repository interface{}, db interface{}, validator interface{}, outbound interface{}) AIService {
	return &aiService{}
}
