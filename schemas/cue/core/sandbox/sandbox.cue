package sandbox

Sandbox: {
	name: string

	nestedStruct: {
		foo: string
	}

	anythingPlz: { ... }

	someMap?: {
		[string]: int32
	}

	#OperatorState: {
			// lastEvaluation is the ResourceVersion last evaluated
			lastEvaluation: string
			// state describes the state of the lastEvaluation.
			// It is limited to three possible states for machine evaluation.
			state: "success" | "in_progress" | "failed"
			// descriptiveState is an optional more descriptive state field which has no requirements on format
			descriptiveState?: string
			// details contains any extra information that is operator-specific
			details?: {
				[string]: _
			}
		}
		// operatorStates is a map of operator ID to operator state evaluations.
		// Any operator which consumes this kind SHOULD add its state evaluation information to this field.
		operatorStates?: {
			[string]: #OperatorState
		}
}
