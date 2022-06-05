package public

const (
	ValidatorKey         = "ValidatorKey"
	TranslatorKey        = "TranslatorKey"
	AdminSessionInfoKey  = "AdminSessionInfoKey" //服务端session key，会被哈希
	ClientSessionInfoKey = "gatewaysession"      //客户端session key, 存储在客户端到cookie
	SessionInfoTime      = 60 * 60 * 24 * 3      //服务端session过期时间
)
