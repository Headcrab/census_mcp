package census

// MockCensusAPI - это реализация API Census для тестов, которая возвращает тестовые данные
// Реализует интерфейс CensusAPIClient
type MockCensusAPI struct{}

// NewMockCensusAPI создает новый экземпляр мок-клиента API Census
func NewMockCensusAPI() *MockCensusAPI {
	return &MockCensusAPI{}
}

// GetStatePopulation возвращает тестовые данные о населении штатов
func (m *MockCensusAPI) GetStatePopulation(stateID string) ([]PopulationData, error) {
	// Тестовые данные о населении штатов
	states := []PopulationData{
		{
			Name:       "Alabama",
			Population: "5024279",
			State:      "01",
		},
		{
			Name:       "Alaska",
			Population: "733391",
			State:      "02",
		},
		{
			Name:       "Arizona",
			Population: "7151502",
			State:      "04",
		},
		{
			Name:       "California",
			Population: "39538223",
			State:      "06",
		},
		{
			Name:       "New York",
			Population: "20201249",
			State:      "36",
		},
		{
			Name:       "Texas",
			Population: "29145505",
			State:      "48",
		},
	}

	// Если указан конкретный штат, возвращаем только его
	if stateID != "" {
		for _, state := range states {
			if state.State == stateID {
				return []PopulationData{state}, nil
			}
		}
		// Если запрошенный штат не найден, возвращаем штат по умолчанию
		return []PopulationData{states[0]}, nil
	}

	return states, nil
}

// GetCountyPopulation возвращает тестовые данные о населении округов
func (m *MockCensusAPI) GetCountyPopulation(stateID string) ([]PopulationData, error) {
	// Тестовые данные о населении округов в Калифорнии
	counties := []PopulationData{
		{
			Name:       "Los Angeles County",
			Population: "10014009",
			State:      "06",
			County:     "037",
		},
		{
			Name:       "San Diego County",
			Population: "3298634",
			State:      "06",
			County:     "073",
		},
		{
			Name:       "Orange County",
			Population: "3186989",
			State:      "06",
			County:     "059",
		},
		{
			Name:       "King County",
			Population: "2252782",
			State:      "53",
			County:     "033",
		},
		{
			Name:       "Harris County",
			Population: "4713325",
			State:      "48",
			County:     "201",
		},
	}

	// Если указан конкретный штат, возвращаем только его округа
	if stateID != "" {
		var result []PopulationData
		for _, county := range counties {
			if county.State == stateID {
				result = append(result, county)
			}
		}
		if len(result) > 0 {
			return result, nil
		}
		// Если запрошенный штат не найден, возвращаем округ по умолчанию
		return []PopulationData{counties[0]}, nil
	}

	return counties, nil
}

// SearchStateByName ищет штаты по названию в тестовых данных
func (m *MockCensusAPI) SearchStateByName(name string) ([]PopulationData, error) {
	states, _ := m.GetStatePopulation("")

	var result []PopulationData
	for _, state := range states {
		// Простой поиск подстроки в названии (без учета регистра)
		if contains(state.Name, name) {
			result = append(result, state)
		}
	}

	return result, nil
}

// GetAvailableDatasets возвращает список доступных наборов данных (тестовые данные)
func (m *MockCensusAPI) GetAvailableDatasets() ([]DatasetInfo, error) {
	// Тестовые данные о доступных наборах данных
	datasets := []DatasetInfo{
		{
			Title:          "American Community Survey 1-Year Estimates",
			Description:    "Annual survey covering demographic, social, economic, and housing data",
			Dataset:        "acs/acs1",
			YearsAvailable: []string{"2019", "2020", "2021"},
		},
		{
			Title:          "Decennial Census",
			Description:    "Complete count of the US population conducted every 10 years",
			Dataset:        "dec/sf1",
			YearsAvailable: []string{"2000", "2010", "2020"},
		},
		{
			Title:          "Population Estimates Program",
			Description:    "Annual population estimates between decennial censuses",
			Dataset:        "pep/population",
			YearsAvailable: []string{"2018", "2019", "2020", "2021"},
		},
	}

	return datasets, nil
}

// GetVariables возвращает список доступных переменных для набора данных (тестовые данные)
func (m *MockCensusAPI) GetVariables(dataset, year string) (map[string]VariableInfo, error) {
	// Тестовые данные о доступных переменных
	variables := map[string]VariableInfo{
		"B01001_001E": {
			Name:        "B01001_001E",
			Label:       "Total Population",
			Concept:     "SEX BY AGE",
			Description: "Total population count",
			Group:       "B01001",
		},
		"B01002_001E": {
			Name:        "B01002_001E",
			Label:       "Median Age",
			Concept:     "MEDIAN AGE BY SEX",
			Description: "Median age of total population",
			Group:       "B01002",
		},
		"B02001_001E": {
			Name:        "B02001_001E",
			Label:       "Total Race Population",
			Concept:     "RACE",
			Description: "Total population count for race estimates",
			Group:       "B02001",
		},
		"B19013_001E": {
			Name:        "B19013_001E",
			Label:       "Median Household Income",
			Concept:     "MEDIAN HOUSEHOLD INCOME IN THE PAST 12 MONTHS",
			Description: "Median household income in the past 12 months (in inflation-adjusted dollars)",
			Group:       "B19013",
		},
	}

	return variables, nil
}

// GetGeographyLevels возвращает доступные географические уровни для набора данных (тестовые данные)
func (m *MockCensusAPI) GetGeographyLevels(dataset, year string) ([]GeographyLevel, error) {
	// Тестовые данные о доступных географических уровнях
	geoLevels := []GeographyLevel{
		{
			Name:        "state",
			Description: "States and Equivalents",
			RequiredFor: []string{"county", "tract", "block"},
			Wildcards:   true,
		},
		{
			Name:        "county",
			Description: "Counties and Equivalents",
			RequiredFor: []string{"tract", "block"},
			Wildcards:   true,
		},
		{
			Name:        "tract",
			Description: "Census Tracts",
			RequiredFor: []string{"block"},
			Wildcards:   true,
		},
		{
			Name:        "block",
			Description: "Census Blocks",
			RequiredFor: []string{},
			Wildcards:   true,
		},
		{
			Name:        "us",
			Description: "United States",
			RequiredFor: []string{},
			Wildcards:   false,
		},
	}

	return geoLevels, nil
}

// GetCustomData позволяет запросить пользовательские данные (тестовые данные)
func (m *MockCensusAPI) GetCustomData(request CustomDataRequest) ([]map[string]string, error) {
	// Тестовые данные для пользовательского запроса
	data := []map[string]string{
		{
			"NAME":        "California",
			"B01001_001E": "39538223",
			"B19013_001E": "78672",
			"state":       "06",
		},
		{
			"NAME":        "New York",
			"B01001_001E": "20201249",
			"B19013_001E": "71117",
			"state":       "36",
		},
		{
			"NAME":        "Texas",
			"B01001_001E": "29145505",
			"B19013_001E": "63826",
			"state":       "48",
		},
	}

	// Если запрошены конкретные переменные, фильтруем данные
	if len(request.Variables) > 0 {
		// Всегда добавляем NAME и географические идентификаторы
		needVars := map[string]bool{"NAME": true}
		for _, geo := range []string{request.GeoLevel} {
			needVars[geo] = true
		}

		// Добавляем запрошенные переменные
		for _, v := range request.Variables {
			needVars[v] = true
		}

		// Фильтруем данные по запрошенным переменным
		for i, item := range data {
			filtered := make(map[string]string)
			for k, v := range item {
				if needVars[k] {
					filtered[k] = v
				}
			}
			data[i] = filtered
		}
	}

	return data, nil
}

// contains проверяет, содержит ли строка подстроку без учета регистра
func contains(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return indexOf(s, substr) >= 0
}

// toLower конвертирует строку в нижний регистр
func toLower(s string) string {
	result := ""
	for _, ch := range s {
		if ch >= 'A' && ch <= 'Z' {
			result += string(ch + ('a' - 'A'))
		} else {
			result += string(ch)
		}
	}
	return result
}

// indexOf находит индекс подстроки в строке
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
