package sandbox

type Sandbox interface {
	IsAvailable() bool
	IsAuthenticated() bool
	Authenticate() error
	ExecuteClaude(prompt string, model string) error
}

func New(sandboxType string) (Sandbox, error) {
	if sandboxType == "docker" {
		return NewDocker(), nil
	}
	return nil, nil
}
