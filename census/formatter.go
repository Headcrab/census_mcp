package census

import (
	"fmt"
	"strings"
)

// Formatter определяет интерфейс для форматирования результатов
type Formatter interface {
	Format(data interface{}) string
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
	default:
		return fmt.Sprintf("Неподдерживаемый тип данных: %T", data)
	}
}

// NewTextFormatter создает новый форматтер в текстовом виде
func NewTextFormatter() Formatter {
	return &TextFormatter{}
}

// boolToYesNo преобразует булево значение в строку "Да" или "Нет"
func boolToYesNo(value bool) string {
	if value {
		return "Да"
	}
	return "Нет"
}
