package models

type Omafile struct {
	Id     int    `json:"id"`
	Brand  string `json:"Brand"`
	Model  string `json:"Model"`
	OMAKey string `json:"OMAKey"`
	ImgKey string `json:"ImgKey"`
	ImgURL string `json:"ImgURL"`
}

func NewOmafile(parametr ...string) Omafile {
	return Omafile{
		Brand:  parametr[0],
		Model:  parametr[1],
		OMAKey: parametr[2],
		ImgKey: parametr[3],
	}
}

func GetOmafile(id int, parametr ...string) Omafile {
	return Omafile{
		Id:     id,
		Brand:  parametr[0],
		Model:  parametr[1],
		OMAKey: parametr[2],
		ImgKey: parametr[3],
	}
}
