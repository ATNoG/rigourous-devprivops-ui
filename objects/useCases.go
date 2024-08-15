package objects

type UseCase struct {
	UseCase      string        `json:"use case" yaml:"use case"`
	IsMisuseCase bool          `json:"is misuse case" yaml:"is misuse case"`
	Requirements []Requirement `json:"requirements" yaml:"requirements"`
}

type Requirement struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Query       string `json:"query" yaml:"query"`
}

/*
type UseCase struct {
	UseCase      string        `yaml:"use case"`
	IsMisuseCase bool          `yaml:"is misuse case"`
	Requirements []Requirement `yaml:"requirements"`
}

type Requirement struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
	Query       string `yaml:"query"`
}
*/
