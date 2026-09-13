package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type ServerConfig struct {
	BaseDir   string
	StaticDir string
	Port      string
}

type Exercise struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Force            *string  `json:"force"`
	Level            string   `json:"level"`
	Mechanic         *string  `json:"mechanic"`
	Equipment        *string  `json:"equipment"`
	PrimaryMuscles   []string `json:"primaryMuscles"`
	SecondaryMuscles []string `json:"secondaryMuscles"`
	Instructions     []string `json:"instructions"`
	Category         string   `json:"category"`
	Images           []string `json:"images"`
}

type WorkoutExerciseConfig struct {
	ExerciseID      string  `json:"exercise_id"`
	Sets            int     `json:"sets"`
	Reps            string  `json:"reps"`
	WeightKg        float64 `json:"weight_kg"`
	RestTimeSeconds int     `json:"rest_time_seconds"`
	IsTimeBased     bool    `json:"is_time_based"`
	TimeSeconds     int     `json:"time_seconds"`
}

type WorkoutDefinition struct {
	ID                   string                  `json:"id"`
	Type                 string                  `json:"type,omitempty"`
	NameKey              string                  `json:"name_key"`
	CategoryKey          string                  `json:"category_key"`
	Days                 []string                `json:"days"`
	EstimatedDurationMin int                     `json:"estimated_duration_min"`
	CoverImageURL        string                  `json:"cover_image_url,omitempty"`
	Exercises            []WorkoutExerciseConfig `json:"exercises"`
}

type BuildResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error,omitempty"`
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func getBaseDir() string {
	if dir := os.Getenv("FITGUIDE_BASE_DIR"); dir != "" {
		return dir
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	if strings.HasSuffix(cwd, "admin/server") {
		return filepath.Clean(filepath.Join(cwd, "../.."))
	}
	if strings.HasSuffix(cwd, "admin") {
		return filepath.Clean(filepath.Join(cwd, ".."))
	}
	return cwd
}

func getStaticDir() string {
	if dir := os.Getenv("FITGUIDE_STATIC_DIR"); dir != "" {
		return dir
	}
	base := getBaseDir()
	distWeb := filepath.Join(base, "admin", "web", "dist")
	if _, err := os.Stat(distWeb); err == nil {
		return distWeb
	}
	return filepath.Join(base, "web_dist")
}

func handleExercises(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		exercisesDir := filepath.Join(baseDir, "database", "exercises", "pt")

		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				filePath := filepath.Join(exercisesDir, fmt.Sprintf("%s.json", id))
				data, err := os.ReadFile(filePath)
				if err != nil {
					http.Error(w, `{"error": "exercise not found"}`, http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return
			}

			files, err := os.ReadDir(exercisesDir)
			if err != nil {
				http.Error(w, `{"error": "failed to read exercises directory"}`, http.StatusInternalServerError)
				return
			}

			exercises := make([]Exercise, 0, len(files))
			for _, file := range files {
				if !strings.HasSuffix(file.Name(), ".json") {
					continue
				}
				content, err := os.ReadFile(filepath.Join(exercisesDir, file.Name()))
				if err != nil {
					continue
				}
				var ex Exercise
				if err := json.Unmarshal(content, &ex); err == nil {
					exercises = append(exercises, ex)
				}
			}

			sort.Slice(exercises, func(i, j int) bool {
				return exercises[i].ID < exercises[j].ID
			})

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(exercises)

		case http.MethodPost, http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
				return
			}

			var ex Exercise
			if err := json.Unmarshal(body, &ex); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "invalid JSON: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}

			if strings.TrimSpace(ex.ID) == "" {
				http.Error(w, `{"error": "exercise id is required"}`, http.StatusBadRequest)
				return
			}

			if strings.TrimSpace(ex.Name) == "" {
				http.Error(w, `{"error": "exercise name is required"}`, http.StatusBadRequest)
				return
			}

			if ex.PrimaryMuscles == nil {
				ex.PrimaryMuscles = []string{}
			}
			if ex.SecondaryMuscles == nil {
				ex.SecondaryMuscles = []string{}
			}
			if ex.Instructions == nil {
				ex.Instructions = []string{}
			}
			if ex.Images == nil {
				ex.Images = []string{}
			}

			formatted, err := json.MarshalIndent(ex, "", "  ")
			if err != nil {
				http.Error(w, `{"error": "failed to serialize JSON"}`, http.StatusInternalServerError)
				return
			}
			formatted = append(formatted, '\n')

			filePath := filepath.Join(exercisesDir, fmt.Sprintf("%s.json", ex.ID))
			if err := os.WriteFile(filePath, formatted, 0644); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to write file: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(formatted)

		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if strings.TrimSpace(id) == "" {
				http.Error(w, `{"error": "id parameter is required"}`, http.StatusBadRequest)
				return
			}

			filePath := filepath.Join(exercisesDir, fmt.Sprintf("%s.json", id))
			if err := os.Remove(filePath); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to delete file: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))

		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleMetadata(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		metadataDir := filepath.Join(baseDir, "database", "metadata")
		topics := []string{"categories", "levels", "mechanics", "equipments", "forces", "muscles"}
		response := make(map[string]interface{})

		for _, topic := range topics {
			ptFile := filepath.Join(metadataDir, topic, "pt.json")
			enFile := filepath.Join(metadataDir, topic, "en.json")

			topicData := make(map[string]interface{})
			if data, err := os.ReadFile(ptFile); err == nil {
				var parsed interface{}
				if err := json.Unmarshal(data, &parsed); err == nil {
					topicData["pt"] = parsed
				}
			}
			if data, err := os.ReadFile(enFile); err == nil {
				var parsed interface{}
				if err := json.Unmarshal(data, &parsed); err == nil {
					topicData["en"] = parsed
				}
			}
			response[topic] = topicData
		}

		anatomyFile := filepath.Join(metadataDir, "muscles", "anatomy.json")
		if data, err := os.ReadFile(anatomyFile); err == nil {
			var parsed interface{}
			if err := json.Unmarshal(data, &parsed); err == nil {
				response["anatomy"] = parsed
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func handleWorkouts(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workoutsDir := filepath.Join(baseDir, "database", "workouts")

		switch r.Method {
		case http.MethodGet:
			result := map[string][]WorkoutDefinition{
				"seeds":     {},
				"templates": {},
			}

			readDir := func(subDir string, key string) {
				path := filepath.Join(workoutsDir, subDir)
				files, err := os.ReadDir(path)
				if err != nil {
					return
				}
				for _, f := range files {
					if !strings.HasSuffix(f.Name(), ".json") {
						continue
					}
					content, err := os.ReadFile(filepath.Join(path, f.Name()))
					if err != nil {
						continue
					}
					var wDef WorkoutDefinition
					if err := json.Unmarshal(content, &wDef); err == nil {
						wDef.Type = key
						result[key] = append(result[key], wDef)
					}
				}
			}

			readDir("seeds", "seeds")
			readDir("templates", "templates")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)

		case http.MethodPost, http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, `{"error": "invalid request body"}`, http.StatusBadRequest)
				return
			}

			var wDef WorkoutDefinition
			if err := json.Unmarshal(body, &wDef); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "invalid JSON: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}

			if strings.TrimSpace(wDef.ID) == "" {
				http.Error(w, `{"error": "workout id is required"}`, http.StatusBadRequest)
				return
			}

			subDir := "templates"
			if wDef.Type == "seeds" || strings.HasPrefix(wDef.ID, "workout_") {
				subDir = "seeds"
			}

			if wDef.Days == nil {
				wDef.Days = []string{}
			}
			if wDef.Exercises == nil {
				wDef.Exercises = []WorkoutExerciseConfig{}
			}

			targetDir := filepath.Join(workoutsDir, subDir)
			os.MkdirAll(targetDir, 0755)

			filePath := filepath.Join(targetDir, fmt.Sprintf("%s.json", wDef.ID))
			formatted, err := json.MarshalIndent(wDef, "", "  ")
			if err != nil {
				http.Error(w, `{"error": "failed to serialize JSON"}`, http.StatusInternalServerError)
				return
			}
			formatted = append(formatted, '\n')

			if err := os.WriteFile(filePath, formatted, 0644); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to write file: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(formatted)

		case http.MethodDelete:
			wType := r.URL.Query().Get("type")
			id := r.URL.Query().Get("id")
			if strings.TrimSpace(id) == "" {
				http.Error(w, `{"error": "id parameter is required"}`, http.StatusBadRequest)
				return
			}
			if wType != "seeds" && wType != "templates" {
				wType = "templates"
			}

			filePath := filepath.Join(workoutsDir, wType, fmt.Sprintf("%s.json", id))
			if err := os.Remove(filePath); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to delete workout: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"success": true}`))

		default:
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleBuild(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		scripts := []string{
			"scripts/compile_en.py",
			"scripts/compile_pt.py",
			"scripts/extract_metadata.py",
			"scripts/compile_workouts.py",
		}

		var fullOutput strings.Builder
		for _, script := range scripts {
			scriptPath := filepath.Join(baseDir, script)
			if _, err := os.Stat(scriptPath); err != nil {
				continue
			}

			cmd := exec.Command("python3", scriptPath)
			cmd.Dir = baseDir
			out, err := cmd.CombinedOutput()
			fullOutput.WriteString(fmt.Sprintf("==> Executando %s\n", script))
			fullOutput.Write(out)
			fullOutput.WriteString("\n")

			if err != nil {
				resp := BuildResponse{
					Success: false,
					Output:  fullOutput.String(),
					Error:   fmt.Sprintf("Falha ao executar %s: %s", script, err.Error()),
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		resp := BuildResponse{
			Success: true,
			Output:  fullOutput.String(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func handleUploadImage(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error":"Método não permitido"}`, http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			http.Error(w, `{"error":"Erro ao processar formulário"}`, http.StatusBadRequest)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			http.Error(w, `{"error":"Arquivo de imagem não encontrado"}`, http.StatusBadRequest)
			return
		}
		defer file.Close()

		relativePath := strings.TrimSpace(r.FormValue("relative_path"))
		if relativePath == "" {
			relativePath = filepath.Clean(handler.Filename)
		} else {
			relativePath = filepath.Clean(relativePath)
		}

		if strings.HasPrefix(relativePath, "..") || strings.HasPrefix(relativePath, "/") {
			http.Error(w, `{"error":"Caminho de imagem inválido"}`, http.StatusBadRequest)
			return
		}

		targetEn := filepath.Join(baseDir, "database", "exercises", "en", relativePath)
		targetPt := filepath.Join(baseDir, "database", "exercises", "pt", relativePath)

		if err := os.MkdirAll(filepath.Dir(targetEn), 0755); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao criar diretório: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		if err := os.MkdirAll(filepath.Dir(targetPt), 0755); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao criar diretório: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao ler imagem: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(targetEn, fileBytes, 0644); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao salvar imagem em exercises_en: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}
		if err := os.WriteFile(targetPt, fileBytes, 0644); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao salvar imagem em exercises_pt: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"path":    relativePath,
			"url":     "/images/" + relativePath,
		})
	}
}

func handleServeImages(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		relPath := strings.TrimPrefix(r.URL.Path, "/images/")
		relPath = filepath.Clean(relPath)

		if strings.HasPrefix(relPath, "..") {
			http.NotFound(w, r)
			return
		}

		candidates := []string{
			filepath.Join(baseDir, "database", "exercises", "en", relPath),
			filepath.Join(baseDir, "database", "exercises", "pt", relPath),
		}

		for _, path := range candidates {
			if info, err := os.Stat(path); err == nil && !info.IsDir() {
				http.ServeFile(w, r, path)
				return
			}
		}

		http.NotFound(w, r)
	}
}

func handleSPA(staticDir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(staticDir))
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/images/") {
			http.NotFound(w, r)
			return
		}

		path := filepath.Join(staticDir, filepath.Clean(r.URL.Path))
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		indexPath := filepath.Join(staticDir, "index.html")
		if _, err := os.Stat(indexPath); err == nil {
			http.ServeFile(w, r, indexPath)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head><title>FitGuide Admin API</title></head>
<body style="font-family:sans-serif;padding:40px;background:#0f172a;color:#f8fafc">
  <h1>FitGuide Admin API</h1>
  <p>O servidor Go está em execução.</p>
  <p>Endpoints disponíveis:</p>
  <ul>
    <li><code>GET /api/exercises</code></li>
    <li><code>GET /api/metadata</code></li>
    <li><code>GET /api/workouts</code></li>
    <li><code>POST /api/upload-image</code></li>
    <li><code>POST /api/build</code></li>
    <li><code>GET /images/*</code></li>
  </ul>
</body>
</html>`))
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	baseDir := getBaseDir()
	staticDir := getStaticDir()

	mux := http.NewServeMux()

	mux.HandleFunc("/api/exercises", enableCORS(handleExercises(baseDir)))
	mux.HandleFunc("/api/metadata", enableCORS(handleMetadata(baseDir)))
	mux.HandleFunc("/api/workouts", enableCORS(handleWorkouts(baseDir)))
	mux.HandleFunc("/api/upload-image", enableCORS(handleUploadImage(baseDir)))
	mux.HandleFunc("/api/build", enableCORS(handleBuild(baseDir)))
	mux.HandleFunc("/images/", enableCORS(handleServeImages(baseDir)))
	mux.HandleFunc("/", handleSPA(staticDir))

	addr := fmt.Sprintf("0.0.0.0:%s", port)
	fmt.Printf("FitGuide Admin Server rodando em http://localhost:%s\n", port)
	fmt.Printf("Base directory: %s\n", baseDir)
	fmt.Printf("Static directory: %s\n", staticDir)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Printf("Erro ao iniciar servidor: %s\n", err.Error())
		os.Exit(1)
	}
}
