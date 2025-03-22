package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"
)

// Константы для ключей логирования
const (
	key_level  = "level"
	key_output = "output"
)

// Config содержит настройки логгера
type Config struct {
	// Level уровень логирования (debug, info, warn, error)
	Level string
	// FilePath путь к файлу для сохранения логов (если пустой, лог выводится в stdout)
	FilePath string
}

// SetupLogger настраивает глобальный логгер приложения
func SetupLogger(config Config) error {
	// Определяем уровень логирования
	logLevel := slog.LevelInfo

	switch strings.ToLower(config.Level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	}

	// Определяем вывод для логов
	var output io.Writer = os.Stdout

	// Если указан путь к файлу, настраиваем вывод в файл
	if config.FilePath != "" {
		// Создаем директории для лога, если их нет
		dir := config.FilePath[:strings.LastIndex(config.FilePath, "/")]
		if dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("ошибка создания директории для логов: %w", err)
			}
		}

		// Открываем файл для записи логов (с добавлением)
		file, err := os.OpenFile(config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return fmt.Errorf("ошибка открытия файла логов: %w", err)
		}

		// Если нужно, можно настроить вывод как в файл, так и в консоль
		// output = io.MultiWriter(os.Stdout, file)
		output = file

		// Логгируем информацию о запуске
		fmt.Printf("Логи будут сохраняться в файл: %s\n", config.FilePath)
	}

	// Настраиваем текстовый логгер с указанием времени и настроенным уровнем
	handler := slog.NewTextHandler(output, &slog.HandlerOptions{
		Level: logLevel,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(time.Now().Format("2006-01-02 15:04:05.000")),
				}
			}
			return a
		},
	})

	// Устанавливаем логгер по умолчанию
	slog.SetDefault(slog.New(handler))

	// Логгируем информацию о настройках
	slog.Info("Логгер настроен",
		key_level, config.Level,
		key_output, config.FilePath)

	return nil
}

// GetLogLevelFromEnv получает уровень логирования из переменных окружения или аргументов
func GetLogLevelFromEnv(flagValue string) string {
	// Проверяем флаг командной строки
	if flagValue != "" {
		return flagValue
	}

	// Если флаг не указан, проверяем переменную окружения
	envLogLevel := os.Getenv("LOG_LEVEL")
	if envLogLevel != "" {
		return envLogLevel
	}

	// Возвращаем значение по умолчанию
	return "info"
}

// GetLogFileFromEnv получает путь к файлу логов из переменных окружения
func GetLogFileFromEnv() string {
	return os.Getenv("LOG_FILE")
}
