package transfer

const (
	Queued State = iota
	Running
	Done
	Failed
	Canceled
)

// Label is the state's display name.
func (s State) Label() string {
	return [...]string{"Queued", "Transferring", "Done", "Failed", "Canceled"}[s]
}

// Finished reports whether the job won't change any more.
func (s State) Finished() bool { return s >= Done }
