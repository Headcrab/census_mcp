package main

import (
	"census_mcp/app"
	"flag"
	"log/slog"
	"os"
)

func main() {
	var transport string
	var testMode bool
	var apiKey string

	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.BoolVar(&testMode, "test", false, "Run in test mode")
	flag.StringVar(&apiKey, "k", "", "Census API key (if not provided, will use CENSUS_API_KEY env var)")
	flag.StringVar(&apiKey, "key", "", "Census API key (if not provided, will use CENSUS_API_KEY env var)")
	flag.Parse()

	// Настраиваем текстовый логгер с указанием времени
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

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

	// Создаем и запускаем сервер
	server, err := app.NewServer(config)
	if err != nil {
		slog.Error("Ошибка при создании сервера", "err", err)
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		slog.Error("Ошибка сервера", "err", err)
		os.Exit(1)
	}
}
