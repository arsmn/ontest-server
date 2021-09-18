package handler

type jsonData struct {
	Data interface{} `json:"data"`
}

func data(d interface{}) *jsonData {
	return &jsonData{d}
}
