package error

type WertErr struct {
	Val any
}

func NewWertErr(val any) *WertErr {
	return &WertErr{
		Val: val,
	}
}

func (r *WertErr) Error() string {
	return "<WertErr>"
}
