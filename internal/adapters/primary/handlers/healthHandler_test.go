package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNovoHealthHandler(t *testing.T) {
	version := "1.0.0"
	handler := NovoHealthHandler(version)
	
	if handler == nil {
		t.Fatal("NovoHealthHandler() retornou nil")
	}
	
	if handler.version != version {
		t.Errorf("NovoHealthHandler() version = %s, esperado %s", handler.version, version)
	}
}

func TestNovoHealthHandler_EmptyVersion(t *testing.T) {
	handler := NovoHealthHandler("")
	
	if handler == nil {
		t.Fatal("NovoHealthHandler() retornou nil com versão vazia")
	}
	
	if handler.version != "" {
		t.Error("NovoHealthHandler() não manteve versão vazia")
	}
}

func TestHealthHandler_HealthCheck(t *testing.T) {
	version := "1.0.0"
	handler := NovoHealthHandler(version)
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	handler.HealthCheck(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("HealthCheck() status = %d, esperado %d", w.Code, http.StatusOK)
	}
	
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("HealthCheck() Content-Type = %s, esperado application/json", contentType)
	}
	
	var response HealthResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}
	
	if response.Status != "UP" {
		t.Errorf("HealthCheck() Status = %s, esperado UP", response.Status)
	}
	
	if response.Version != version {
		t.Errorf("HealthCheck() Version = %s, esperado %s", response.Version, version)
	}
	
	if response.Message != "API está funcionando normalmente" {
		t.Errorf("HealthCheck() Message = %s, esperado 'API está funcionando normalmente'", response.Message)
	}
	
	now := time.Now()
	diff := now.Sub(response.Timestamp)
	if diff > 5*time.Second || diff < 0 {
		t.Errorf("HealthCheck() Timestamp muito antigo ou no futuro: %v", response.Timestamp)
	}
}

func TestHealthHandler_HealthCheck_MultipleRequests(t *testing.T) {
	handler := NovoHealthHandler("1.0.0")
	
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		
		handler.HealthCheck(w, req)
		
		if w.Code != http.StatusOK {
			t.Errorf("Request %d: status = %d, esperado %d", i, w.Code, http.StatusOK)
		}
		
		var response HealthResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			t.Errorf("Request %d: erro ao decodificar resposta: %v", i, err)
			continue
		}
		
		if response.Status != "UP" {
			t.Errorf("Request %d: Status = %s, esperado UP", i, response.Status)
		}
	}
}

func TestHealthHandler_HealthCheck_WithDifferentHTTPMethods(t *testing.T) {
	handler := NovoHealthHandler("1.0.0")
	
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}
	
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/health", nil)
			w := httptest.NewRecorder()
			
			handler.HealthCheck(w, req)
			
			if w.Code != http.StatusOK {
				t.Errorf("Método %s: status = %d, esperado %d", method, w.Code, http.StatusOK)
			}
		})
	}
}

func TestHealthHandler_HealthCheck_ResponseStructure(t *testing.T) {
	handler := NovoHealthHandler("2.5.1")
	
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	
	handler.HealthCheck(w, req)
	
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Erro ao decodificar resposta: %v", err)
	}
	
	requiredFields := []string{"status", "timestamp", "message", "version"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("Campo obrigatório '%s' não encontrado na resposta", field)
		}
	}
	
	if status, ok := response["status"].(string); !ok {
		t.Error("Campo 'status' não é string")
	} else if status != "UP" {
		t.Errorf("Status = %s, esperado UP", status)
	}
	
	if version, ok := response["version"].(string); !ok {
		t.Error("Campo 'version' não é string")
	} else if version != "2.5.1" {
		t.Errorf("Version = %s, esperado 2.5.1", version)
	}
	
	if message, ok := response["message"].(string); !ok {
		t.Error("Campo 'message' não é string")
	} else if message == "" {
		t.Error("Message está vazio")
	}
	
	if _, ok := response["timestamp"].(string); !ok {
		t.Error("Campo 'timestamp' não é string")
	}
}

func TestHealthHandler_HealthCheck_VersionVariations(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "Versão semântica",
			version: "1.2.3",
		},
		{
			name:    "Versão com v",
			version: "v2.0.0",
		},
		{
			name:    "Versão de desenvolvimento",
			version: "dev",
		},
		{
			name:    "Versão vazia",
			version: "",
		},
		{
			name:    "Versão com sufixo",
			version: "1.0.0-beta",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NovoHealthHandler(tt.version)
			
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			
			handler.HealthCheck(w, req)
			
			var response HealthResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Fatalf("Erro ao decodificar resposta: %v", err)
			}
			
			if response.Version != tt.version {
				t.Errorf("Version = %s, esperado %s", response.Version, tt.version)
			}
		})
	}
}

func TestHealthHandler_HealthCheck_ConcurrentRequests(t *testing.T) {
	handler := NovoHealthHandler("1.0.0")
	
	const numRequests = 100
	done := make(chan bool, numRequests)
	
	for i := 0; i < numRequests; i++ {
		go func(index int) {
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			w := httptest.NewRecorder()
			
			handler.HealthCheck(w, req)
			
			if w.Code != http.StatusOK {
				t.Errorf("Request %d: status = %d, esperado %d", index, w.Code, http.StatusOK)
			}
			
			done <- true
		}(i)
	}
	
	for i := 0; i < numRequests; i++ {
		<-done
	}
}

func TestHealthResponse_JSONMarshaling(t *testing.T) {
	response := HealthResponse{
		Status:    "UP",
		Timestamp: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		Message:   "Test message",
		Version:   "1.0.0",
	}
	
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Erro ao serializar HealthResponse: %v", err)
	}
	
	var decoded HealthResponse
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Erro ao desserializar HealthResponse: %v", err)
	}
	
	if decoded.Status != response.Status {
		t.Errorf("Status após marshal/unmarshal = %s, esperado %s", decoded.Status, response.Status)
	}
	
	if decoded.Version != response.Version {
		t.Errorf("Version após marshal/unmarshal = %s, esperado %s", decoded.Version, response.Version)
	}
	
	if decoded.Message != response.Message {
		t.Errorf("Message após marshal/unmarshal = %s, esperado %s", decoded.Message, response.Message)
	}
}

func BenchmarkHealthHandler_HealthCheck(b *testing.B) {
	handler := NovoHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.HealthCheck(w, req)
	}
}

func BenchmarkHealthHandler_HealthCheck_WithJSONDecoding(b *testing.B) {
	handler := NovoHealthHandler("1.0.0")
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.HealthCheck(w, req)
		
		var response HealthResponse
		json.NewDecoder(w.Body).Decode(&response)
	}
}
