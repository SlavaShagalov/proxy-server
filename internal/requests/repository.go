package requests

type Req struct {
	Method     string
	Path       string
	GetParams  map[string][]string
	Headers    map[string][]string
	Cookies    map[string]string
	PostParams map[string][]string
	Body       string
}

type Resp struct {
	Code    int
	Message string
	Headers map[string][]string
	Body    string
}

type Request struct {
	ID   string `json:"id"`
	Req  Req    `json:"request"`
	Resp Resp   `json:"response"`
}

type Repository interface {
	Create(*Request) error
	Get(id string) (*Request, error)
	List() ([]Request, error)
}
