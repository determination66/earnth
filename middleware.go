package earnth

type Middleware func(next HandleFunc) HandleFunc
