package out

type SubscribeError struct {
	Code    int
	Details string
}

func (e *SubscribeError) Error() string {
	return e.Details
}
