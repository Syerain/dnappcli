package model

// RegisterBody 与 server internal/model.RegisterBody 结构一致，
// 用于构造 /api/register 的请求体。
// 注意：server 的 model 位于 internal/ 下无法被外部模块引用，
// 因此在此按契约复刻一份。
type RegisterBody struct {
	Username     string
	Nickname     string
	Password     string
	Registercode string
	Registerway  string
}

// Response 描述 server 统一的响应外壳。
// 当前 server 返回 {"message": "..."}。
type Response struct {
	Message string `json:"message"`
}

// LoginBody 与 server internal/model/LoginBody 结构一致，
// 用于构造 /api/login 的请求体。
type LoginBody struct {
	Loginway string `json:"loginway"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse 对应 /api/login 成功时的返回体。
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// TokenData 对应保存在 data.yaml 中的 token 信息。
type TokenData struct {
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
}
