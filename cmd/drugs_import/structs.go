package drugs_import

import (
	"database/sql"
	"github.com/warete/pharm/config"
)

type MedElementImporter struct {
	serviceUrl      string            //URL к методу получения лекарств
	requestParams   map[string]string //Параметры запроса
	maxElementCount int               //Общее количество лекарств
	timeoutSeconds  int               //Количество секунд для таймаута
	timeoutItemsCnt int               //Количество элементов, через которое нужно делать таймаут
	AppConfig       *config.Config
	DBConnection    *sql.DB
}

type MedElementResponse struct {
	Data   string `json:"data"`
	Result string `json:"result"`
}
