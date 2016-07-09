package internal

type key int

var (
	ConfigContextKey    key = 0
	AuthContextKey      key = 1
	ClaimSetContextKey  key = 2
	ParamContextKey     key = 3
	AuthCheckContextKey key = 4
)
