package apirequest

type RefreshTokenRequest struct {
	Token string `json:"refresh_token" binding:"required"`
}
