package health

type InfoResponse struct {
	Status string `json:"status"`
	Commit string `json:"commit"`
	Uptime string `json:"uptime"`
}

func NewInfoResponse(status, commit, uptime string) *InfoResponse {
	return &InfoResponse{
		Status: status,
		Commit: commit,
		Uptime: uptime,
	}
}