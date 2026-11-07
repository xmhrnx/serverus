package terminal

// frame is what one paint needs, captured under the lock.
type frame struct {
	lines         [][]cell
	cols          int
	cursorX       int
	cursorY       int // row in lines; -1 when off screen
	cursorVisible bool
	scrollback    int // lines of history available
	top           int // absolute index of lines[0]
}
