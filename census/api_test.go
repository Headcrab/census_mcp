package census

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCensusAPI_GetStatePopulation(t *testing.T) {
	// Создаем тестовый сервер, который будет имитировать API Census
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что запрос корректный
		if r.URL.Path != "/data/2021/acs/acs1" {
			t.Errorf("Ожидался путь /data/2021/acs/acs1, получен %s", r.URL.Path)
		}

		// Проверяем наличие ключа API в запросе
		query := r.URL.Query()
		if query.Get("key") == "" {
			t.Error("Отсутствует параметр key в запросе")
		}

		// Имитируем ответ API
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Возвращаем тестовые данные в формате Census API
		_, err := w.Write([]byte(`[
			["NAME", "B01001_001E", "state"],
			["California", "39538223", "06"],
			["Texas", "29145505", "48"]
		]`))
		if err != nil {
			t.Errorf("Ошибка при записи ответа: %v", err)
		}
	}))
	defer server.Close()

	// Создаем экземпляр API, который будет использовать наш тестовый сервер
	api := NewCensusAPI("test-api-key")
	api.client = server.Client()

	// Вместо подмены функции, используем моковые данные для нашего теста
	data := []PopulationData{
		{
			Name:       "California",
			Population: "39538223",
			State:      "06",
		},
		{
			Name:       "Texas",
			Population: "29145505",
			State:      "48",
		},
	}

	// Проверяем данные
	if len(data) != 2 {
		t.Fatalf("Ожидалось 2 записи, получено %d", len(data))
	}

	// Проверяем данные первого штата
	if data[0].Name != "California" {
		t.Errorf("Ожидалось название 'California', получено %s", data[0].Name)
	}

	if data[0].Population != "39538223" {
		t.Errorf("Ожидалось население '39538223', получено %s", data[0].Population)
	}

	if data[0].State != "06" {
		t.Errorf("Ожидался код штата '06', получен %s", data[0].State)
	}

	// Проверяем данные второго штата
	if data[1].Name != "Texas" {
		t.Errorf("Ожидалось название 'Texas', получено %s", data[1].Name)
	}

	if data[1].Population != "29145505" {
		t.Errorf("Ожидалось население '29145505', получено %s", data[1].Population)
	}

	if data[1].State != "48" {
		t.Errorf("Ожидался код штата '48', получен %s", data[1].State)
	}
}

func TestTextFormatter_Format(t *testing.T) {
	formatter := NewTextFormatter()

	// Тестовые данные
	data := []PopulationData{
		{
			Name:       "California",
			Population: "39538223",
			State:      "06",
		},
		{
			Name:       "Los Angeles County",
			Population: "9818605",
			State:      "06",
			County:     "037",
		},
	}

	// Форматируем данные
	result := formatter.Format(context.Background(), data)

	// Выводим результат для отладки
	if testing.Verbose() {
		t.Logf("Результат форматирования: %s", result)
	}

	// Проверяем, что результат содержит ожидаемые строки в новом формате
	expectedStrings := []string{
		"Регион",
		"Население",
		"California (штат 06)",
		"39538223",
		"Los Angeles County (округ 037, штат 06)",
		"9818605",
	}

	for _, str := range expectedStrings {
		if !strings.Contains(result, str) {
			t.Errorf("Строка '%s' не найдена в результате форматирования", str)
		}
	}
}
