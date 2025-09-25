package function

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type requestBody struct {
	CPF string `json:"cpf"`
}

type cliente struct {
	ID    string `json:"id"`
	Nome  string `json:"nome"`
	CPF   string `json:"cpf"`
	Email string `json:"email"`
}

type responseBody struct {
	Token   string  `json:"token"`
	Cliente cliente `json:"cliente"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func connect() (*sql.DB, error) {
	// Use Unix socket for Cloud SQL Auth Proxy
	dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/sylvan-willow-473218-f8:us-central1:soat-fiap-mysql)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func validateCPF(raw string) (string, bool) {
	r := regexp.MustCompile("[^0-9]")
	cpf := r.ReplaceAllString(raw, "")
	if len(cpf) != 11 {
		return "", false
	}
	return cpf, true
}

func findClienteByCPF(db *sql.DB, cpf string) (*cliente, error) {
	row := db.QueryRow("SELECT id, nome, cpf, email FROM clientes WHERE cpf = ?", cpf)
	c := &cliente{}
	if err := row.Scan(&c.ID, &c.Nome, &c.CPF, &c.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func generateJWT(c *cliente, secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   c.ID,
		"cpf":   c.CPF,
		"nome":  c.Nome,
		"email": c.Email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthByCPF(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorResponse{Error: "method not allowed"})
		return
	}
	var req requestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid json"})
		return
	}
	cpf, ok := validateCPF(req.CPF)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "cpf inválido"})
		return
	}
	db, err := connect()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "db connect error"})
		return
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "db ping error"})
		return
	}
	c, err := findClienteByCPF(db, cpf)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "db query error"})
		return
	}
	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse{Error: "cliente não encontrado"})
		return
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "change-me"
	}
	tkn, err := generateJWT(c, secret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "jwt error"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseBody{Token: tkn, Cliente: *c})
}
