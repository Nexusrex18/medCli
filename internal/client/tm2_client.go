package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Nexusrex18/medCli/internal/config"
	"github.com/Nexusrex18/medCli/internal/models"
	"github.com/Nexusrex18/medCli/internal/repository"
	"github.com/patrickmn/go-cache"
)

type TM2Client struct {
	repo   *repository.CSVRepository
	config *config.Config
	cache  *cache.Cache
	hits   int
	misses int
}

type SearchResult struct {
	Records []models.MedicineRecord `json:"records"`
	Count   int                     `json:"count"`
}

type SymptomSearchResult struct {
	Records []models.MedicineRecord `json:"records"`
	Count   int                     `json:"count"`
}

func NewTM2Client(cfg *config.Config) (*TM2Client, error) {
	// Load CSV data
	repo, err := repository.NewCSVRepository(cfg.CSV.FilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CSV data: %w", err)
	}

	cacheTTL, err := time.ParseDuration(cfg.Cache.TTL)
	if err != nil {
		return nil, fmt.Errorf("invalid cache TTL format: %w", err)
	}

	return &TM2Client{
		repo:   repo,
		config: cfg,
		cache:  cache.New(cacheTTL, 10*time.Minute),
		hits:   0,
		misses: 0,
	}, nil
}

func (c *TM2Client) SearchByCode(ctx context.Context, code string, searchType string) (*SearchResult, error) {
	cacheKey := "search:" + code + ":" + searchType
	if cached, found := c.cache.Get(cacheKey); found {
		c.hits++
		return cached.(*SearchResult), nil
	}
	c.misses++

	records := c.repo.SearchByCode(code)

	result := &SearchResult{
		Records: records,
		Count:   len(records),
	}

	c.cache.Set(cacheKey, result, cache.DefaultExpiration)
	return result, nil
}

func (c *TM2Client) SearchBySymptoms(ctx context.Context, symptoms []string) (*SymptomSearchResult, error) {
	cacheKey := "symptoms:" + strings.Join(symptoms, ",")
	if cached, found := c.cache.Get(cacheKey); found {
		c.hits++
		return cached.(*SymptomSearchResult), nil
	}
	c.misses++

	records := c.repo.SearchBySymptoms(symptoms)

	result := &SymptomSearchResult{
		Records: records,
		Count:   len(records),
	}

	c.cache.Set(cacheKey, result, cache.DefaultExpiration)
	return result, nil
}

func (c *TM2Client) GetCacheStats() (hits, misses, items int) {
	return c.hits, c.misses, c.cache.ItemCount()
}

func (c *TM2Client) GetRepoStats() map[string]int {
	return c.repo.GetStats()
}