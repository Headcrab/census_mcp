package main

import (
	"census_mcp/app"
	"census_mcp/logger"
	"flag"
	"log/slog"
	"os"
)

// Константы для ключей логирования
const (
	key_transport = "transport"
	key_test_mode = "test_mode"
	key_log_level = "log_level"
	key_config    = "config"
	key_err       = "err"
)

func main() {
	var transport string
	var testMode bool
	var apiKey string
	var logLevelFlag string

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.BoolVar(&testMode, "test", false, "Run in test mode")
	flag.StringVar(&apiKey, "k", "", "Census API key (if not provided, will use CENSUS_API_KEY env var)")
	flag.StringVar(&apiKey, "key", "", "Census API key (if not provided, will use CENSUS_API_KEY env var)")
	flag.StringVar(&logLevelFlag, "log-level", "", "Log level (debug, info, warn, error)")
	flag.Parse()

	// Настраиваем логирование
	logLevel := logger.GetLogLevelFromEnv(logLevelFlag)
	logFile := logger.GetLogFileFromEnv()

	err := logger.SetupLogger(logger.Config{
		Level:    logLevel,
		FilePath: logFile,
	})

	if err != nil {
		// Не можем использовать slog, так как он еще не настроен
		panic("Ошибка настройки логгера: " + err.Error())
	}

	slog.Info("Запуск Census MCP API",
		key_transport, transport,
		key_test_mode, testMode,
		key_log_level, logLevel)

	// Информация о доступных возможностях API
	slog.Info("Census MCP API поддерживает следующие инструменты:")
	slog.Info("- get_state_population: получение данных о населении штатов")
	slog.Info("- get_county_population: получение данных о населении округов")
	slog.Info("- search_state_by_name: поиск штата по названию")
	slog.Info("- get_available_datasets: получение списка доступных наборов данных")
	slog.Info("- get_variables: получение списка переменных для набора данных")
	slog.Info("- get_geography_levels: получение доступных географических уровней")
	slog.Info("- get_custom_data: выполнение пользовательских запросов к Census API")

	// Конфигурация сервера
	config := app.ServerConfig{
		Transport: transport,
		TestMode:  testMode,
		APIKey:    apiKey,
	}

	slog.Debug("Создание сервера с конфигурацией",
		key_config, config)

	// Создаем и запускаем сервер
	server, err := app.NewServer(config)
	if err != nil {
		slog.Error("Ошибка при создании сервера",
			key_err, err)
		os.Exit(1)
	}

	slog.Info("Сервер успешно создан, запускаем...")

	if err := server.Start(); err != nil {
		slog.Error("Ошибка сервера",
			key_err, err)
		os.Exit(1)
	}

	slog.Info("Сервер завершил работу")
}
