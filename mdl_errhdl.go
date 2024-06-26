package earnth

type ErrHdlMiddlewareBuilder struct {
	resp map[int][]byte
}

func NewErrHdlMiddlewareBuilder() *ErrHdlMiddlewareBuilder {
	return &ErrHdlMiddlewareBuilder{
		// 这里可以非常大方，因为在预计中用户会关心的错误码不可能超过 64
		resp: make(map[int][]byte, 64),
	}
}

func (m *ErrHdlMiddlewareBuilder) RegisterError(code int, resp []byte) *ErrHdlMiddlewareBuilder {
	m.resp[code] = resp
	return m
}

func (m *ErrHdlMiddlewareBuilder) Build() MiddlewareFunc {
	return func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
