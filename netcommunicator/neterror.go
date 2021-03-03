package netcommunicator

type NetError struct {
	Message string
}

func (netErr NetError) Error() string {
	return netErr.Message
}
