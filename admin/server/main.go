package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
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

type CategoryPayload struct {
	Key     string `json:"key"`
	LabelPt string `json:"label_pt"`
	LabelEn string `json:"label_en"`
}

type TranslateRequest struct {
	Text   string `json:"text"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type TranslateResponse struct {
	Translated string `json:"translated"`
	Source     string `json:"source"`
}

type TranslationMismatch struct {
	ID        string `json:"id"`
	NamePt    string `json:"name_pt"`
	NameEn    string `json:"name_en"`
	Field     string `json:"field"`
	ValuePt   string `json:"value_pt"`
	ValueEn   string `json:"value_en"`
	Severity  string `json:"severity"`
	Message   string `json:"message"`
}

type TranslationAuditReport struct {
	Languages       []string              `json:"languages"`
	Counts          map[string]int        `json:"counts"`
	MissingInEn     []string              `json:"missing_in_en"`
	MissingInPt     []string              `json:"missing_in_pt"`
	StepMismatches  []TranslationMismatch `json:"step_mismatches"`
	ImageMismatches []TranslationMismatch `json:"image_mismatches"`
	TotalDiscrepancies int                `json:"total_discrepancies"`
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
		lang := r.URL.Query().Get("lang")
		if lang == "" {
			lang = "pt"
		}
		lang = filepath.Clean(lang)

		exercisesBaseDir := filepath.Join(baseDir, "database", "exercises")

		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id != "" {
				filePath := filepath.Join(exercisesBaseDir, id, fmt.Sprintf("%s.json", lang))
				data, err := os.ReadFile(filePath)
				if err != nil {
					http.Error(w, `{"error": "exercise not found"}`, http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
				return
			}

			entries, err := os.ReadDir(exercisesBaseDir)
			if err != nil {
				http.Error(w, `{"error": "failed to read exercises directory"}`, http.StatusInternalServerError)
				return
			}

			exercises := make([]Exercise, 0, len(entries))
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				filePath := filepath.Join(exercisesBaseDir, entry.Name(), fmt.Sprintf("%s.json", lang))
				content, err := os.ReadFile(filePath)
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

			exFolder := filepath.Join(exercisesBaseDir, ex.ID)
			if err := os.MkdirAll(exFolder, 0755); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to create directory: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			filePath := filepath.Join(exFolder, fmt.Sprintf("%s.json", lang))
			if err := os.WriteFile(filePath, formatted, 0644); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to write file: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(formatted)

		case http.MethodDelete:
			id := strings.TrimSpace(r.URL.Query().Get("id"))
			if id == "" || id == "undefined" || id == "null" {
				http.Error(w, `{"error": "id parameter is required and must be valid"}`, http.StatusBadRequest)
				return
			}
			id = filepath.Clean(id)
			if strings.Contains(id, "..") || strings.Contains(id, "/") || strings.Contains(id, "\\") {
				http.Error(w, `{"error": "invalid id parameter"}`, http.StatusBadRequest)
				return
			}

			exFolder := filepath.Join(exercisesBaseDir, id)
			if _, err := os.Stat(exFolder); os.IsNotExist(err) {
				http.Error(w, fmt.Sprintf(`{"error": "exercise %s not found"}`, id), http.StatusNotFound)
				return
			}

			if err := os.RemoveAll(exFolder); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "failed to delete exercise: %s"}`, err.Error()), http.StatusInternalServerError)
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

func sanitizeMetadataKey(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	var b strings.Builder
	for _, r := range s {
		switch r {
		case 'á', 'à', 'ã', 'â', 'ä':
			b.WriteRune('a')
		case 'é', 'è', 'ê', 'ë':
			b.WriteRune('e')
		case 'í', 'ì', 'î', 'ï':
			b.WriteRune('i')
		case 'ó', 'ò', 'õ', 'ô', 'ö':
			b.WriteRune('o')
		case 'ú', 'ù', 'û', 'ü':
			b.WriteRune('u')
		case 'ç':
			b.WriteRune('c')
		case '_':
			b.WriteRune('_')
		default:
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
			}
		}
	}
	res := b.String()
	for strings.Contains(res, "__") {
		res = strings.ReplaceAll(res, "__", "_")
	}
	return strings.Trim(res, "_")
}

func executeTranslation(baseDir, text, source, target string) string {
	if strings.TrimSpace(text) == "" {
		return ""
	}
	script := fmt.Sprintf(`import sys, os
sys.path.insert(0, os.getcwd())
try:
    from scripts.translator import translate_term
    print(translate_term(%q, source=%q, target=%q), end="")
except Exception:
    print(%q, end="")
`, text, source, target, text)

	cmd := exec.Command("python3", "-c", script)
	cmd.Dir = baseDir
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return text
	}
	return strings.TrimSpace(string(out))
}

func handleMetadataCategory(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ptPath := filepath.Join(baseDir, "database", "metadata", "categories", "pt.json")
		enPath := filepath.Join(baseDir, "database", "metadata", "categories", "en.json")

		loadMap := func(path string) map[string]string {
			m := make(map[string]string)
			data, err := os.ReadFile(path)
			if err == nil {
				json.Unmarshal(data, &m)
			}
			return m
		}

		saveMap := func(path string, m map[string]string) error {
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			formatted, err := json.MarshalIndent(m, "", "  ")
			if err != nil {
				return err
			}
			formatted = append(formatted, '\n')
			return os.WriteFile(path, formatted, 0644)
		}

		switch r.Method {
		case http.MethodPost, http.MethodPut:
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, `{"error": "corpo de requisição inválido"}`, http.StatusBadRequest)
				return
			}

			var payload CategoryPayload
			if err := json.Unmarshal(body, &payload); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "JSON inválido: %s"}`, err.Error()), http.StatusBadRequest)
				return
			}

			payload.LabelPt = strings.TrimSpace(payload.LabelPt)
			payload.LabelEn = strings.TrimSpace(payload.LabelEn)
			payload.Key = strings.TrimSpace(payload.Key)

			if payload.LabelPt == "" && payload.LabelEn == "" && payload.Key == "" {
				http.Error(w, `{"error": "nome da categoria é obrigatório"}`, http.StatusBadRequest)
				return
			}

			if payload.Key == "" {
				if payload.LabelPt != "" {
					payload.Key = sanitizeMetadataKey(payload.LabelPt)
				} else {
					payload.Key = sanitizeMetadataKey(payload.LabelEn)
				}
			} else {
				payload.Key = sanitizeMetadataKey(payload.Key)
			}

			if payload.Key == "" {
				http.Error(w, `{"error": "chave de categoria inválida"}`, http.StatusBadRequest)
				return
			}

			if payload.LabelPt == "" && payload.LabelEn != "" {
				payload.LabelPt = executeTranslation(baseDir, payload.LabelEn, "en", "pt")
			}

			if payload.LabelEn == "" && payload.LabelPt != "" {
				payload.LabelEn = executeTranslation(baseDir, payload.LabelPt, "pt", "en")
			}

			ptMap := loadMap(ptPath)
			enMap := loadMap(enPath)

			ptMap[payload.Key] = payload.LabelPt
			enMap[payload.Key] = payload.LabelEn

			if err := saveMap(ptPath, ptMap); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "falha ao salvar pt.json: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			if err := saveMap(enPath, enMap); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "falha ao salvar en.json: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"category": payload,
			})

		case http.MethodDelete:
			key := r.URL.Query().Get("key")
			if strings.TrimSpace(key) == "" {
				var bodyData struct {
					Key string `json:"key"`
				}
				if body, err := io.ReadAll(r.Body); err == nil {
					json.Unmarshal(body, &bodyData)
					key = bodyData.Key
				}
			}

			key = sanitizeMetadataKey(key)
			if key == "" {
				http.Error(w, `{"error": "parâmetro key é obrigatório"}`, http.StatusBadRequest)
				return
			}

			ptMap := loadMap(ptPath)
			enMap := loadMap(enPath)

			delete(ptMap, key)
			delete(enMap, key)

			if err := saveMap(ptPath, ptMap); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "falha ao atualizar pt.json: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			if err := saveMap(enPath, enMap); err != nil {
				http.Error(w, fmt.Sprintf(`{"error": "falha ao atualizar en.json: %s"}`, err.Error()), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"key":     key,
			})

		default:
			http.Error(w, `{"error": "método não permitido"}`, http.StatusMethodNotAllowed)
		}
	}
}

func handleTranslate(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, `{"error": "método não permitido"}`, http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, `{"error": "corpo de requisição inválido"}`, http.StatusBadRequest)
			return
		}

		var req TranslateRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "JSON inválido: %s"}`, err.Error()), http.StatusBadRequest)
			return
		}

		req.Text = strings.TrimSpace(req.Text)
		if req.Text == "" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(TranslateResponse{Translated: "", Source: ""})
			return
		}

		if req.Source == "" {
			req.Source = "pt"
		}
		if req.Target == "" {
			req.Target = "en"
		}

		translated := executeTranslation(baseDir, req.Text, req.Source, req.Target)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TranslateResponse{
			Translated: translated,
			Source:     req.Text,
		})
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

func handleTranslationAudit(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		exercisesBaseDir := filepath.Join(baseDir, "database", "exercises")

		ptMap := make(map[string]Exercise)
		enMap := make(map[string]Exercise)

		entries, err := os.ReadDir(exercisesBaseDir)
		if err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				exID := entry.Name()

				ptFile := filepath.Join(exercisesBaseDir, exID, "pt.json")
				if content, err := os.ReadFile(ptFile); err == nil {
					var ex Exercise
					if err := json.Unmarshal(content, &ex); err == nil && ex.ID != "" {
						ptMap[ex.ID] = ex
					}
				}

				enFile := filepath.Join(exercisesBaseDir, exID, "en.json")
				if content, err := os.ReadFile(enFile); err == nil {
					var ex Exercise
					if err := json.Unmarshal(content, &ex); err == nil && ex.ID != "" {
						enMap[ex.ID] = ex
					}
				}
			}
		}

		report := TranslationAuditReport{
			Languages:       []string{"pt", "en"},
			Counts:          map[string]int{"pt": len(ptMap), "en": len(enMap)},
			MissingInEn:     []string{},
			MissingInPt:     []string{},
			StepMismatches:  []TranslationMismatch{},
			ImageMismatches: []TranslationMismatch{},
		}

		for id, ptEx := range ptMap {
			enEx, ok := enMap[id]
			if !ok {
				report.MissingInEn = append(report.MissingInEn, id)
				continue
			}

			if len(ptEx.Instructions) != len(enEx.Instructions) {
				report.StepMismatches = append(report.StepMismatches, TranslationMismatch{
					ID:       id,
					NamePt:   ptEx.Name,
					NameEn:   enEx.Name,
					Field:    "instructions",
					ValuePt:  fmt.Sprintf("%d passos", len(ptEx.Instructions)),
					ValueEn:  fmt.Sprintf("%d passos", len(enEx.Instructions)),
					Severity: "warning",
					Message:  fmt.Sprintf("Diferença no número de instruções (PT: %d vs EN: %d)", len(ptEx.Instructions), len(enEx.Instructions)),
				})
			}

			if len(ptEx.Images) != len(enEx.Images) {
				report.ImageMismatches = append(report.ImageMismatches, TranslationMismatch{
					ID:       id,
					NamePt:   ptEx.Name,
					NameEn:   enEx.Name,
					Field:    "images",
					ValuePt:  fmt.Sprintf("%d fotos", len(ptEx.Images)),
					ValueEn:  fmt.Sprintf("%d fotos", len(enEx.Images)),
					Severity: "warning",
					Message:  fmt.Sprintf("Diferença na quantidade de fotos (PT: %d vs EN: %d)", len(ptEx.Images), len(enEx.Images)),
				})
			}
		}

		for id := range enMap {
			if _, ok := ptMap[id]; !ok {
				report.MissingInPt = append(report.MissingInPt, id)
			}
		}

		sort.Strings(report.MissingInEn)
		sort.Strings(report.MissingInPt)

		report.TotalDiscrepancies = len(report.MissingInEn) + len(report.MissingInPt) + len(report.StepMismatches) + len(report.ImageMismatches)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(report)
	}
}

func handleTranslationCompare(baseDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if strings.TrimSpace(id) == "" {
			http.Error(w, `{"error": "id parameter is required"}`, http.StatusBadRequest)
			return
		}

		ptFile := filepath.Join(baseDir, "database", "exercises", id, "pt.json")
		enFile := filepath.Join(baseDir, "database", "exercises", id, "en.json")

		result := map[string]interface{}{
			"id":        id,
			"languages": make(map[string]interface{}),
		}

		langsMap := result["languages"].(map[string]interface{})

		if data, err := os.ReadFile(ptFile); err == nil {
			var ptEx Exercise
			if err := json.Unmarshal(data, &ptEx); err == nil {
				langsMap["pt"] = ptEx
			}
		}

		if data, err := os.ReadFile(enFile); err == nil {
			var enEx Exercise
			if err := json.Unmarshal(data, &enEx); err == nil {
				langsMap["en"] = enEx
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
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

		ext := strings.ToLower(filepath.Ext(relativePath))
		if ext != ".jpg" && ext != ".jpeg" {
			relativePath = strings.TrimSuffix(relativePath, filepath.Ext(relativePath)) + ".jpg"
		} else if ext == ".jpeg" {
			relativePath = strings.TrimSuffix(relativePath, ".jpeg") + ".jpg"
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao ler imagem: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		var finalBytes []byte
		img, _, decodeErr := image.Decode(bytes.NewReader(fileBytes))
		if decodeErr == nil {
			var buf bytes.Buffer
			if encodeErr := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); encodeErr == nil {
				finalBytes = buf.Bytes()
			} else {
				finalBytes = fileBytes
			}
		} else {
			finalBytes = fileBytes
		}

		targetFile := filepath.Join(baseDir, "database", "exercises", relativePath)

		if err := os.MkdirAll(filepath.Dir(targetFile), 0755); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao criar diretório: %s"}`, err.Error()), http.StatusInternalServerError)
			return
		}

		if err := os.WriteFile(targetFile, finalBytes, 0644); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"Erro ao salvar imagem: %s"}`, err.Error()), http.StatusInternalServerError)
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

		targetFile := filepath.Join(baseDir, "database", "exercises", relPath)
		if info, err := os.Stat(targetFile); err == nil && !info.IsDir() {
			http.ServeFile(w, r, targetFile)
			return
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
    <li><code>GET /api/translations/audit</code></li>
    <li><code>GET /api/translations/compare</code></li>
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
	mux.HandleFunc("/api/metadata/category", enableCORS(handleMetadataCategory(baseDir)))
	mux.HandleFunc("/api/translate", enableCORS(handleTranslate(baseDir)))
	mux.HandleFunc("/api/workouts", enableCORS(handleWorkouts(baseDir)))
	mux.HandleFunc("/api/translations/audit", enableCORS(handleTranslationAudit(baseDir)))
	mux.HandleFunc("/api/translations/compare", enableCORS(handleTranslationCompare(baseDir)))
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
