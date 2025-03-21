package app

import (
	"census_mcp/census"
	"census_mcp/mcp"
	"fmt"
	"log/slog"
	"os"

	mcpsdk "github.com/mark3labs/mcp-go/server"
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
	// Настраиваем логгер
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	// Создаем форматтер
	formatter := census.NewTextFormatter()

	var api census.CensusAPIClient
	var tools mcp.CensusToolHandler

	// В тестовом режиме используем мок-клиент
	if config.TestMode {
		mockAPI := census.NewMockCensusAPI()
		api = mockAPI
		tools = mcp.NewCensusToolHandler(mockAPI, formatter)
		slog.Info("Используется тестовый клиент Census API (мок-данные)")
	} else {
		// Создаем реальный клиент Census API
		var censusAPI *census.CensusAPI
		var err error

		if config.APIKey != "" {
			censusAPI = census.NewCensusAPI(config.APIKey)
		} else {
			censusAPI, err = census.NewCensusAPIFromEnv()
			if err != nil {
				return nil, fmt.Errorf("ошибка при создании Census API клиента: %w", err)
			}
		}

		api = censusAPI
		tools = mcp.NewCensusToolHandler(censusAPI, formatter)
	}

	// Создаем MCP сервер
	mcpServer := mcpsdk.NewMCPServer(
		"census-api",         // имя сервера
		"1.0.0",              // версия
		mcpsdk.WithLogging(), // включаем логирование
	)

	// Регистрируем инструменты
	mcp.RegisterCensusTools(mcpServer, tools)

	slog.Info("Инструменты Census API добавлены")

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
	mockAPI := census.NewMockCensusAPI()

	// Тестируем получение данных о населении штатов
	fmt.Println("=== Тестирование получения данных о населении штатов ===")
	states, err := mockAPI.GetStatePopulation("")
	if err != nil {
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
		fmt.Println(formatter.Format(limitedStates))
	}

	// Тестируем поиск штата по названию
	fmt.Println("=== Тестирование поиска штата по названию ===")
	searchTerm := "york"
	searchResults, err := mockAPI.SearchStateByName(searchTerm)
	if err != nil {
		fmt.Printf("Ошибка при поиске '%s': %v\n", searchTerm, err)
	} else {
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(searchResults))
	}

	// Тестируем получение доступных наборов данных
	fmt.Println("=== Тестирование получения доступных наборов данных ===")
	datasets, err := mockAPI.GetAvailableDatasets()
	if err != nil {
		fmt.Printf("Ошибка при получении наборов данных: %v\n", err)
	} else {
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(datasets))
	}

	// Тестируем получение переменных набора данных
	fmt.Println("=== Тестирование получения переменных набора данных ===")
	variables, err := mockAPI.GetVariables("acs/acs1", "2021")
	if err != nil {
		fmt.Printf("Ошибка при получении переменных: %v\n", err)
	} else {
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(variables))
	}

	// Тестируем получение географических уровней
	fmt.Println("=== Тестирование получения географических уровней ===")
	geoLevels, err := mockAPI.GetGeographyLevels("acs/acs1", "2021")
	if err != nil {
		fmt.Printf("Ошибка при получении географических уровней: %v\n", err)
	} else {
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(geoLevels))
	}

	// Тестируем получение пользовательских данных
	fmt.Println("=== Тестирование получения пользовательских данных ===")
	customRequest := census.CustomDataRequest{
		Variables: []string{"NAME", "B01001_001E", "B19013_001E"},
		Dataset:   "acs/acs1",
		Year:      "2021",
		GeoLevel:  "state",
		GeoFilter: map[string]string{"state": "*"},
	}
	customData, err := mockAPI.GetCustomData(customRequest)
	if err != nil {
		fmt.Printf("Ошибка при получении пользовательских данных: %v\n", err)
	} else {
		formatter := census.NewTextFormatter()
		fmt.Println(formatter.Format(customData))
	}

	// Информация по использованию через MCP клиент
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
		s.RunTests()
		return nil
	}

	if s.config.Transport == "sse" {
		sseServer := mcpsdk.NewSSEServer(s.mcpServer, mcpsdk.WithBaseURL("http://localhost:8080"))
		slog.Info("SSE server listening on :8080")
		if err := sseServer.Start(":8080"); err != nil {
			return fmt.Errorf("ошибка запуска SSE сервера: %w", err)
		}
	} else {
		slog.Info("Запуск Census API сервера через stdio")
		if err := mcpsdk.ServeStdio(s.mcpServer); err != nil {
			return fmt.Errorf("ошибка запуска stdio сервера: %w", err)
		}
	}

	return nil
}
