package userapi

// LoginReq 用于编码成JSON，作为BODY发起登录请求
type LoginReq struct {
	Clientid     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Email        string `json:"email"`
	Password     string `json:"password"`
}

// LoginResp 解析登录接口返回的JSON
type LoginResp struct {
	AccessToken  string `json:"access_token"`
	Email        string `json:"email"`
	ExpiresIn    int64  `json:"expires_in"`
	Message      string `json:"message"`
	RefreshToken string `json:"refresh_token"`
	ResCode      int    `json:"res_code"`
	TokenType    string `json:"token_type"`
	UserID       string `json:"user_id"`
}

// NewLoginReq hh.
func NewLoginReq(clientid, secret, email, pwd string) interface{} {
	req := &LoginReq{
		Clientid:     clientid,
		ClientSecret: secret,
		Email:        email,
		Password:     pwd,
	}
	return req
}
