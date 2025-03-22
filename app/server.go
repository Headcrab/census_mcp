package app

import (
	"census_mcp/census"
	"census_mcp/mcp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	mcpsdk "github.com/mark3labs/mcp-go/server"
)

// Константы для ключей логирования
const (
	key_test_mode   = "test_mode"
	key_transport   = "transport"
	key_err         = "err"
	key_duration    = "duration"
	key_count       = "count"
	key_displayed   = "displayed"
	key_search_term = "search_term"
	key_dataset     = "dataset"
	key_year        = "year"
	key_variables   = "variables"
	key_geo_level   = "geo_level"
	key_uptime      = "uptime"
)

// ServerConfig содержит конфигурацию сервера
type ServerConfig struct {
	Transport string
	TestMode  bool
	APIKey    string
}

// Server инкапсулирует логику запуска и настройки MCP сервера
type Server struct {
	config    ServerConfig
	mcpServer *mcpsdk.MCPServer
	tools     mcp.CensusToolHandler
	api       census.CensusAPIClient
}

// NewServer создает новый экземпляр сервера
func NewServer(config ServerConfig) (*Server, error) {
	slog.Debug("Создание нового сервера",
		key_test_mode, config.TestMode,
		key_transport, config.Transport)

	// Создаем форматтер
	formatter := census.NewTextFormatter()
	slog.Debug("Создан текстовый форматтер для данных Census API")

	var api census.CensusAPIClient
	var tools mcp.CensusToolHandler

	startTime := time.Now()

	// В тестовом режиме используем мок-клиент
	if config.TestMode {
		slog.Info("Инициализация тестового режима с мок-данными")
		mockAPI := census.NewMockCensusAPI()
		api = mockAPI
		tools = mcp.NewCensusToolHandler(mockAPI, formatter)
		slog.Info("Используется тестовый клиент Census API (мок-данные)")
	} else {
		// Создаем реальный клиент Census API
		slog.Info("Инициализация режима работы с реальным Census API")
		var censusAPI *census.CensusAPI
		var err error

		if config.APIKey != "" {
			slog.Debug("Использование ключа API из конфигурации")
			censusAPI = census.NewCensusAPI(config.APIKey)
		} else {
			slog.Debug("Попытка получить ключ API из переменной окружения")
			censusAPI, err = census.NewCensusAPIFromEnv()
			if err != nil {
				slog.Error("Не удалось создать клиент Census API",
					key_err, err)
				return nil, fmt.Errorf("ошибка при создании Census API клиента: %w", err)
			}
		}

		api = censusAPI
		tools = mcp.NewCensusToolHandler(censusAPI, formatter)
		slog.Info("Клиент Census API успешно инициализирован")
	}

	initDuration := time.Since(startTime)
	slog.Debug("Инициализация клиента API заняла",
		key_duration, initDuration)

	// Создаем MCP сервер
	slog.Debug("Создание MCP сервера")
	mcpServer := mcpsdk.NewMCPServer(
		"census-api",         // имя сервера
		"1.0.0",              // версия
		mcpsdk.WithLogging(), // включаем логирование
	)

	// Регистрируем инструменты
	slog.Debug("Регистрация инструментов Census API")
	mcp.RegisterCensusTools(mcpServer, tools)

	slog.Info("Инструменты Census API добавлены")

	totalInitTime := time.Since(startTime)
	slog.Debug("Сервер полностью инициализирован",
		key_duration, totalInitTime)

	return &Server{
		config:    config,
		mcpServer: mcpServer,
		tools:     tools,
		api:       api,
	}, nil
}

// RunTests запускает тестовые примеры
func (s *Server) RunTests() {
	slog.Info("Запуск тестовых примеров")
	fmt.Println("### ЗАПУСК ТЕСТОВЫХ ПРИМЕРОВ ###")

	// Получаем мок-API для тестов
	slog.Debug("Инициализация мок-API для тестов")
	mockAPI := census.NewMockCensusAPI()

	// Тестируем получение данных о населении штатов
	slog.Info("Тестирование получения данных о населении штатов")
	fmt.Println("=== Тестирование получения данных о населении штатов ===")
	states, err := mockAPI.GetStatePopulation("")
	if err != nil {
		slog.Error("Ошибка при тестировании получения данных о населении штатов",
			key_err, err)
		fmt.Printf("Ошибка: %v\n", err)
	} else {
		// Выводим только первые 3 результата для краткости
		formatter := census.NewTextFormatter()
		var limitedStates []census.PopulationData
		if len(states) > 3 {
			limitedStates = states[:3]
		} else {
			limitedStates = states
		}
		slog.Debug("Получены тестовые данные о населении штатов",
			key_count, len(states),
			key_displayed, len(limitedStates))
		fmt.Println(formatter.Format(context.Background(), limitedStates))
	}

	// Тестируем поиск штата по названию
	slog.Info("Тестирование поиска штата по названию")
	fmt.Println("=== Тестирование поиска штата по названию ===")
	searchTerm := "york"
	slog.Debug("Поиск штата по названию",
		key_search_term, searchTerm)
	searchResults, err := mockAPI.SearchStateByName(searchTerm)
	if err != nil {
		slog.Error("Ошибка при тестировании поиска штата по названию",
			key_err, err,
			key_search_term, searchTerm)
		fmt.Printf("Ошибка при поиске '%s': %v\n", searchTerm, err)
	} else {
		slog.Debug("Получены результаты поиска штата",
			key_search_term, searchTerm,
			key_count, len(searchResults))
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(context.Background(), searchResults))
	}

	// Тестируем получение доступных наборов данных
	slog.Info("Тестирование получения доступных наборов данных")
	fmt.Println("=== Тестирование получения доступных наборов данных ===")
	datasets, err := mockAPI.GetAvailableDatasets()
	if err != nil {
		slog.Error("Ошибка при тестировании получения доступных наборов данных",
			key_err, err)
		fmt.Printf("Ошибка при получении наборов данных: %v\n", err)
	} else {
		slog.Debug("Получены доступные наборы данных",
			key_count, len(datasets))
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(context.Background(), datasets))
	}

	// Тестируем получение переменных набора данных
	slog.Info("Тестирование получения переменных набора данных")
	fmt.Println("=== Тестирование получения переменных набора данных ===")
	datasetName := "acs/acs1"
	year := "2021"
	slog.Debug("Получение переменных набора данных",
		key_dataset, datasetName,
		key_year, year)
	variables, err := mockAPI.GetVariables(datasetName, year)
	if err != nil {
		slog.Error("Ошибка при тестировании получения переменных набора данных",
			key_err, err,
			key_dataset, datasetName,
			key_year, year)
		fmt.Printf("Ошибка при получении переменных: %v\n", err)
	} else {
		slog.Debug("Получены переменные набора данных",
			key_count, len(variables),
			key_dataset, datasetName,
			key_year, year)
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(context.Background(), variables))
	}

	// Тестируем получение географических уровней
	slog.Info("Тестирование получения географических уровней")
	fmt.Println("=== Тестирование получения географических уровней ===")
	geoLevels, err := mockAPI.GetGeographyLevels(datasetName, year)
	if err != nil {
		slog.Error("Ошибка при тестировании получения географических уровней",
			key_err, err,
			key_dataset, datasetName,
			key_year, year)
		fmt.Printf("Ошибка при получении географических уровней: %v\n", err)
	} else {
		slog.Debug("Получены географические уровни",
			key_count, len(geoLevels),
			key_dataset, datasetName,
			key_year, year)
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(context.Background(), geoLevels))
	}

	// Тестируем получение пользовательских данных
	slog.Info("Тестирование получения пользовательских данных")
	fmt.Println("=== Тестирование получения пользовательских данных ===")
	customRequest := census.CustomDataRequest{
		Variables: []string{"NAME", "B01001_001E", "B19013_001E"},
		Dataset:   "acs/acs1",
		Year:      "2021",
		GeoLevel:  "state",
		GeoFilter: map[string]string{"state": "*"},
	}
	slog.Debug("Запрос пользовательских данных",
		key_variables, customRequest.Variables,
		key_dataset, customRequest.Dataset,
		key_year, customRequest.Year,
		key_geo_level, customRequest.GeoLevel)
	customData, err := mockAPI.GetCustomData(customRequest)
	if err != nil {
		slog.Error("Ошибка при тестировании получения пользовательских данных",
			key_err, err)
		fmt.Printf("Ошибка при получении пользовательских данных: %v\n", err)
	} else {
		slog.Debug("Получены пользовательские данные",
			key_count, len(customData))
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(context.Background(), customData))
	}

	// Информация по использованию через MCP клиент
	slog.Info("Завершение тестирования, вывод примеров запросов через MCP")
	fmt.Println("\n=== Тестирование через MCP сервер ===")
	fmt.Println("Для тестирования через MCP клиент можно использовать запросы:")
	fmt.Println(`1. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_state_population","arguments":{}}}`)
	fmt.Println(`2. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_state_population","arguments":{"stateID":"06"}}}`)
	fmt.Println(`3. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"search_state_by_name","arguments":{"name":"california"}}}`)
	fmt.Println(`4. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_available_datasets","arguments":{}}}`)
	fmt.Println(`5. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_variables","arguments":{"dataset":"acs/acs1","year":"2021"}}}`)
	fmt.Println(`6. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_geography_levels","arguments":{"dataset":"acs/acs1","year":"2021"}}}`)
	fmt.Println(`7. {"jsonrpc":"2.0","id":"test","method":"mcp.call","params":{"tool":"get_custom_data","arguments":{"dataset":"acs/acs1","year":"2021","geoLevel":"state","variables":["NAME","B01001_001E"]}}}`)
	fmt.Println("\nЗапустите сервер без флага -test и отправьте запрос через клиент MCP")
}

// Start запускает сервер
func (s *Server) Start() error {
	if s.config.TestMode {
		slog.Info("Запуск сервера в тестовом режиме")
		s.RunTests()
		return nil
	}

	startTime := time.Now()
	slog.Info("Запуск сервера Census MCP API",
		key_transport, s.config.Transport)

	if s.config.Transport == "sse" {
		slog.Info("Инициализация SSE сервера на порту 8080")
		sseServer := mcpsdk.NewSSEServer(s.mcpServer, mcpsdk.WithBaseURL("http://localhost:8080"))

		// Настраиваем обработчик для /health
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{"status":"ok"}`))
			if err != nil {
				slog.Error("Ошибка при записи HTTP ответа",
					key_err, err)
			}
		})

		slog.Info("SSE server listening on :8080")
		if err := sseServer.Start(":8080"); err != nil {
			slog.Error("Ошибка при запуске SSE сервера",
				key_err, err,
				key_uptime, time.Since(startTime))
			return fmt.Errorf("ошибка запуска SSE сервера: %w", err)
		}
	} else {
		slog.Info("Запуск Census API сервера через stdio")
		if err := mcpsdk.ServeStdio(s.mcpServer); err != nil {
			slog.Error("Ошибка при запуске stdio сервера",
				key_err, err,
				key_uptime, time.Since(startTime))
			return fmt.Errorf("ошибка запуска stdio сервера: %w", err)
		}
	}

	slog.Info("Сервер завершил работу",
		key_uptime, time.Since(startTime))
	return nil
}
