package auth

type Response struct {
	Token string `json:"token" example:"jwt token string"`
}

func NewResponse(token string) Response {
	return Response{
		Token: token,
	}
}