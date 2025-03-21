package mcp

import (
	"census_mcp/census"
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// ToolHandler определяет интерфейс для обработчика инструментов MCP
type ToolHandler interface {
	// HandleCountLettersTool обрабатывает запрос на подсчет букв
	HandleCountLettersTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

// DefaultToolHandler - стандартная реализация обработчика инструментов
type DefaultToolHandler struct {
	counter   census.Formatter
	formatter census.Formatter
}

// NewToolHandler создает новый экземпляр обработчика инструментов
func NewToolHandler(counter census.Formatter, formatter census.Formatter) ToolHandler {
	return &DefaultToolHandler{
		counter:   counter,
		formatter: formatter,
	}
}

// HandleCountLettersTool обрабатывает запрос на подсчет букв
func (h *DefaultToolHandler) HandleCountLettersTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	word, ok1 := arguments["word"].(string)
	letters, ok2 := arguments["letters"].(string)

	if !ok1 || !ok2 {
		return mcp.NewToolResultError("Необходимо указать параметры 'word' и 'letters'"), nil
	}

	// Использование счетчика букв для получения результатов
	counts := h.counter.CountLetters(word, letters)

	// Форматирование результатов
	result := h.formatter.Format(counts)

	// Возвращаем результат
	return mcp.NewToolResultText(result), nil
}

// RegisterTools регистрирует инструменты MCP
func RegisterTools(mcpServer *server.MCPServer, handler ToolHandler) {
	mcpServer.AddTool(mcp.NewTool("count_letters",
		mcp.WithDescription("Подсчитывает количество определённых букв в слове"),
		mcp.WithString("word",
			mcp.Description("Слово, в котором нужно посчитать буквы"),
			mcp.Required(),
		),
		mcp.WithString("letters",
			mcp.Description("Буквы, которые нужно посчитать (например 'аео')"),
			mcp.Required(),
		),
	), handler.HandleCountLettersTool)
}
