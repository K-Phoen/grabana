package sandbox

type OperatorState struct {
	// lastEvaluation is the ResourceVersion last evaluated
LastEvaluation string `json:"lastEvaluation"`
	// state describes the state of the lastEvaluation.
// It is limited to three possible states for machine evaluation.
State SandboxState `json:"state"`
	// descriptiveState is an optional more descriptive state field which has no requirements on format
DescriptiveState *string `json:"descriptiveState,omitempty"`
	// details contains any extra information that is operator-specific
Details any `json:"details,omitempty"`
}

type Sandbox struct {
	Name string `json:"name"`
	NestedStruct struct {
	Foo string `json:"foo"`
} `json:"nestedStruct"`
	AnythingPlz any `json:"anythingPlz"`
	SomeMap map[string]int32 `json:"someMap,omitempty"`
	// operatorStates is a map of operator ID to operator state evaluations.
// Any operator which consumes this kind SHOULD add its state evaluation information to this field.
OperatorStates map[string]OperatorState `json:"operatorStates,omitempty"`
}

// state describes the state of the lastEvaluation.
// It is limited to three possible states for machine evaluation.
type SandboxState string
const (
	Success SandboxState = "success"
	In_progress SandboxState = "in_progress"
	Failed SandboxState = "failed"
)

