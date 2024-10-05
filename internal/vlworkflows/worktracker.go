package vlworkflows

type workTracker struct {
	workWasDone bool
}

func (wt *workTracker) Work() {
	wt.workWasDone = true
}

func (wt *workTracker) WorkIfNoError(err error) {
	if err == nil {
		wt.workWasDone = true
	}
}

func (wt *workTracker) WorkWasDone() bool {
	return wt.workWasDone
}

func (wt *workTracker) AwaitFunc() func() bool {
	return func() bool {
		return wt.WorkWasDone()
	}
}
