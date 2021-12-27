package api

type API struct {
}

type APIInterface interface {
}

func NewAPI() APIInterface {
	return &API{}
}
