package tagflag

type userError struct {
	msg string
}

func (ue userError) Error() string {
	return ue.msg
}
