package census

import (
	"context"
	"fmt"
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
	result := formatter.Format(context.Background(), popData)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"Регион",
		"Население",
		"California (штат 06)",
		"39538223",
		"Los Angeles County (округ 037, штат 06)",
		"10014009",
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
	result := formatter.Format(context.Background(), datasets)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"# Доступные наборы данных",
		"## American Community Survey 1-Year Estimates",
		"**ID набора**: acs/acs1",
		"**Описание**: Annual survey covering demographic data",
		"**Доступные годы**: 2019, 2020, 2021",
		"## Decennial Census",
		"**ID набора**: dec/sf1",
		"**Описание**: Complete count of the US population",
		"**Доступные годы**: 2000, 2010, 2020",
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
	result := formatter.Format(context.Background(), variables)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"# Доступные переменные",
		"## B01001_001E: Total Population",
		"**Описание**: Total population count",
		"**Концепция**: SEX BY AGE",
		"**Группа**: B01001",
		"## NAME: Geographic Area Name",
		"**Описание**: Name of the geographic area",
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
	result := formatter.Format(context.Background(), levels)

	// Проверяем, что результат содержит ожидаемые строки
	expectedStrings := []string{
		"# Доступные географические уровни",
		"## state",
		"**Описание**: States and Equivalent",
		"**Поддержка подстановочных знаков**: true",
		"## county",
		"**Описание**: Counties and Equivalent",
		"**Требуется для**: state",
		"**Поддержка подстановочных знаков**: true",
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
	result := formatter.Format(context.Background(), customData)

	// Проверяем наличие заголовков таблицы и данных
	expectedStrings := []string{
		"NAME",
		"B01001_001E",
		"state",
		"California",
		"39538223",
		"06",
		"Texas",
		"29145505",
		"48",
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
	result := formatter.Format(context.Background(), unsupportedData)

	// Проверяем, что результат просто содержит строковое представление числа
	assert.Equal(t, "123", result)
}

func TestNewTextFormatter(t *testing.T) {
	formatter := NewTextFormatter()
	assert.NotNil(t, formatter)
	assert.Equal(t, "*census.TextFormatter", fmt.Sprintf("%T", formatter))
}
