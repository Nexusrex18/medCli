package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Nexusrex18/medCli/internal/models"
)

type CSVRepository struct {
	records      []models.MedicineRecord
	codeIndex    map[string][]models.MedicineRecord // code -> records
	tm2CodeIndex map[string][]models.MedicineRecord // tm2_code -> records
	mu           sync.RWMutex
}

func NewCSVRepository(csvFilePath string) (*CSVRepository, error) {
	repo := &CSVRepository{
		codeIndex:    make(map[string][]models.MedicineRecord),
		tm2CodeIndex: make(map[string][]models.MedicineRecord),
	}

	if err := repo.loadCSV(csvFilePath); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *CSVRepository) loadCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or has only headers")
	}

	headers := records[0]
	var dataRecords []models.MedicineRecord

	for i, record := range records[1:] {
		if len(record) != len(headers) {
			return fmt.Errorf("record %d has incorrect number of fields", i+1)
		}

		medicine := models.MedicineRecord{}
		for j, header := range headers {
			value := record[j]
			switch strings.ToLower(header) {
			case "tm2_code":
				medicine.TM2Code = value
			case "code":
				medicine.Code = value
			case "tm2_title":
				medicine.TM2Title = value
			case "tm2_definition":
				medicine.TM2Definition = value
			case "code_title":
				medicine.CodeTitle = value
			case "code_description":
				medicine.Description = value
			case "confidence_score":
				if score, err := strconv.ParseFloat(value, 64); err == nil {
					medicine.ConfidenceScore = score
				}
			case "type":
				medicine.Type = value
			case "tm2_link":
				medicine.TM2Link = value
			}
		}
		dataRecords = append(dataRecords, medicine)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = dataRecords
	r.buildIndexes()

	return nil
}

func (r *CSVRepository) buildIndexes() {
	for _, record := range r.records {
		// Index by traditional code (lowercase for case-insensitive search)
		codeKey := strings.ToLower(record.Code)
		r.codeIndex[codeKey] = append(r.codeIndex[codeKey], record)

		// Index by TM2 code (lowercase for case-insensitive search)
		tm2Key := strings.ToLower(record.TM2Code)
		r.tm2CodeIndex[tm2Key] = append(r.tm2CodeIndex[tm2Key], record)
	}
}

func (r *CSVRepository) SearchByCode(code string) []models.MedicineRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	code = strings.ToLower(strings.TrimSpace(code))
	var results []models.MedicineRecord

	// Search in both code indexes
	if records, exists := r.codeIndex[code]; exists {
		results = append(results, records...)
	}
	if records, exists := r.tm2CodeIndex[code]; exists {
		results = append(results, records...)
	}

	return results
}

func (r *CSVRepository) SearchBySymptoms(symptoms []string) []models.MedicineRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []models.MedicineRecord
	seen := make(map[string]bool) // To avoid duplicates

	for _, record := range r.records {
		// Combine all searchable fields into one string for this record
		searchableText := strings.ToLower(
			record.TM2Title + " " +
				record.Description + " " +
				record.TM2Definition + " " +
				record.CodeTitle,
		)

		// Check if this record matches ALL symptoms (AND logic)
		recordMatchesAll := true

		for _, symptom := range symptoms {
			symptom = strings.ToLower(strings.TrimSpace(symptom))
			if symptom == "" {
				continue
			}

			// Split symptom into individual words
			words := strings.Fields(symptom)
			if len(words) == 0 {
				continue
			}

			// Check if ALL words in this symptom exist in the searchable text
			allWordsFound := true
			for _, word := range words {
				if len(word) < 2 { // Skip very short words
					continue
				}
				if !strings.Contains(searchableText, word) {
					allWordsFound = false
					break
				}
			}

			// If any symptom doesn't match all its words, this record fails
			if !allWordsFound {
				recordMatchesAll = false
				break
			}
		}

		// If record matches all symptoms, add it to results
		if recordMatchesAll && len(symptoms) > 0 {
			key := record.TM2Code + ":" + record.Code
			if !seen[key] {
				results = append(results, record)
				seen[key] = true
			}
		}
	}

	return results
}

func (r *CSVRepository) GetAllRecords() []models.MedicineRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.records
}

func (r *CSVRepository) GetStats() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return map[string]int{
		"total_records":    len(r.records),
		"unique_codes":     len(r.codeIndex),
		"unique_tm2_codes": len(r.tm2CodeIndex),
	}
}