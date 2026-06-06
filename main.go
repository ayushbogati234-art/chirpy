package main

import (
	"encoding/json"
        "log"
	"fmt"
	"net/http"
	"strings"
	"database/sql"
	"os"
"time"

        "github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	 "chirpy/internal/database"
"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	 platform       string
	dbQueries      *database.Queries
}

type chirpRequest struct {
	Body string `json:"body"`
    
}

type createUserRequest struct {
        Email string `json:"email"`
}

type User struct {
        ID        uuid.UUID `json:"id"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
        Email     string    `json:"email"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func cleanProfanity(text string) string {
	words := strings.Split(text, " ")

	for i, word := range words {
		switch strings.ToLower(word) {
		case "kerfuffle", "sharbert", "fornax":
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	var req chirpRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong")
		return
	}

	    

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	respondWithJSON(w, http.StatusOK, response{
		CleanedBody: cleanProfanity(req.Body),
	})
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
        var req createUserRequest

        err := json.NewDecoder(r.Body).Decode(&req)
        if err != nil {
                respondWithError(w, http.StatusBadRequest, "Invalid request")
                return
        }

        dbUser, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Could not create user")
                return
        }

        user := User{
                ID:        dbUser.ID,
                CreatedAt: dbUser.CreatedAt,
                UpdatedAt: dbUser.UpdatedAt,
                Email:     dbUser.Email,
        }

        respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	hits := cfg.fileserverHits.Load()

	html := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	 if cfg.platform != "dev" {
                w.WriteHeader(http.StatusForbidden)
                return
        }

        cfg.fileserverHits.Store(0)

        err := cfg.dbQueries.DeleteAllUsers(r.Context())
        if err != nil {
                respondWithError(w, http.StatusInternalServerError, "Failed to delete users")
                return
        }

        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Hits: 0"))
}

func main() {
        err := godotenv.Load()
        if err != nil {
                log.Fatal("Error loading .env file")
        }

        dbURL := os.Getenv("DB_URL")
        platform := os.Getenv("PLATFORM")
        db, err := sql.Open("postgres", dbURL)
        if err != nil {
                log.Fatal(err)
        }

        dbQueries := database.New(db)

        mux := http.NewServeMux()

        apiCfg := &apiConfig{
                fileserverHits: atomic.Int32{},
                dbQueries:      dbQueries,
				platform:       platform,
        }

        // Health check endpoint
        mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Content-Type", "text/plain; charset=utf-8")
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("OK"))
        })

        // Chirp validation endpoint
        mux.HandleFunc("POST /api/validate_chirp", apiCfg.validateChirpHandler)
        mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)
        // File server
        fs := http.FileServer(http.Dir("."))

        mux.Handle(
                "/app/",
                http.StripPrefix(
                        "/app/",
                        apiCfg.middlewareMetricsInc(fs),
                ),
        )

        // Admin endpoints
        mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
        mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
		
       
        log.Fatal(http.ListenAndServe(":8080", mux))
}
