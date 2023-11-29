package models

type Application struct {
	Username        string `json:"username"`
	Firstname       string `json:"firstname"`
	Lastname        string `json:"lastname"`
	Birthday        string `json:"birthday"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Country         string `json:"country"`
	State           string `json:"state"`
	Zipcode         string `json:"zipcode"`
	Twitter         string `json:"twitter"`
	Instagram       string `json:"instagram"`
	Website         string `json:"website"`
	Documenttype    string `json:"document_type"`
	ExplicitContent string `json:"explicit_content"`
	FileFront       string `json:"file_front"`
	FileBack        string `json:"file_back"`
	FileHandwritten string `json:"file_handwritten"`
	FileVideo       string `json:"file_video"`
}
