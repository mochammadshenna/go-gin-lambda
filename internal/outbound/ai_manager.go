package outbound

type AIManager interface {
	// AI manager methods will be defined here
}

type aiManager struct {
	// dependencies will be injected
}

func NewAIManager() AIManager {
	return &aiManager{}
}
