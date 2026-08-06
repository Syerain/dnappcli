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

// InfoMe 对应 server GET /api/v1/user/me 成功时返回的用户信息体。
// server 的 InfoMeResponse 内嵌 model.InfoMe，JSON 序列化会平铺展开，
// 因此这里以扁平结构对齐（RegisterTime 对应 server 的 time.Time -> JSON 字符串）。
type InfoMe struct {
	Uid          uint   `json:"uid"`
	Username     string `json:"username"`
	Nickname     string `json:"nickname"`
	Email        string `json:"email"`
	RegisterTime string `json:"register_time"`
	Role         string `json:"role"`
	GitHubID     *int64 `json:"github_id"`
	GitHubLogin  string `json:"github_login"`
}

// MeResponse 对应 GET /api/v1/user/me 的响应体。
type MeResponse struct {
	InfoMe
}
