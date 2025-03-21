package census

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextFormatter_CountLetters(t *testing.T) {
	formatter := NewTextFormatter()

	tests := []struct {
		name     string
		word     string
		letters  string
		expected map[rune]int
	}{
		{
			name:    "Базовый подсчет",
			word:    "привет",
			letters: "пр",
			expected: map[rune]int{
				'п': 1,
				'р': 1,
			},
		},
		{
			name:    "Повторяющиеся буквы",
			word:    "тестирование",
			letters: "те",
			expected: map[rune]int{
				'т': 2,
				'е': 2,
			},
		},
		{
			name:    "Пустое слово",
			word:    "",
			letters: "абв",
			expected: map[rune]int{
				'а': 0,
				'б': 0,
				'в': 0,
			},
		},
		{
			name:     "Пустые буквы для поиска",
			word:     "привет",
			letters:  "",
			expected: map[rune]int{},
		},
		{
			name:    "Буквы отсутствуют в слове",
			word:    "привет",
			letters: "xyz",
			expected: map[rune]int{
				'x': 0,
				'y': 0,
				'z': 0,
			},
		},
		{
			name:    "Регистрозависимый поиск",
			word:    "Привет",
			letters: "п",
			expected: map[rune]int{
				'п': 0, // 'П' и 'п' - разные символы
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatter.CountLetters(tt.word, tt.letters)
			assert.Equal(t, tt.expected, result)
		})
	}
}

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
		"Запись 1:",
		"NAME: California",
		"B01001_001E: 39538223",
		"state: 06",
		"Запись 2:",
		"NAME: Texas",
		"B01001_001E: 29145505",
		"state: 48",
	}

	for _, str := range expectedStrings {
		assert.Contains(t, result, str)
	}
}

func TestTextFormatter_Format_LetterCount(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные подсчета букв
	letterCount := map[rune]int{
		'а': 2,
		'б': 0,
		'в': 1,
	}

	// Форматируем данные
	result := formatter.Format(letterCount)

	// Проверяем, что результат содержит все буквы и их количество
	assert.Contains(t, result, "'а': 2")
	assert.Contains(t, result, "'б': 0")
	assert.Contains(t, result, "'в': 1")
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
