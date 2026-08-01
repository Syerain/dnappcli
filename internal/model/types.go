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
