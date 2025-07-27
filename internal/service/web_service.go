package service

type WebService interface {
	// Web service methods will be defined here
}

type webService struct {
	// dependencies will be injected
}

func NewWebService(repository interface{}, db interface{}, validator interface{}, aiService interface{}) WebService {
	return &webService{}
}
