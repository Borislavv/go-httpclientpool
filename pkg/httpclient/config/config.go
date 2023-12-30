package httpclientconfig

type Config struct {
	ReqTimeoutMs int `arg:"-t,env:HTTP_CLIENT_POOL_REQ_TIMEOUT" default:"1000"`
	PoolInitSize int `arg:"-c,env:HTTP_CLIENT_POOL_INIT_SIZE"   default:"32"`
	PoolMaxSize  int `arg:"-c,env:HTTP_CLIENT_POOL_MAX_SIZE"    default:"10240"`
}
