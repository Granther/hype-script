package error

type ReturnErr struct {
	Val any
}

func NewReturnErr(val any) *ReturnErr {
	return &ReturnErr{
		Val: val,
	}
}

func (r *ReturnErr) Error() string {
	return "ReturnErr"
}