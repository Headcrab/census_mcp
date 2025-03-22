package census

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sort"
	"strings"
)

// Константы для ключей логирования
const (
	key_type       = "type"
	key_item_count = "item_count"
)

// Formatter определяет интерфейс для форматирования данных Census API
type Formatter interface {
	Format(ctx context.Context, data interface{}) string
}

// TextFormatter представляет текстовый форматтер для данных Census API
type TextFormatter struct{}

// NewTextFormatter создает новый экземпляр текстового форматтера
func NewTextFormatter() *TextFormatter {
	slog.DebugContext(context.Background(), "Создание нового текстового форматтера")
	return &TextFormatter{}
}

// Format форматирует данные Census API в текстовом виде
func (f *TextFormatter) Format(ctx context.Context, data interface{}) string {
	slog.DebugContext(ctx, "Форматирование данных",
		key_type, reflect.TypeOf(data))

	if data == nil {
		slog.WarnContext(ctx, "Попытка форматирования nil данных")
		return "Нет данных"
	}

	switch v := data.(type) {
	case []PopulationData:
		return f.formatPopulationData(ctx, v)
	case []DatasetInfo:
		return f.formatDatasetInfo(ctx, v)
	case map[string]VariableInfo:
		return f.formatVariableInfo(ctx, v)
	case []GeographyLevel:
		return f.formatGeographyLevel(ctx, v)
	case []map[string]string:
		return f.formatCustomData(ctx, v)
	default:
		slog.WarnContext(ctx, "Неизвестный тип данных для форматирования",
			key_type, reflect.TypeOf(data))
		return fmt.Sprintf("%v", data)
	}
}

// formatPopulationData форматирует данные о населении
func (f *TextFormatter) formatPopulationData(ctx context.Context, data []PopulationData) string {
	slog.DebugContext(ctx, "Форматирование данных о населении",
		key_item_count, len(data))

	if len(data) == 0 {
		return "Нет данных о населении"
	}

	var sb strings.Builder
	sb.WriteString("| Регион | Население |\n")
	sb.WriteString("|--------|----------|\n")

	for _, item := range data {
		// Определяем, это штат или округ
		var regionName string
		if item.County != "" {
			regionName = fmt.Sprintf("%s (округ %s, штат %s)", item.Name, item.County, item.State)
		} else {
			regionName = fmt.Sprintf("%s (штат %s)", item.Name, item.State)
		}

		sb.WriteString(fmt.Sprintf("| %s | %s |\n", regionName, item.Population))
	}

	slog.DebugContext(ctx, "Форматирование данных о населении завершено",
		key_item_count, len(data))
	return sb.String()
}

// formatDatasetInfo форматирует информацию о наборах данных
func (f *TextFormatter) formatDatasetInfo(ctx context.Context, data []DatasetInfo) string {
	slog.DebugContext(ctx, "Форматирование информации о наборах данных",
		key_item_count, len(data))

	if len(data) == 0 {
		return "Нет данных о наборах данных"
	}

	var sb strings.Builder
	sb.WriteString("# Доступные наборы данных\n\n")

	for _, item := range data {
		sb.WriteString(fmt.Sprintf("## %s\n", item.Title))
		sb.WriteString(fmt.Sprintf("- **ID набора**: %s\n", item.Dataset))
		sb.WriteString(fmt.Sprintf("- **Описание**: %s\n", item.Description))
		sb.WriteString("- **Доступные годы**: ")
		if len(item.YearsAvailable) > 0 {
			sb.WriteString(strings.Join(item.YearsAvailable, ", "))
		} else {
			sb.WriteString("Нет информации")
		}
		sb.WriteString("\n\n")
	}

	slog.DebugContext(ctx, "Форматирование информации о наборах данных завершено",
		key_item_count, len(data))
	return sb.String()
}

// formatVariableInfo форматирует информацию о переменных
func (f *TextFormatter) formatVariableInfo(ctx context.Context, data map[string]VariableInfo) string {
	slog.DebugContext(ctx, "Форматирование информации о переменных",
		key_item_count, len(data))

	if len(data) == 0 {
		return "Нет данных о переменных"
	}

	var sb strings.Builder
	sb.WriteString("# Доступные переменные\n\n")

	// Для стабильного вывода, сортируем ключи
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		item := data[key]
		sb.WriteString(fmt.Sprintf("## %s: %s\n", key, item.Label))
		if item.Description != "" {
			sb.WriteString(fmt.Sprintf("- **Описание**: %s\n", item.Description))
		}
		if item.Concept != "" {
			sb.WriteString(fmt.Sprintf("- **Концепция**: %s\n", item.Concept))
		}
		if item.Group != "" {
			sb.WriteString(fmt.Sprintf("- **Группа**: %s\n", item.Group))
		}
		sb.WriteString("\n")
	}

	slog.DebugContext(ctx, "Форматирование информации о переменных завершено",
		key_item_count, len(data))
	return sb.String()
}

// formatGeographyLevel форматирует информацию о географических уровнях
func (f *TextFormatter) formatGeographyLevel(ctx context.Context, data []GeographyLevel) string {
	slog.DebugContext(ctx, "Форматирование информации о географических уровнях",
		key_item_count, len(data))

	if len(data) == 0 {
		return "Нет данных о географических уровнях"
	}

	var sb strings.Builder
	sb.WriteString("# Доступные географические уровни\n\n")

	for _, item := range data {
		sb.WriteString(fmt.Sprintf("## %s\n", item.Name))
		sb.WriteString(fmt.Sprintf("- **Описание**: %s\n", item.Description))

		if len(item.RequiredFor) > 0 {
			sb.WriteString(fmt.Sprintf("- **Требуется для**: %s\n", strings.Join(item.RequiredFor, ", ")))
		}

		sb.WriteString(fmt.Sprintf("- **Поддержка подстановочных знаков**: %t\n", item.Wildcards))
		sb.WriteString("\n")
	}

	slog.DebugContext(ctx, "Форматирование информации о географических уровнях завершено",
		key_item_count, len(data))
	return sb.String()
}

// formatCustomData форматирует пользовательские данные
func (f *TextFormatter) formatCustomData(ctx context.Context, data []map[string]string) string {
	slog.DebugContext(ctx, "Форматирование пользовательских данных",
		key_item_count, len(data))

	if len(data) == 0 {
		return "Нет данных"
	}

	var sb strings.Builder

	// Получаем все возможные заголовки из всех записей
	var headers []string
	headerMap := make(map[string]bool)

	for _, item := range data {
		for key := range item {
			if !headerMap[key] {
				headerMap[key] = true
				headers = append(headers, key)
			}
		}
	}

	// Сортируем заголовки для стабильного вывода
	sort.Strings(headers)

	// Создаем заголовок таблицы
	sb.WriteString("| ")
	for _, header := range headers {
		sb.WriteString(header + " | ")
	}
	sb.WriteString("\n")

	// Создаем разделительную линию
	sb.WriteString("| ")
	for range headers {
		sb.WriteString("--- | ")
	}
	sb.WriteString("\n")

	// Добавляем строки данных
	for _, item := range data {
		sb.WriteString("| ")
		for _, header := range headers {
			value, ok := item[header]
			if !ok {
				value = "N/A"
			}
			sb.WriteString(value + " | ")
		}
		sb.WriteString("\n")
	}

	slog.DebugContext(ctx, "Форматирование пользовательских данных завершено",
		key_item_count, len(data))
	return sb.String()
}
