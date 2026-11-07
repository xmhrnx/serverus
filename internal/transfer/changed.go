package transfer

func (q *Queue) changed() {
	if q.onChange != nil {
		q.onChange()
	}
}
