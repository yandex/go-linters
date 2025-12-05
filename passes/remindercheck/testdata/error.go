package testdata

/*
TODO: multi-line comments are ignored
*/

// Error

// TODO// want `TODO must be contains right parts. Required template: '// TODO: TASKID-1: comment'`

// FIXME:							// want `FIXME must include task id. Required template: '// FIXME: TASKID-1: comment'`

// tODo: TASKID-1 make a coffee		// want `keyword 'tODo' must be upper case. Required template: '// TODO: TASKID-1: comment'`
// TODO: TASKID-0 make a coffee		// want `TODO must use task id number greater zero: TASKID-0. Required template: '// TODO: TASKID-1: comment'`
// TODO: (TASKID-1) make a coffee 	// want `TODO must be contains right parts. Required template: '// TODO: TASKID-1: comment'`
// TODO(TASKID-1): make a coffee	// want `TODO must be contains right parts. Required template: '// TODO: TASKID-1: comment'`
// TODO: TASKID make a coffee	    // want `TODO must be contains right parts. Required template: '// TODO: TASKID-1: comment'`
// TODO: TASKID100                  // want `TODO must use valid task id: TASKID100. Required template: '// TODO: TASKID-1: comment'`
// FIXME: TASKID-1          		// want `FIXME must describe what needs to remind. Required template: '// FIXME: TASKID-1: comment'`
func makeBadTea() string {
	// TODO: tmp-ticket make a coffee   // want `TODO must be contains right parts. Required template: '// TODO: TASKID-1: comment'`
	panic("make a bad tea is not implemented")
}
