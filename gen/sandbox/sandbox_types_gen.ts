// state describes the state of the lastEvaluation.
// It is limited to three possible states for machine evaluation.
export enum SandboxState {
	Success = "success",
	In_progress = "in_progress",
	Failed = "failed",
}

export interface OperatorState {
	// lastEvaluation is the ResourceVersion last evaluated
	lastEvaluation: string;
	// state describes the state of the lastEvaluation.
	// It is limited to three possible states for machine evaluation.
	state: SandboxState;
	// descriptiveState is an optional more descriptive state field which has no requirements on format
	descriptiveState?: string;
	// details contains any extra information that is operator-specific
	details?: any;
}

export interface Sandbox {
	name: string;
	nestedStruct: struct;
	anythingPlz: any;
	someMap?: map;
	// operatorStates is a map of operator ID to operator state evaluations.
	// Any operator which consumes this kind SHOULD add its state evaluation information to this field.
	operatorStates?: map;
}

