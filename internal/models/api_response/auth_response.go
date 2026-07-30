package apiresponse

type RefreshResponse struct {
	AccessToken  string `json:"accees_token"`
	RefreshToken string `json:"refresh_token"`
}
