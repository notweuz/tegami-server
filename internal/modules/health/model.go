package health

type Health struct {
	Status string
	Uptime string
	Commit string
}

func NewHealth(status, commit, uptime string) *Health {
	return &Health{
		Status: status,
		Uptime: uptime,
		Commit: commit,
	}
}
