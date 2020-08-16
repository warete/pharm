package drugs_import

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/warete/pharm/cmd/pharm/repository/drug"
	"github.com/warete/pharm/config"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

//Конструктор
func Init(serviceUrl string, requestParams map[string]string, appConfig *config.Config) (*MedElementImporter, error) {
	connection, err := sql.Open(
		appConfig.DB.Type,
		fmt.Sprintf(
			"%s:%s@tcp(%s)/%s",
			appConfig.DB.User,
			appConfig.DB.Password,
			appConfig.DB.Host,
			appConfig.DB.DBName,
		),
	)
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	return &MedElementImporter{
		serviceUrl:      serviceUrl,
		requestParams:   requestParams,
		maxElementCount: 0,
		timeoutSeconds:  3,
		timeoutItemsCnt: 1500,
		AppConfig:       appConfig,
		DBConnection:    connection,
	}, nil
}

//Запускает импорт
func (i *MedElementImporter) Run() error {
	//Мапа со всеми лекарствами
	drugItems := make(map[string]*drug.Drug)
	//Мьютекс для конкурентной записи в мапу
	drugItemsMutex := sync.RWMutex{}

	//Получаем первую пачку лекарств и смотрим сколько их всего
	firstDrugs, err := i.getDrugs(0)
	if err != nil {
		return err
	}
	for _, item := range firstDrugs {
		drugItems[item.Guid] = item
	}
	drugsStepCnt := len(firstDrugs)

	//Получаем все остальные лекарства до конца
	var wg sync.WaitGroup
	if i.maxElementCount > 0 {
		for k := drugsStepCnt; k <= i.maxElementCount; k = k + drugsStepCnt {
			//Делаем таймаут, чтобы не возвращались ошибки из-за большого количества запросов
			if i.timeoutItemsCnt > 0 && k%i.timeoutItemsCnt == 0 && i.timeoutSeconds > 0 {
				time.Sleep(time.Duration(i.timeoutSeconds) * time.Second)
			}

			wg.Add(1)
			go func(i *MedElementImporter, k int) {
				defer wg.Done()

				drugs, err := i.getDrugs(k)
				if err != nil {
					panic(err)
				}

				drugItemsMutex.Lock()
				for _, item := range drugs {
					drugItems[item.Guid] = item
				}
				drugItemsMutex.Unlock()
			}(i, k)
		}
	}
	wg.Wait()
	fmt.Println(len(drugItems))

	//TODO: write to DB

	return nil
}

//Получение данных с сервера
func (i *MedElementImporter) getDataFromService(customParams map[string]string) (string, error) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", i.serviceUrl+"?"+GetQueryStringFromMap(customParams), nil)

	if err != nil {
		return "", err
	}
	request.Header.Set("x-requested-with", "XMLHttpRequest")
	res, err := client.Do(request)
	if res == nil {
		return "", errors.New("empty response")
	}
	if res.StatusCode != 200 {
		return "", err
	}
	defer res.Body.Close()

	respData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(respData), nil
}

//Обработка ответа от сервера
func (i *MedElementImporter) processData(data string) *MedElementResponse {
	drugsResponse := &MedElementResponse{}
	json.Unmarshal([]byte(data), drugsResponse)

	//Либа для парсинга принимает только цельный html-документ
	//поэтому добавляем теги html, head и body
	drugsResponse.Data = strings.ReplaceAll(drugsResponse.Data, "\n", "")
	drugsResponse.Data = strings.ReplaceAll(drugsResponse.Data, "        ", "")
	drugsResponse.Data = "<html><head></head><body>" + drugsResponse.Data + "</body></html>"

	return drugsResponse
}

//Парсинг лекарств
func (i *MedElementImporter) parseItemsFromData(data *MedElementResponse) ([]*drug.Drug, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data.Data))
	if err != nil {
		return nil, err
	}

	//Получим общее количество лекарств
	if i.maxElementCount == 0 {
		searchedResultsStr := strings.TrimSpace(doc.Find(".results .row.sort .pull-right").Text())

		if strings.Contains(searchedResultsStr, "Найдено ") {
			i.maxElementCount, err = strconv.Atoi(strings.ReplaceAll(searchedResultsStr, "Найдено ", ""))
			if err != nil {
				return nil, err
			}
		}
	}

	//Ищем товары
	var drugs []*drug.Drug
	doc.Find(".row.results__result").Each(func(i int, s *goquery.Selection) {
		drugElement := &drug.Drug{}
		drugElement.Name = s.Find("a.results__title-link").Text()
		s.Find("span.text-muted").Each(func(i int, s *goquery.Selection) {
			nodeText := strings.TrimSpace(s.Text())
			if strings.Contains(nodeText, "Производитель:") {
				drugElement.Vendor = strings.ReplaceAll(nodeText, "Производитель: ", "")
			}
			if strings.Contains(nodeText, "АТХ:") {
				drugElement.ATH = strings.ReplaceAll(nodeText, "АТХ: ", "")
			}
		})
		//Сгенерим uuid
		drugElement.Guid = uuid.New().String()
		drugElement.Active = true
		drugs = append(drugs, drugElement)
	})

	return drugs, nil
}

//Обёртка для получения среза лекарств
func (i *MedElementImporter) getDrugs(skip int) ([]*drug.Drug, error) {
	customParams := make(map[string]string)
	for key, value := range i.requestParams {
		customParams[key] = value
	}
	customParams["skip"] = strconv.Itoa(skip)

	jsonDataFromService, err := i.getDataFromService(customParams)
	if err != nil {
		return nil, err
	}

	drugs, err := i.parseItemsFromData(i.processData(jsonDataFromService))
	if err != nil {
		return nil, err
	}

	return drugs, nil
}

func (i *MedElementImporter) GetAll() ([]*drug.Drug, error) {
	drugRepo, err := drug.New(i.DBConnection)
	if err != nil {
		return nil, err
	}
	return drugRepo.GetAllActiveDrugs()
}
