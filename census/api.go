package census

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// CensusAPI представляет собой клиент для API переписи населения
type CensusAPI struct {
	apiKey string
	client *http.Client
}

// NewCensusAPI создает новый экземпляр клиента CensusAPI
func NewCensusAPI(apiKey string) *CensusAPI {
	return &CensusAPI{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

// NewCensusAPIFromEnv создает новый экземпляр клиента CensusAPI, используя ключ API из переменной окружения
func NewCensusAPIFromEnv() (*CensusAPI, error) {
	apiKey := os.Getenv("CENSUS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("переменная окружения CENSUS_API_KEY не установлена")
	}
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

	resp, err := c.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка при отправке запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API вернул статус %d", resp.StatusCode)
	}

	// Census API возвращает массив массивов, первый массив содержит заголовки
	var rawData [][]string
	if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
		return nil, fmt.Errorf("ошибка при декодировании ответа: %w", err)
	}

	if len(rawData) < 2 {
		return nil, fmt.Errorf("API вернул пустой результат")
	}

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
	endpoint := "https://api.census.gov/data/2021/acs/acs1"

	params := url.Values{}
	params.Add("get", "NAME,B01001_001E")

	if stateID != "" {
		params.Add("for", "county:*")
		params.Add("in", fmt.Sprintf("state:%s", stateID))
	} else {
		params.Add("for", "county:*")
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
	states, err := c.GetStatePopulation("")
	if err != nil {
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

// Formatter определяет интерфейс для форматирования результатов
type Formatter interface {
	Format(data interface{}) string
	CountLetters(word, letters string) map[rune]int
}

// TextFormatter реализует форматирование в текстовом виде
type TextFormatter struct{}

// Format форматирует данные о населении в текстовом виде
func (f *TextFormatter) Format(data interface{}) string {
	switch v := data.(type) {
	case []PopulationData:
		result := "Результаты запроса к Census API:\n\n"

		for _, item := range v {
			result += fmt.Sprintf("Название: %s\n", item.Name)
			result += fmt.Sprintf("Население: %s\n", item.Population)

			if item.State != "" {
				result += fmt.Sprintf("Код штата: %s\n", item.State)
			}

			if item.County != "" {
				result += fmt.Sprintf("Код округа: %s\n", item.County)
			}

			result += "\n"
		}

		return result
	case []DatasetInfo:
		result := "Доступные наборы данных Census API:\n\n"

		for _, item := range v {
			result += fmt.Sprintf("Идентификатор: %s\n", item.Dataset)
			result += fmt.Sprintf("Название: %s\n", item.Title)

			if item.Description != "" {
				result += fmt.Sprintf("Описание: %s\n", item.Description)
			}

			result += fmt.Sprintf("Доступные годы: %s\n", strings.Join(item.YearsAvailable, ", "))
			result += "\n"
		}

		return result
	case map[string]VariableInfo:
		result := "Доступные переменные Census API:\n\n"

		for name, info := range v {
			result += fmt.Sprintf("Переменная: %s\n", name)
			result += fmt.Sprintf("Название: %s\n", info.Label)

			if info.Concept != "" {
				result += fmt.Sprintf("Концепция: %s\n", info.Concept)
			}

			if info.Description != "" {
				result += fmt.Sprintf("Описание: %s\n", info.Description)
			}

			if info.Group != "" {
				result += fmt.Sprintf("Группа: %s\n", info.Group)
			}

			result += "\n"
		}

		return result
	case []GeographyLevel:
		result := "Доступные географические уровни Census API:\n\n"

		for _, level := range v {
			result += fmt.Sprintf("Уровень: %s\n", level.Name)
			result += fmt.Sprintf("Описание: %s\n", level.Description)

			if len(level.RequiredFor) > 0 {
				result += fmt.Sprintf("Требуется указать: %s\n", strings.Join(level.RequiredFor, ", "))
			}

			result += fmt.Sprintf("Поддержка wildcard: %s\n", boolToYesNo(level.Wildcards))
			result += "\n"
		}

		return result
	case []map[string]string:
		result := "Результаты пользовательского запроса к Census API:\n\n"

		if len(v) == 0 {
			return result + "Нет данных"
		}

		// Выводим заголовки (ключи из первой записи)
		var headers []string
		for header := range v[0] {
			headers = append(headers, header)
		}

		for i, item := range v {
			result += fmt.Sprintf("Запись %d:\n", i+1)
			for _, header := range headers {
				if value, ok := item[header]; ok {
					result += fmt.Sprintf("%s: %s\n", header, value)
				}
			}
			result += "\n"
		}

		return result
	case map[rune]int:
		result := "Результаты подсчета букв:\n\n"

		for char, count := range v {
			result += fmt.Sprintf("Буква '%c': %d\n", char, count)
		}

		return result
	default:
		return fmt.Sprintf("Неподдерживаемый тип данных: %T", data)
	}
}

// NewTextFormatter создает новый форматтер в текстовом виде
func NewTextFormatter() Formatter {
	return &TextFormatter{}
}

// CountLetters подсчитывает количество определенных букв в слове
func (f *TextFormatter) CountLetters(word, letters string) map[rune]int {
	result := make(map[rune]int)

	// Инициализация счетчиков
	for _, letter := range letters {
		result[letter] = 0
	}

	// Подсчет букв
	for _, char := range word {
		for _, letter := range letters {
			if char == letter {
				result[letter]++
			}
		}
	}

	return result
}

// boolToYesNo преобразует булево значение в строку "Да" или "Нет"
func boolToYesNo(value bool) string {
	if value {
		return "Да"
	}
	return "Нет"
}
