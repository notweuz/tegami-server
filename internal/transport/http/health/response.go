package health

type InfoResponse struct {
	Status string `json:"status" example:"UP"`
	Commit string `json:"commit" example:"2acec66b3b1efbe040f738b18cbb7f60f841fddc"`
	Uptime string `json:"uptime" example:"12.666100458s"`
}

func NewInfoResponse(status, commit, uptime string) *InfoResponse {
	return &InfoResponse{
		Status: status,
		Commit: commit,
		Uptime: uptime,
	}
}
