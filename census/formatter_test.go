package census

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatter_Format_PopulationData(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные о населении
	popData := []PopulationData{
		{
			Name:       "California",
			Population: "39538223",
			State:      "06",
		},
		{
			Name:       "Los Angeles County",
			Population: "10014009",
			State:      "06",
			County:     "037",
		},
	}

	// Форматируем данные
	result := formatter.Format(popData)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Результаты запроса к Census API:",
		"Название: California",
		"Население: 39538223",
		"Код штата: 06",
		"Название: Los Angeles County",
		"Население: 10014009",
		"Код штата: 06",
		"Код округа: 037",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_DatasetInfo(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные о наборах данных
	datasets := []DatasetInfo{
		{
			Title:          "American Community Survey 1-Year Estimates",
			Description:    "Annual survey covering demographic data",
			Dataset:        "acs/acs1",
			YearsAvailable: []string{"2019", "2020", "2021"},
		},
		{
			Title:          "Decennial Census",
			Description:    "Complete count of the US population",
			Dataset:        "dec/sf1",
			YearsAvailable: []string{"2000", "2010", "2020"},
		},
	}

	// Форматируем данные
	result := formatter.Format(datasets)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Доступные наборы данных Census API:",
		"Название: American Community Survey 1-Year Estimates",
		"Описание: Annual survey covering demographic data",
		"Идентификатор: acs/acs1",
		"Доступные годы: 2019, 2020, 2021",
		"Название: Decennial Census",
		"Описание: Complete count of the US population",
		"Идентификатор: dec/sf1",
		"Доступные годы: 2000, 2010, 2020",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_VariableInfo(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные о переменных
	variables := map[string]VariableInfo{
		"B01001_001E": {
			Name:        "B01001_001E",
			Label:       "Total Population",
			Description: "Total population count",
			Concept:     "SEX BY AGE",
			Group:       "B01001",
		},
		"NAME": {
			Name:        "NAME",
			Label:       "Geographic Area Name",
			Description: "Name of the geographic area",
		},
	}

	// Форматируем данные
	result := formatter.Format(variables)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Доступные переменные Census API:",
		"Переменная: B01001_001E",
		"Название: Total Population",
		"Описание: Total population count",
		"Концепция: SEX BY AGE",
		"Группа: B01001",
		"Переменная: NAME",
		"Название: Geographic Area Name",
		"Описание: Name of the geographic area",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_GeographyLevel(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные о географических уровнях
	levels := []GeographyLevel{
		{
			Name:        "state",
			Description: "States and Equivalent",
			Wildcards:   true,
		},
		{
			Name:        "county",
			Description: "Counties and Equivalent",
			RequiredFor: []string{"state"},
			Wildcards:   true,
		},
	}

	// Форматируем данные
	result := formatter.Format(levels)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Доступные географические уровни Census API:",
		"Уровень: state",
		"Описание: States and Equivalent",
		"Поддержка wildcard: Да",
		"Уровень: county",
		"Описание: Counties and Equivalent",
		"Требуется указать: state",
		"Поддержка wildcard: Да",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_CustomData(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые пользовательские данные
	customData := []map[string]string{
		{
			"NAME":        "California",
			"B01001_001E": "39538223",
			"state":       "06",
		},
		{
			"NAME":        "Texas",
			"B01001_001E": "29145505",
			"state":       "48",
		},
	}

	// Форматируем данные
	result := formatter.Format(customData)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Результаты пользовательского запроса к Census API:",
		"NAME: California",
		"B01001_001E: 39538223",
		"state: 06",
		"NAME: Texas",
		"B01001_001E: 29145505",
		"state: 48",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_UnsupportedType(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные неподдерживаемого типа
	unsupportedData := 123

	// Форматируем данные
	result := formatter.Format(unsupportedData)

	// Проверяем, что результат содержит сообщение об ошибке
	assert.True(t, strings.Contains(result, "Неподдерживаемый тип данных") ||
		strings.Contains(result, "не поддерживается"),
		"Результат должен содержать сообщение о неподдерживаемом типе")
}

func TestNewTextFormatter(t *testing.T) {
	formatter := NewTextFormatter()
	assert.NotNil(t, formatter)
	_, ok := formatter.(*TextFormatter)
	assert.True(t, ok)
}
