// Package for the internal objects the website manipulates
package objects

// Represents a use case
// A use case is a list of requirements that have to be met in order for it to be possible
type UseCase struct {
	UseCase      string        `json:"use case" yaml:"use case"`             // The use case title
	IsMisuseCase bool          `json:"is misuse case" yaml:"is misuse case"` // Whether it is a misuse case
	Requirements []Requirement `json:"requirements" yaml:"requirements"`     // The list of requirements
}

// Represents a requirement
// A requirement has a query that finds proof that the system meets it
type Requirement struct {
	Title       string `json:"title" yaml:"title"`             // The requirement's title
	Description string `json:"description" yaml:"description"` // The requirement's description
	Query       string `json:"query" yaml:"query"`             // The requirement's query
}
