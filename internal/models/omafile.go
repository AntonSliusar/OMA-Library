package models

type Omafile struct {
	Id    int    `json:"id"`
	Brand string `json:"Brand"`
	Model string `json:"Model"`
	Key   string `json:"Key"`
}

func NewOmafile(parametr ...string) Omafile {
	return Omafile{
		Brand: parametr[0],
		Model: parametr[1],
		Key:   parametr[2],
	}
}

func GetOmafile(id int, parametr ...string) Omafile {
	return Omafile{
		Id:    id,
		Brand: parametr[0],
		Model: parametr[1],
		Key:   parametr[2],
	}
}
