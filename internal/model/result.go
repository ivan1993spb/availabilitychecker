package model

type Result struct {
	Task *Task

	// Check status (2 values):
	//   - fail (check has failed)
	//   - ok   (check has succeeded)
	Status Status

	// Error message: the reason why the check has failed
	// or an empty string
	FailMessage string
}
