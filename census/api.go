package census

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// Константы для ключей логирования
const (
	key_state_id    = "state_id"
	key_endpoint    = "endpoint"
	key_err         = "err"
	key_status_code = "status_code"
	key_count       = "count"
	key_name        = "name"
	key_search_name = "search_name"
	key_dataset     = "dataset"
	key_year        = "year"
	key_request     = "request"
)

// CensusAPI представляет собой клиент для API переписи населения
type CensusAPI struct {
	apiKey string
	client *http.Client
}

// NewCensusAPI создает новый экземпляр клиента CensusAPI
func NewCensusAPI(apiKey string) *CensusAPI {
	slog.Debug("Создание нового клиента CensusAPI с ключом API")
	return &CensusAPI{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// NewCensusAPIFromEnv создает новый экземпляр клиента CensusAPI, используя ключ API из переменной окружения
func NewCensusAPIFromEnv() (*CensusAPI, error) {
	slog.Debug("Создание клиента CensusAPI из переменной окружения")
	apiKey := os.Getenv("CENSUS_API_KEY")
	if apiKey == "" {
		slog.Error("Переменная окружения CENSUS_API_KEY не установлена")
		return nil, fmt.Errorf("переменная окружения CENSUS_API_KEY не установлена")
	}
	slog.Debug("Ключ API получен из переменной окружения")
	return NewCensusAPI(apiKey), nil
}

// PopulationData представляет собой данные о населении
type PopulationData struct {
	Name       string `json:"NAME"`
	Population string `json:"B01001_001E"`
	State      string `json:"state,omitempty"`
	County     string `json:"county,omitempty"`
}

// CustomDataRequest представляет запрос пользовательских данных переписи
type CustomDataRequest struct {
	Variables []string          // Список переменных для запроса
	Dataset   string            // Набор данных (например, "acs/acs1", "dec/sf1")
	Year      string            // Год данных
	GeoLevel  string            // Географический уровень (state, county, tract и т.д.)
	GeoFilter map[string]string // Фильтр географии (например, {"state": "06", "county": "*"})
}

// DatasetInfo содержит информацию о доступном наборе данных
type DatasetInfo struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Dataset        string   `json:"dataset"`
	YearsAvailable []string `json:"years_available"`
}

// VariableInfo содержит информацию о переменной
type VariableInfo struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Concept     string `json:"concept,omitempty"`
	Group       string `json:"group,omitempty"`
}

// GeographyLevel содержит информацию о доступном географическом уровне
type GeographyLevel struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RequiredFor []string `json:"required_for,omitempty"`
	Wildcards   bool     `json:"wildcards"`
}

// CensusAPIClient определяет интерфейс для клиента Census API
type CensusAPIClient interface {
	// GetStatePopulation возвращает данные о населении для указанного штата
	GetStatePopulation(stateID string) ([]PopulationData, error)
	// GetCountyPopulation возвращает данные о населении для округов в указанном штате
	GetCountyPopulation(stateID string) ([]PopulationData, error)
	// SearchStateByName ищет штат по названию (полному или частичному)
	SearchStateByName(name string) ([]PopulationData, error)
	// GetAvailableDatasets возвращает список доступных наборов данных
	GetAvailableDatasets() ([]DatasetInfo, error)
	// GetVariables возвращает список доступных переменных для набора данных
	GetVariables(dataset, year string) (map[string]VariableInfo, error)
	// GetGeographyLevels возвращает доступные географические уровни для набора данных
	GetGeographyLevels(dataset, year string) ([]GeographyLevel, error)
	// GetCustomData позволяет запросить пользовательские данные
	GetCustomData(request CustomDataRequest) ([]map[string]string, error)
}

// GetStatePopulation возвращает данные о населении для указанного штата
func (c *CensusAPI) GetStatePopulation(stateID string) ([]PopulationData, error) {
	slog.Info("Получение данных о населении штата", key_state_id, stateID)

	endpoint := "https://api.census.gov/data/2021/acs/acs1"

	params := url.Values{}
	params.Add("get", "NAME,B01001_001E")

	if stateID != "" {
		params.Add("for", fmt.Sprintf("state:%s", stateID))
	} else {
		params.Add("for", "state:*")
	}

	params.Add("key", c.apiKey)

	requestURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	slog.Debug("Отправка запроса к Census API", key_endpoint, endpoint)

	resp, err := c.client.Get(requestURL)
	if err != nil {
		slog.Error("Ошибка при отправке запроса",
			key_err, err,
			key_endpoint, endpoint)
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("API вернул неуспешный статус",
			key_status_code, resp.StatusCode,
			key_endpoint, endpoint)
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	// Census API возвращает массив массивов, первый массив содержит заголовки
	var rawData [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		slog.Error("Ошибка при декодировании ответа",
			key_err, err,
			key_endpoint, endpoint)
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(rawData) < 2 {
		slog.Error("API вернул пустой результат",
			key_endpoint, endpoint)
		return nil, fmt.Errorf("API вернул пустой результат")
	}

	slog.Debug("Получены данные из Census API",
		key_count, len(rawData)-1,
		key_endpoint, endpoint)

	// Первый массив - это заголовки
	headers := rawData[0]

	// Преобразуем остальные массивы в структуры PopulationData
	var result []PopulationData
	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		if len(data) != len(headers) {
			continue
		}

		// Создаем мапу для удобного доступа к данным
		dataMap := make(map[string]string)
		for j, header := range headers {
			dataMap[header] = data[j]
		}

		popData := PopulationData{
			Name:       dataMap["NAME"],
			Population: dataMap["B01001_001E"],
			State:      dataMap["state"],
		}

		result = append(result, popData)
	}

	return result, nil
}

// GetCountyPopulation возвращает данные о населении для округов в указанном штате
func (c *CensusAPI) GetCountyPopulation(stateID string) ([]PopulationData, error) {
	slog.Info("Получение данных о населении округов", key_state_id, stateID)

	endpoint := "https://api.census.gov/data/2021/acs/acs1"

	params := url.Values{}
	params.Add("get", "NAME,B01001_001E")

	if stateID != "" {
		params.Add("for", fmt.Sprintf("county:*&in=state:%s", stateID))
	} else {
		// Запрос всех округов во всех штатах
		params.Add("for", "county:*")
	}

	params.Add("key", c.apiKey)

	requestURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	slog.Debug("Отправка запроса к Census API", key_endpoint, endpoint)

	resp, err := c.client.Get(requestURL)
	if err != nil {
		slog.Error("Ошибка при отправке запроса",
			key_err, err,
			key_endpoint, endpoint)
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("API вернул неуспешный статус",
			key_status_code, resp.StatusCode,
			key_endpoint, endpoint)
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	var rawData [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(rawData) < 2 {
		return nil, fmt.Errorf("API вернул пустой результат")
	}

	headers := rawData[0]

	var result []PopulationData
	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		if len(data) != len(headers) {
			continue
		}

		dataMap := make(map[string]string)
		for j, header := range headers {
			dataMap[header] = data[j]
		}

		popData := PopulationData{
			Name:       dataMap["NAME"],
			Population: dataMap["B01001_001E"],
			State:      dataMap["state"],
			County:     dataMap["county"],
		}

		result = append(result, popData)
	}

	return result, nil
}

// SearchStateByName ищет штат по названию (полному или частичному)
func (c *CensusAPI) SearchStateByName(name string) ([]PopulationData, error) {
	slog.Info("Поиск штата по названию", key_name, name)

	// Получаем все штаты
	states, err := c.GetStatePopulation("")
	if err != nil {
		slog.Error("Ошибка при получении данных о штатах для поиска",
			key_err, err,
			key_search_name, name)
		return nil, err
	}

	var result []PopulationData
	for _, state := range states {
		if strings.Contains(strings.ToLower(state.Name), strings.ToLower(name)) {
			result = append(result, state)
		}
	}

	return result, nil
}

// GetAvailableDatasets возвращает список доступных наборов данных
func (c *CensusAPI) GetAvailableDatasets() ([]DatasetInfo, error) {
	slog.Info("Получение списка доступных наборов данных")

	endpoint := "https://api.census.gov/data.json"

	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	type apiResponse struct {
		Dataset []struct {
			C_Dataset []struct {
				Title        string `json:"title"`
				Description  string `json:"description"`
				Distribution []struct {
					Title     string `json:"title"`
					AccessURL string `json:"accessURL"`
				} `json:"distribution"`
			} `json:"c_dataset"`
		} `json:"dataset"`
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	var result []DatasetInfo
	yearMap := make(map[string]map[string][]string) // map[dataset]map[baseURL][]years

	for _, ds := range response.Dataset {
		for _, cds := range ds.C_Dataset {
			for _, dist := range cds.Distribution {
				// Извлекаем информацию о наборе данных из URL
				// URL обычно имеет вид https://api.census.gov/data/[year]/[dataset]
				if dist.AccessURL == "" {
					continue
				}

				parts := strings.Split(dist.AccessURL, "/data/")
				if len(parts) != 2 {
					continue
				}

				pathParts := strings.Split(parts[1], "/")
				if len(pathParts) < 2 {
					continue
				}

				year := pathParts[0]
				dataset := strings.Join(pathParts[1:], "/")
				baseURL := parts[0] + "/data/" + dataset

				if _, ok := yearMap[dataset]; !ok {
					yearMap[dataset] = make(map[string][]string)
				}
				yearMap[dataset][baseURL] = append(yearMap[dataset][baseURL], year)
			}
		}
	}

	// Преобразуем собранную информацию в результат
	for dataset, baseURLMap := range yearMap {
		for _, years := range baseURLMap {
			result = append(result, DatasetInfo{
				Title:          dataset,
				Dataset:        dataset,
				YearsAvailable: years,
			})
		}
	}

	return result, nil
}

// GetVariables возвращает список доступных переменных для набора данных
func (c *CensusAPI) GetVariables(dataset, year string) (map[string]VariableInfo, error) {
	slog.Info("Получение списка переменных",
		key_dataset, dataset,
		key_year, year)

	if dataset == "" || year == "" {
		return nil, fmt.Errorf("необходимо указать набор данных и год")
	}

	endpoint := fmt.Sprintf("https://api.census.gov/data/%s/%s/variables.json", year, dataset)

	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d для %s", resp.StatusCode, endpoint)
	}

	type apiResponse struct {
		Variables map[string]struct {
			Label       string `json:"label"`
			Concept     string `json:"concept"`
			Description string `json:"description,omitempty"`
			Group       string `json:"group,omitempty"`
		} `json:"variables"`
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	result := make(map[string]VariableInfo)
	for name, info := range response.Variables {
		result[name] = VariableInfo{
			Name:        name,
			Label:       info.Label,
			Concept:     info.Concept,
			Description: info.Description,
			Group:       info.Group,
		}
	}

	return result, nil
}

// GetGeographyLevels возвращает доступные географические уровни для набора данных
func (c *CensusAPI) GetGeographyLevels(dataset, year string) ([]GeographyLevel, error) {
	slog.Info("Получение доступных географических уровней",
		key_dataset, dataset,
		key_year, year)

	if dataset == "" || year == "" {
		return nil, fmt.Errorf("необходимо указать набор данных и год")
	}

	endpoint := fmt.Sprintf("https://api.census.gov/data/%s/%s/geography.json", year, dataset)

	resp, err := c.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d для %s", resp.StatusCode, endpoint)
	}

	type apiResponse struct {
		GeographyLevels map[string]struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			RequiredFor []string `json:"required_for,omitempty"`
			Wildcards   bool     `json:"wildcards"`
		} `json:"fips"`
	}

	var response apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	var result []GeographyLevel
	for _, info := range response.GeographyLevels {
		result = append(result, GeographyLevel{
			Name:        info.Name,
			Description: info.Description,
			RequiredFor: info.RequiredFor,
			Wildcards:   info.Wildcards,
		})
	}

	return result, nil
}

// GetCustomData позволяет запросить пользовательские данные
func (c *CensusAPI) GetCustomData(request CustomDataRequest) ([]map[string]string, error) {
	slog.Info("Получение пользовательских данных",
		key_request, request)

	if len(request.Variables) == 0 {
		return nil, fmt.Errorf("необходимо указать хотя бы одну переменную")
	}

	if request.Dataset == "" {
		return nil, fmt.Errorf("необходимо указать набор данных")
	}

	if request.Year == "" {
		return nil, fmt.Errorf("необходимо указать год")
	}

	if request.GeoLevel == "" {
		return nil, fmt.Errorf("необходимо указать географический уровень")
	}

	endpoint := fmt.Sprintf("https://api.census.gov/data/%s/%s", request.Year, request.Dataset)

	params := url.Values{}
	params.Add("get", strings.Join(request.Variables, ","))
	params.Add("for", fmt.Sprintf("%s:%s", request.GeoLevel, request.GeoFilter[request.GeoLevel]))

	for key, value := range request.GeoFilter {
		if key != request.GeoLevel {
			params.Add("in", fmt.Sprintf("%s:%s", key, value))
		}
	}

	params.Add("key", c.apiKey)

	requestURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := c.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	var rawData [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(rawData) < 2 {
		return nil, fmt.Errorf("API вернул пустой результат")
	}

	headers := rawData[0]
	result := make([]map[string]string, 0, len(rawData)-1)

	for i := 1; i < len(rawData); i++ {
		data := rawData[i]
		if len(data) != len(headers) {
			continue
		}

		dataMap := make(map[string]string)
		for j, header := range headers {
			dataMap[header] = data[j]
		}

		result = append(result, dataMap)
	}

	return result, nil
}
