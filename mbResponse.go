package moonapi

type MbResponse struct {
	Data             []Problem   `json:"Data"`
	Total            int         `json:"Total"`
	AggregateResults interface{} `json:"AggregateResults"`
	Errors           interface{} `json:"Errors"`
}
