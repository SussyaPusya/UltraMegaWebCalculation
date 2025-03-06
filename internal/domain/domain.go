package domain

type JsonReq struct {
	Expression string `json:"expression"`
}

type Task struct {
	ID             string  `json:"id"`
	Arg1           float64 `json:"arg1"`
	Arg2           float64 `json:"arg2"`
	Operation      string  `json:"operation"`
	Operation_time string  `json:"operation_time"`

	ResultChan chan float64 `json:"-"`
}

type Expression struct {
	ID         string   `json:"id"`
	Status     string   `json:"status"`
	Result     *float64 `json:"result"`
	Expression string   `json:"-"`
}

var DoneTask []Expression

type Key string

const (
	Logger Key = "logger"

	RequestID Key = "request_id"
)
