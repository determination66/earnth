package earnth

type MiddlewareFunc func(next HandleFunc) HandleFunc
