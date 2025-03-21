# Census API MCP Сервер

[![Go Version](https://img.shields.io/github/go-mod/go-version/Headcrab/census_mcp)](https://go.dev)
[![License](https://github.com/Headcrab/census_mcp/blob/main/LICENSE)](LICENSE)
[![Coverage](https://codecov.io/gh/Headcrab/census_mcp/graph/badge.svg?token=WSRWMHXMTA)](https://codecov.io/gh/Headcrab/census_mcp)

Этот проект представляет собой пример MCP-совместимого сервера для работы с Census API (API переписи населения США).

## Возможности

- Получение данных о населении штатов США
- Получение данных о населении округов в штатах
- Поиск штатов по названию (полному или частичному)
- Получение списка доступных наборов данных Census API
- Получение списка переменных для заданного набора данных
- Получение доступных географических уровней
- Выполнение пользовательских запросов к Census API с указанием набора данных, года, переменных и географии
- Запуск в режиме тестирования для демонстрации работы
- Поддержка различных транспортов (stdio и SSE)

## Структура проекта

Проект имеет четкое разделение ответственности по принципам SOLID:

```tree
census-mcp/
├── app/            # Основная логика приложения
│   └── server.go   # Настройка и запуск сервера
├── census/         # Пакет для работы с Census API
│   ├── api.go      # Клиент Census API
│   └── api_test.go # Тесты Census API
├── mcp/            # Работа с протоколом MCP
│   └── census_tools.go # Инструменты MCP для Census API
└── main.go         # Точка входа
```

## Установка и запуск

### Предварительные требования

- Go 1.19 или выше
- Ключ Census API (получить можно на [официальном сайте](https://api.census.gov/data/key_signup.html))

### Установка

```bash
git clone https://github.com/yourusername/census-mcp.git
cd census-mcp
go build -o census-mcp
```

### Настройка ключа API

Есть два способа указать ключ API:

1. Через переменную окружения:
```bash
export CENSUS_API_KEY=ваш_ключ_api
```

2. Через флаг командной строки:
```bash
./census-mcp -key ваш_ключ_api
```

### Запуск

Запуск в обычном режиме:
```bash
./census-mcp
```

Запуск в тестовом режиме (для демонстрации работы без реальных запросов):
```bash
./census-mcp -test
```

Использование SSE транспорта:
```bash
./census-mcp -transport sse
```

## Развертывание с Docker

### Предварительные требования

- Docker
- Docker Compose (опционально)

### Запуск с использованием Docker

Сборка образа:
```bash
docker build -t census-mcp .
```

Запуск контейнера:
```bash
docker run -p 8080:8080 -e CENSUS_API_KEY=ваш_ключ_api census-mcp
```

### Запуск с использованием Docker Compose

1. Настройте переменные окружения в файле `.env`:
```
CENSUS_API_KEY=ваш_ключ_api
PORT=8080
```

2. Запустите сервис:
```bash
docker-compose up -d
```

3. Проверьте журналы:
```bash
docker-compose logs -f
```

Docker-образ автоматически запускает приложение в режиме SSE на порте 8080 (если не переопределено переменной PORT). Все логи сохраняются в директории `logs/`, которая монтируется как том.

## Настройка MCP клиента

Для использования Census MCP в проектах, использующих протокол MCP, необходимо настроить клиент. Ниже представлен пример конфигурации для использования через stdio:

```json
{
  "mcpServers": {
    "census": {
      "command": "/path/to/census_mcp",
      "args": ["-key", "ваш_ключ_api"],
      "disabled": false,
      "alwaysAllow": []
    }
  }
}
```

## Настройка MCP клиента c SSE

Для использования SSE транспорта (например, при работе с Docker-контейнером):

```json
{
  "mcpServers": {
    "census": {
      "url": "http://localhost:8080/sse",
      "env": {
        "CENSUS_API_KEY": "ваш_ключ_api"
      }
    }
  }
}
```

## Использование с Claude или другими поддерживаемыми LLM

Чтобы использовать Census MCP API с Claude или другими LLM, поддерживающими протокол MCP, необходимо добавить его в конфигурацию `claude_desktop_config.json`. 

Пример интеграции с Claude:

1. Создайте файл конфигурации `claude_desktop_config.json` в домашней директории:

```json
{
  "mcpServers": {
    "census": {
      "command": "path/to/census_mcp",
      "args": ["-t", "stdio", "-key", "your_api_key"],
    }
  }
}
```

2. При использовании с облачным Claude через SSE:

```json
{
  "mcpServers": {
    "census": {
      "url": "https://your-server.example.com/sse",
      "env": {
        "CENSUS_API_KEY": "your_api_key"
      }
    }
  }
}
```

3. Пример запроса к Claude с использованием Census MCP API:

```
Используя Census API, найди штаты, название которых содержит "new".
```

Claude сможет использовать инструмент `search_state_by_name` для выполнения этого запроса.

## Инструменты MCP

Сервер предоставляет следующие инструменты:

1. `get_state_population` - Получение данных о населении штатов
   - Параметр: `stateID` (опционально) - ID штата (например, "06" для Калифорнии)

2. `get_county_population` - Получение данных о населении округов
   - Параметр: `stateID` (опционально) - ID штата (например, "06" для Калифорнии)

3. `search_state_by_name` - Поиск штатов по названию
   - Параметр: `name` (обязательно) - Название для поиска (полное или частичное)

4. `get_available_datasets` - Получение списка доступных наборов данных
   - Параметров нет

5. `get_variables` - Получение списка переменных для набора данных
   - Параметр: `dataset` (обязательно) - Набор данных (например, "acs/acs1")
   - Параметр: `year` (обязательно) - Год данных (например, "2021")

6. `get_geography_levels` - Получение доступных географических уровней
   - Параметр: `dataset` (обязательно) - Набор данных (например, "acs/acs1")
   - Параметр: `year` (обязательно) - Год данных (например, "2021")

7. `get_custom_data` - Выполнение пользовательских запросов к Census API
   - Параметр: `dataset` (обязательно) - Набор данных (например, "acs/acs1")
   - Параметр: `year` (обязательно) - Год данных (например, "2021")
   - Параметр: `geoLevel` (обязательно) - Географический уровень (например, "state")
   - Параметр: `variables` (обязательно) - Массив переменных (например, ["NAME", "B01001_001E"])
   - Параметр: `geoFilter` (опционально) - Объект с фильтрами (например, {"state": "06"})

## Примеры запросов

### Получение данных о населении всех штатов
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_state_population",
    "arguments": {}
  }
}
```

### Получение данных о населении конкретного штата (Калифорния)
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_state_population",
    "arguments": {
      "stateID": "06"
    }
  }
}
```

### Поиск штата по названию
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "search_state_by_name",
    "arguments": {
      "name": "york"
    }
  }
}
```

### Получение списка доступных наборов данных
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_available_datasets",
    "arguments": {}
  }
}
```

### Получение списка переменных для набора данных
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_variables",
    "arguments": {
      "dataset": "acs/acs1",
      "year": "2021"
    }
  }
}
```

### Получение доступных географических уровней
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_geography_levels",
    "arguments": {
      "dataset": "acs/acs1",
      "year": "2021"
    }
  }
}
```

### Выполнение пользовательского запроса к Census API
```json
{
  "jsonrpc": "2.0",
  "id": "test",
  "method": "mcp.call",
  "params": {
    "tool": "get_custom_data",
    "arguments": {
      "dataset": "acs/acs1",
      "year": "2021",
      "geoLevel": "state",
      "variables": ["NAME", "B01001_001E"],
      "geoFilter": {
        "state": "06"
      }
    }
  }
}
```

## Лицензия

Этот проект лицензирован под MIT License - см. файл [LICENSE](LICENSE) для деталей.

## Вклад в проект

1. Форкните репозиторий
2. Создайте ветку для ваших изменений
3. Внесите изменения и создайте pull request

## Контакты

Создайте issue в репозитории для сообщения о проблемах или предложений по улучшению.

## Спасибо

- [@Headcrab](https://github.com/Headcrab)
