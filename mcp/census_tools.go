package mcp

import (
	"census_mcp/census"
	"context"
	"log/slog"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Константы для ключей логирования
const (
	key_state_id      = "state_id"
	key_err           = "err"
	key_count         = "count"
	key_name          = "name"
	key_search_name   = "search_name"
	key_dataset       = "dataset"
	key_year          = "year"
	key_geo_level     = "geo_level"
	key_valid         = "valid"
	key_dataset_valid = "dataset_valid"
	key_year_valid    = "year_valid"
)

// CensusToolHandler определяет интерфейс для обработчика инструментов Census MCP
type CensusToolHandler interface {
	// HandleGetStatePopulationTool обрабатывает запрос на получение данных о населении штата
	HandleGetStatePopulationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleGetCountyPopulationTool обрабатывает запрос на получение данных о населении округа
	HandleGetCountyPopulationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleSearchStateByNameTool обрабатывает запрос на поиск штата по названию
	HandleSearchStateByNameTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleGetAvailableDatasetsTool обрабатывает запрос на получение доступных наборов данных
	HandleGetAvailableDatasetsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleGetVariablesTool обрабатывает запрос на получение доступных переменных набора данных
	HandleGetVariablesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleGetGeographyLevelsTool обрабатывает запрос на получение доступных географических уровней
	HandleGetGeographyLevelsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	// HandleGetCustomDataTool обрабатывает запрос на получение пользовательских данных
	HandleGetCustomDataTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// CensusDefaultToolHandler - стандартная реализация обработчика инструментов
type CensusDefaultToolHandler struct {
	api       census.CensusAPIClient
	formatter census.Formatter
}

// NewCensusToolHandler создает новый экземпляр обработчика инструментов
func NewCensusToolHandler(api census.CensusAPIClient, formatter census.Formatter) CensusToolHandler {
	return &CensusDefaultToolHandler{
		api:       api,
		formatter: formatter,
	}
}

// HandleGetStatePopulationTool обрабатывает запрос на получение данных о населении штата
func (h *CensusDefaultToolHandler) HandleGetStatePopulationTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения данных о населении штата")

	arguments := request.Params.Arguments
	stateID, _ := arguments["stateID"].(string)

	slog.DebugContext(ctx, "Параметры инструмента получения данных о населении штата",
		key_state_id, stateID)

	// Получение данных о населении штата
	population, err := h.api.GetStatePopulation(stateID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении данных о населении штата",
			key_err, err,
			key_state_id, stateID)
		return mcp.NewToolResultError("Ошибка при получении данных о населении: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Получены данные о населении штата",
		key_count, len(population),
		key_state_id, stateID)

	// Форматирование результатов
	result := h.formatter.Format(ctx, population)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleGetCountyPopulationTool обрабатывает запрос на получение данных о населении округа
func (h *CensusDefaultToolHandler) HandleGetCountyPopulationTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения данных о населении округов")

	arguments := request.Params.Arguments
	stateID, _ := arguments["stateID"].(string)

	slog.DebugContext(ctx, "Параметры инструмента получения данных о населении округов",
		key_state_id, stateID)

	// Получение данных о населении округов
	population, err := h.api.GetCountyPopulation(stateID)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении данных о населении округов",
			key_err, err,
			key_state_id, stateID)
		return mcp.NewToolResultError("Ошибка при получении данных о населении округов: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Получены данные о населении округов",
		key_count, len(population),
		key_state_id, stateID)

	// Форматирование результатов
	result := h.formatter.Format(ctx, population)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleSearchStateByNameTool обрабатывает запрос на поиск штата по названию
func (h *CensusDefaultToolHandler) HandleSearchStateByNameTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента поиска штата по названию")

	arguments := request.Params.Arguments
	name, ok := arguments["name"].(string)

	slog.DebugContext(ctx, "Параметры инструмента поиска штата по названию",
		key_name, name,
		key_valid, ok)

	if !ok || name == "" {
		slog.ErrorContext(ctx, "Отсутствует обязательный параметр name для поиска штата")
		return mcp.NewToolResultError("Необходимо указать параметр 'name'"), nil
	}

	// Поиск штата по названию
	states, err := h.api.SearchStateByName(name)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при поиске штата по названию",
			key_err, err,
			key_search_name, name)
		return mcp.NewToolResultError("Ошибка при поиске штата: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Результаты поиска штата по названию",
		key_count, len(states),
		key_search_name, name)

	if len(states) == 0 {
		slog.InfoContext(ctx, "Штаты не найдены по запросу",
			key_name, name)
		return mcp.NewToolResultText("Штаты не найдены по запросу: " + name), nil
	}

	// Форматирование результатов
	result := h.formatter.Format(ctx, states)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleGetAvailableDatasetsTool обрабатывает запрос на получение доступных наборов данных
func (h *CensusDefaultToolHandler) HandleGetAvailableDatasetsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения доступных наборов данных")

	// Получение данных о доступных наборах данных
	datasets, err := h.api.GetAvailableDatasets()
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении данных о доступных наборах данных",
			key_err, err)
		return mcp.NewToolResultError("Ошибка при получении данных о доступных наборах данных: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Получены данные о доступных наборах данных",
		key_count, len(datasets))

	// Форматирование результатов
	result := h.formatter.Format(ctx, datasets)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleGetVariablesTool обрабатывает запрос на получение доступных переменных набора данных
func (h *CensusDefaultToolHandler) HandleGetVariablesTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения переменных набора данных")

	arguments := request.Params.Arguments
	dataset, ok1 := arguments["dataset"].(string)
	year, ok2 := arguments["year"].(string)

	slog.DebugContext(ctx, "Параметры инструмента получения переменных набора данных",
		key_dataset, dataset,
		key_year, year,
		key_dataset_valid, ok1,
		key_year_valid, ok2)

	if !ok1 || !ok2 || dataset == "" || year == "" {
		slog.ErrorContext(ctx, "Отсутствуют обязательные параметры для получения переменных набора данных")
		return mcp.NewToolResultError("Необходимо указать параметры 'dataset' и 'year'"), nil
	}

	// Получение данных о доступных переменных
	variables, err := h.api.GetVariables(dataset, year)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении данных о доступных переменных",
			key_err, err,
			key_dataset, dataset,
			key_year, year)
		return mcp.NewToolResultError("Ошибка при получении данных о доступных переменных: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Получены данные о доступных переменных",
		key_count, len(variables),
		key_dataset, dataset,
		key_year, year)

	// Форматирование результатов
	result := h.formatter.Format(ctx, variables)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleGetGeographyLevelsTool обрабатывает запрос на получение доступных географических уровней
func (h *CensusDefaultToolHandler) HandleGetGeographyLevelsTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения географических уровней")

	arguments := request.Params.Arguments
	dataset, ok1 := arguments["dataset"].(string)
	year, ok2 := arguments["year"].(string)

	slog.DebugContext(ctx, "Параметры инструмента получения географических уровней",
		key_dataset, dataset,
		key_year, year,
		key_dataset_valid, ok1,
		key_year_valid, ok2)

	if !ok1 || !ok2 || dataset == "" || year == "" {
		slog.ErrorContext(ctx, "Отсутствуют обязательные параметры для получения географических уровней")
		return mcp.NewToolResultError("Необходимо указать параметры 'dataset' и 'year'"), nil
	}

	// Получение данных о доступных географических уровнях
	levels, err := h.api.GetGeographyLevels(dataset, year)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении данных о доступных географических уровнях",
			key_err, err,
			key_dataset, dataset,
			key_year, year)
		return mcp.NewToolResultError("Ошибка при получении данных о доступных географических уровнях: " + err.Error()), nil
	}

	slog.DebugContext(ctx, "Получены данные о доступных географических уровнях",
		key_count, len(levels),
		key_dataset, dataset,
		key_year, year)

	// Форматирование результатов
	result := h.formatter.Format(ctx, levels)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// HandleGetCustomDataTool обрабатывает запрос на получение пользовательских данных
func (h *CensusDefaultToolHandler) HandleGetCustomDataTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	slog.InfoContext(ctx, "Обработка инструмента получения пользовательских данных")

	arguments := request.Params.Arguments

	// Извлечение обязательных параметров
	dataset, ok1 := arguments["dataset"].(string)
	year, ok2 := arguments["year"].(string)
	geoLevel, ok3 := arguments["geoLevel"].(string)
	variables, ok4 := arguments["variables"].([]interface{})

	if !ok1 || !ok2 || !ok3 || !ok4 || dataset == "" || year == "" || geoLevel == "" || len(variables) == 0 {
		return mcp.NewToolResultError("Необходимо указать параметры 'dataset', 'year', 'geoLevel' и 'variables'"), nil
	}

	// Преобразование списка переменных
	varList := make([]string, 0, len(variables))
	for _, v := range variables {
		if vs, ok := v.(string); ok {
			varList = append(varList, vs)
		}
	}

	// Извлечение и преобразование гео-фильтров
	geoFilterMap := make(map[string]string)
	if geoFilter, ok := arguments["geoFilter"].(map[string]interface{}); ok {
		for k, v := range geoFilter {
			if vs, ok := v.(string); ok {
				geoFilterMap[k] = vs
			}
		}
	} else {
		// Если не указан фильтр, добавляем wildcard для текущего уровня
		geoFilterMap[geoLevel] = "*"
	}

	// Создание запроса пользовательских данных
	customRequest := census.CustomDataRequest{
		Variables: varList,
		Dataset:   dataset,
		Year:      year,
		GeoLevel:  geoLevel,
		GeoFilter: geoFilterMap,
	}

	// Получение пользовательских данных
	customData, err := h.api.GetCustomData(customRequest)
	if err != nil {
		return mcp.NewToolResultError("Ошибка при получении пользовательских данных: " + err.Error()), nil
	}

	// Форматирование результатов
	result := h.formatter.Format(ctx, customData)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// RegisterCensusTools регистрирует инструменты Census MCP
func RegisterCensusTools(mcpServer *server.MCPServer, handler CensusToolHandler) {
	// Функция RegisterCensusTools не имеет контекста в параметрах,
	// поэтому используем background контекст
	ctx := context.Background()
	slog.InfoContext(ctx, "Регистрация инструментов Census API")

	// Инструмент для получения данных о населении штата
	mcpServer.AddTool(mcp.NewTool("get_state_population",
		mcp.WithDescription("Получает данные о населении штата США"),
		mcp.WithString("stateID",
			mcp.Description("ID штата (например, '06' для Калифорнии). Если не указан, возвращает данные для всех штатов"),
		),
	), handler.HandleGetStatePopulationTool)

	// Инструмент для получения данных о населении округа
	mcpServer.AddTool(mcp.NewTool("get_county_population",
		mcp.WithDescription("Получает данные о населении округов в штате США"),
		mcp.WithString("stateID",
			mcp.Description("ID штата (например, '06' для Калифорнии). Если не указан, возвращает данные для всех округов"),
		),
	), handler.HandleGetCountyPopulationTool)

	// Инструмент для поиска штата по названию
	mcpServer.AddTool(mcp.NewTool("search_state_by_name",
		mcp.WithDescription("Ищет штат по названию (полному или частичному)"),
		mcp.WithString("name",
			mcp.Description("Название штата для поиска"),
			mcp.Required(),
		),
	), handler.HandleSearchStateByNameTool)

	// Инструмент для получения доступных наборов данных
	mcpServer.AddTool(mcp.NewTool("get_available_datasets",
		mcp.WithDescription("Получает список доступных наборов данных Census API"),
	), handler.HandleGetAvailableDatasetsTool)

	// Инструмент для получения переменных набора данных
	mcpServer.AddTool(mcp.NewTool("get_variables",
		mcp.WithDescription("Получает список доступных переменных для указанного набора данных"),
		mcp.WithString("dataset",
			mcp.Description("Набор данных (например, 'acs/acs1')"),
			mcp.Required(),
		),
		mcp.WithString("year",
			mcp.Description("Год данных (например, '2021')"),
			mcp.Required(),
		),
	), handler.HandleGetVariablesTool)

	// Инструмент для получения географических уровней
	mcpServer.AddTool(mcp.NewTool("get_geography_levels",
		mcp.WithDescription("Получает список доступных географических уровней для указанного набора данных"),
		mcp.WithString("dataset",
			mcp.Description("Набор данных (например, 'acs/acs1')"),
			mcp.Required(),
		),
		mcp.WithString("year",
			mcp.Description("Год данных (например, '2021')"),
			mcp.Required(),
		),
	), handler.HandleGetGeographyLevelsTool)

	// Инструмент для получения пользовательских данных
	mcpServer.AddTool(mcp.NewTool("get_custom_data",
		mcp.WithDescription("Позволяет делать пользовательские запросы к Census API с указанием набора данных, года, переменных и географического уровня"),
		mcp.WithString("dataset",
			mcp.Description("Набор данных (например, 'acs/acs1')"),
			mcp.Required(),
		),
		mcp.WithString("year",
			mcp.Description("Год данных (например, '2021')"),
			mcp.Required(),
		),
		mcp.WithString("geoLevel",
			mcp.Description("Географический уровень (например, 'state' или 'county')"),
			mcp.Required(),
		),
		mcp.WithArray("variables",
			mcp.Description("Список переменных для запроса (например, ['NAME', 'B01001_001E'])"),
			mcp.Required(),
		),
		mcp.WithObject("geoFilter",
			mcp.Description("Фильтр географии (например, {\"state\": \"06\", \"county\": \"*\"}). Если не указан, будет использован wildcard для указанного географического уровня"),
		),
	), handler.HandleGetCustomDataTool)
}
