package drugs_import

type MedElementImporter struct {
	serviceUrl      string            //URL к методу получения лекарств
	requestParams   map[string]string //Параметры запроса
	maxElementCount int               //Общее количество лекарств
	timeoutSeconds  int               //Количество секунд для таймаута
	timeoutItemsCnt int               //Количество элементов, через которое нужно делать таймаут
}

type MedElementResponse struct {
	Data   string `json:"data"`
	Result string `json:"result"`
}

type Drug struct {
	Guid   string
	Name   string
	Vendor string
	ATH    string
}
