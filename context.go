package earnth

import "net/http"

type Context struct {
	Req    *http.Request
	Writer http.ResponseWriter
	*matchInfo
	statusCode int
}
