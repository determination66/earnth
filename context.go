package earnth

import "net/http"

type Context struct {
	Req   *http.Request
	Write http.ResponseWriter
}
