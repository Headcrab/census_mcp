package mcp

import (
	"census_mcp/census"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
)

// MockCensusAPIClient - мок для интерфейса CensusAPIClient
type MockCensusAPIClient struct {
	GetStatePopulationFunc   func(stateID string) ([]census.PopulationData, error)
	GetCountyPopulationFunc  func(stateID string) ([]census.PopulationData, error)
	SearchStateByNameFunc    func(name string) ([]census.PopulationData, error)
	GetAvailableDatasetsFunc func() ([]census.DatasetInfo, error)
	GetVariablesFunc         func(dataset, year string) (map[string]census.VariableInfo, error)
	GetGeographyLevelsFunc   func(dataset, year string) ([]census.GeographyLevel, error)
	GetCustomDataFunc        func(request census.CustomDataRequest) ([]map[string]string, error)
}

// MockFormatter - мок для интерфейса Formatter
type MockFormatter struct {
	FormatFunc func(data interface{}) string
}

func (m *MockFormatter) Format(data interface{}) string {
	return m.FormatFunc(data)
}

func (m *MockCensusAPIClient) GetStatePopulation(stateID string) ([]census.PopulationData, error) {
	return m.GetStatePopulationFunc(stateID)
}

func (m *MockCensusAPIClient) GetCountyPopulation(stateID string) ([]census.PopulationData, error) {
	return m.GetCountyPopulationFunc(stateID)
}

func (m *MockCensusAPIClient) SearchStateByName(name string) ([]census.PopulationData, error) {
	return m.SearchStateByNameFunc(name)
}

func (m *MockCensusAPIClient) GetAvailableDatasets() ([]census.DatasetInfo, error) {
	return m.GetAvailableDatasetsFunc()
}

func (m *MockCensusAPIClient) GetVariables(dataset, year string) (map[string]census.VariableInfo, error) {
	return m.GetVariablesFunc(dataset, year)
}

func (m *MockCensusAPIClient) GetGeographyLevels(dataset, year string) ([]census.GeographyLevel, error) {
	return m.GetGeographyLevelsFunc(dataset, year)
}

func (m *MockCensusAPIClient) GetCustomData(request census.CustomDataRequest) ([]map[string]string, error) {
	return m.GetCustomDataFunc(request)
}

// CreateMockCallToolRequest создает моковый запрос для тестирования
func CreateMockCallToolRequest(args map[string]interface{}) mcp.CallToolRequest {
	mockRequest := mcp.CallToolRequest{}
	mockRequest.Params.Arguments = args
	return mockRequest
}

// GetContentAsString извлекает текстовое содержимое из Content
func GetContentAsString(content []mcp.Content) string {
	if len(content) == 0 {
		return ""
	}

	// Проверяем есть ли TextContent
	for _, c := range content {
		if tc, ok := c.(mcp.TextContent); ok {
			return tc.Text
		}
	}

	return ""
}

func TestCensusDefaultToolHandler_HandleGetStatePopulationTool(t *testing.T) {
	tests := []struct {
		name           string
		stateID        string
		mockData       []census.PopulationData
		mockError      error
		expectedOutput string
		expectError    bool
	}{
		{
			name:    "Успешное получение данных о населении штата",
			stateID: "06",
			mockData: []census.PopulationData{
				{
					Name:       "California",
					Population: "39538223",
					State:      "06",
				},
			},
			mockError:      nil,
			expectedOutput: "Форматированные данные о населении",
			expectError:    false,
		},
		{
			name:           "Ошибка при получении данных",
			stateID:        "99",
			mockData:       nil,
			mockError:      errors.New("ошибка API"),
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-API
			mockAPI := &MockCensusAPIClient{
				GetStatePopulationFunc: func(stateID string) ([]census.PopulationData, error) {
					assert.Equal(t, tt.stateID, stateID)
					return tt.mockData, tt.mockError
				},
			}

			// Создаем мок-форматтер
			mockFormatter := &MockFormatter{
				FormatFunc: func(data interface{}) string {
					if tt.expectError {
						t.Fatalf("Formatter не должен вызываться при ошибке")
					}
					popData, ok := data.([]census.PopulationData)
					assert.True(t, ok)
					assert.Equal(t, tt.mockData, popData)
					return tt.expectedOutput
				},
			}

			// Создаем обработчик
			handler := NewCensusToolHandler(mockAPI, mockFormatter)

			// Создаем запрос с правильной структурой
			request := CreateMockCallToolRequest(map[string]interface{}{
				"stateID": tt.stateID,
			})

			// Вызываем тестируемый метод
			result, err := handler.HandleGetStatePopulationTool(context.Background(), request)

			// Проверяем результаты
			assert.NoError(t, err)

			// Получаем текстовое содержимое
			contentText := GetContentAsString(result.Content)

			if tt.expectError {
				assert.True(t, strings.Contains(contentText, "Ошибка при получении данных"),
					"Ожидалось сообщение об ошибке с текстом 'Ошибка при получении данных', получено: %s", contentText)
				if tt.mockError != nil {
					assert.True(t, strings.Contains(contentText, tt.mockError.Error()),
						"Ожидалось сообщение с текстом ошибки '%s', получено: %s", tt.mockError.Error(), contentText)
				}
			} else {
				assert.Equal(t, tt.expectedOutput, contentText)
			}
		})
	}
}

func TestCensusDefaultToolHandler_HandleGetCountyPopulationTool(t *testing.T) {
	tests := []struct {
		name           string
		stateID        string
		mockData       []census.PopulationData
		mockError      error
		expectedOutput string
		expectError    bool
	}{
		{
			name:    "Успешное получение данных о населении округов",
			stateID: "06",
			mockData: []census.PopulationData{
				{
					Name:       "Los Angeles County",
					Population: "10014009",
					State:      "06",
					County:     "037",
				},
			},
			mockError:      nil,
			expectedOutput: "Форматированные данные о населении округов",
			expectError:    false,
		},
		{
			name:           "Ошибка при получении данных",
			stateID:        "99",
			mockData:       nil,
			mockError:      errors.New("ошибка API"),
			expectedOutput: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-API
			mockAPI := &MockCensusAPIClient{
				GetCountyPopulationFunc: func(stateID string) ([]census.PopulationData, error) {
					assert.Equal(t, tt.stateID, stateID)
					return tt.mockData, tt.mockError
				},
			}

			// Создаем мок-форматтер
			mockFormatter := &MockFormatter{
				FormatFunc: func(data interface{}) string {
					if tt.expectError {
						t.Fatalf("Formatter не должен вызываться при ошибке")
					}
					popData, ok := data.([]census.PopulationData)
					assert.True(t, ok)
					assert.Equal(t, tt.mockData, popData)
					return tt.expectedOutput
				},
			}

			// Создаем обработчик
			handler := NewCensusToolHandler(mockAPI, mockFormatter)

			// Создаем запрос с правильной структурой
			request := CreateMockCallToolRequest(map[string]interface{}{
				"stateID": tt.stateID,
			})

			// Вызываем тестируемый метод
			result, err := handler.HandleGetCountyPopulationTool(context.Background(), request)

			// Проверяем результаты
			assert.NoError(t, err)

			// Получаем текстовое содержимое
			contentText := GetContentAsString(result.Content)

			if tt.expectError {
				assert.True(t, strings.Contains(contentText, "Ошибка при получении данных"),
					"Ожидалось сообщение об ошибке с текстом 'Ошибка при получении данных', получено: %s", contentText)
				if tt.mockError != nil {
					assert.True(t, strings.Contains(contentText, tt.mockError.Error()),
						"Ожидалось сообщение с текстом ошибки '%s', получено: %s", tt.mockError.Error(), contentText)
				}
			} else {
				assert.Equal(t, tt.expectedOutput, contentText)
			}
		})
	}
}

func TestCensusDefaultToolHandler_HandleSearchStateByNameTool(t *testing.T) {
	tests := []struct {
		name           string
		stateName      string
		mockData       []census.PopulationData
		mockError      error
		expectedOutput string
		expectError    bool
		emptyParams    bool
	}{
		{
			name:      "Успешный поиск штата",
			stateName: "California",
			mockData: []census.PopulationData{
				{
					Name:       "California",
					Population: "39538223",
					State:      "06",
				},
			},
			mockError:      nil,
			expectedOutput: "Форматированные данные о штате",
			expectError:    false,
			emptyParams:    false,
		},
		{
			name:           "Ошибка при поиске",
			stateName:      "Unknown",
			mockData:       nil,
			mockError:      errors.New("ошибка API"),
			expectedOutput: "",
			expectError:    true,
			emptyParams:    false,
		},
		{
			name:           "Пустое название штата",
			stateName:      "",
			mockData:       nil,
			mockError:      nil,
			expectedOutput: "",
			expectError:    true,
			emptyParams:    true,
		},
		{
			name:           "Штат не найден",
			stateName:      "Unknown",
			mockData:       []census.PopulationData{},
			mockError:      nil,
			expectedOutput: "Штаты не найдены по запросу: Unknown",
			expectError:    false,
			emptyParams:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-API
			mockAPI := &MockCensusAPIClient{
				SearchStateByNameFunc: func(name string) ([]census.PopulationData, error) {
					if !tt.emptyParams {
						assert.Equal(t, tt.stateName, name)
					}
					return tt.mockData, tt.mockError
				},
			}

			// Создаем мок-форматтер
			mockFormatter := &MockFormatter{
				FormatFunc: func(data interface{}) string {
					if tt.expectError || len(tt.mockData) == 0 {
						return tt.expectedOutput
					}
					popData, ok := data.([]census.PopulationData)
					assert.True(t, ok)
					assert.Equal(t, tt.mockData, popData)
					return tt.expectedOutput
				},
			}

			// Создаем обработчик
			handler := NewCensusToolHandler(mockAPI, mockFormatter)

			// Создаем запрос с правильной структурой
			var request mcp.CallToolRequest
			if tt.emptyParams {
				request = CreateMockCallToolRequest(map[string]interface{}{})
			} else {
				request = CreateMockCallToolRequest(map[string]interface{}{
					"name": tt.stateName,
				})
			}

			// Вызываем тестируемый метод
			result, err := handler.HandleSearchStateByNameTool(context.Background(), request)

			// Проверяем результаты
			assert.NoError(t, err)

			// Получаем текстовое содержимое
			contentText := GetContentAsString(result.Content)

			if tt.emptyParams {
				assert.True(t, strings.Contains(contentText, "Необходимо указать параметр"),
					"Ожидалось сообщение о необходимости указать параметр, получено: %s", contentText)
			} else if tt.expectError {
				assert.True(t, strings.Contains(contentText, "Ошибка при поиске"),
					"Ожидалось сообщение об ошибке с текстом 'Ошибка при поиске', получено: %s", contentText)
				if tt.mockError != nil {
					assert.True(t, strings.Contains(contentText, tt.mockError.Error()),
						"Ожидалось сообщение с текстом ошибки '%s', получено: %s", tt.mockError.Error(), contentText)
				}
			} else if len(tt.mockData) == 0 {
				assert.True(t, strings.Contains(contentText, "не найдены"),
					"Ожидалось сообщение о том, что штаты не найдены, получено: %s", contentText)
				assert.True(t, strings.Contains(contentText, tt.stateName),
					"Ожидалось сообщение, содержащее название штата '%s', получено: %s", tt.stateName, contentText)
			} else {
				assert.Equal(t, tt.expectedOutput, contentText)
			}
		})
	}
}

func TestNewCensusToolHandler(t *testing.T) {
	// Создаем мок-объекты
	mockAPI := &MockCensusAPIClient{}
	mockFormatter := &MockFormatter{}

	// Проверяем создание обработчика
	handler := NewCensusToolHandler(mockAPI, mockFormatter)

	// Проверяем, что обработчик не nil и имеет правильный тип
	assert.NotNil(t, handler)
	_, ok := handler.(*CensusDefaultToolHandler)
	assert.True(t, ok)
}

// Тест для RegisterCensusTools - проверка регистрации инструментов
func TestRegisterCensusTools(t *testing.T) {
	// Пропускаем тест для функции регистрации инструментов
	// Тестирование RegisterCensusTools требует создания полного мока для server.MCPServer,
	// что выходит за рамки данных тестов. В будущем можно реализовать этот тест при необходимости.
	t.Skip("Требуется реализация моков для MCPServer")
}

// MockCensusToolHandler - мок для интерфейса CensusToolHandler
type MockCensusToolHandler struct {
	HandleGetStatePopulationToolFunc   func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleGetCountyPopulationToolFunc  func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleSearchStateByNameToolFunc    func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleGetAvailableDatasetsToolFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleGetVariablesToolFunc         func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleGetGeographyLevelsToolFunc   func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HandleGetCustomDataToolFunc        func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func (m *MockCensusToolHandler) HandleGetStatePopulationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetStatePopulationToolFunc != nil {
		return m.HandleGetStatePopulationToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleGetCountyPopulationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetCountyPopulationToolFunc != nil {
		return m.HandleGetCountyPopulationToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleSearchStateByNameTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleSearchStateByNameToolFunc != nil {
		return m.HandleSearchStateByNameToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleGetAvailableDatasetsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetAvailableDatasetsToolFunc != nil {
		return m.HandleGetAvailableDatasetsToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleGetVariablesTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetVariablesToolFunc != nil {
		return m.HandleGetVariablesToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleGetGeographyLevelsTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetGeographyLevelsToolFunc != nil {
		return m.HandleGetGeographyLevelsToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}

func (m *MockCensusToolHandler) HandleGetCustomDataTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if m.HandleGetCustomDataToolFunc != nil {
		return m.HandleGetCustomDataToolFunc(ctx, request)
	}
	return mcp.NewToolResultText("mock"), nil
}
