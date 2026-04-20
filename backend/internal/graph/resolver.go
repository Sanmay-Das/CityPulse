package graph

import (
	"fmt"

	"github.com/graphql-go/graphql"

	"chicago-safety-index/internal/db"
)

// Resolver holds the database reference shared by all field resolvers.
type Resolver struct {
	DB *db.Database
}

// cityArg extracts the city argument, defaulting to "Chicago".
func cityArg(p graphql.ResolveParams) string {
	if c, ok := p.Args["city"].(string); ok && c != "" {
		return c
	}
	return "Chicago"
}

// allZipStats resolves Query.allZipStats
func (r *Resolver) allZipStats(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	year, _ := p.Args["year"].(int)
	rows, err := r.DB.GetAllZipStats(p.Context, city, year)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rows))
	for i, row := range rows {
		result[i] = map[string]interface{}{
			"city":         row.City,
			"zipCode":      row.ZipCode,
			"totalCrimes":  row.TotalCrimes,
			"topCrimeType": row.TopCrimeType,
		}
	}
	return result, nil
}

// zipStats resolves Query.zipStats
func (r *Resolver) zipStats(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	zipCode, ok := p.Args["zipCode"].(string)
	if !ok || zipCode == "" {
		return nil, fmt.Errorf("zipCode is required")
	}

	row, err := r.DB.GetZipStats(p.Context, city, zipCode)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"city":         row.City,
		"zipCode":      row.ZipCode,
		"totalCrimes":  row.TotalCrimes,
		"topCrimeType": row.TopCrimeType,
	}, nil
}

// crimeTimeline resolves Query.crimeTimeline
func (r *Resolver) crimeTimeline(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	zipCode, ok := p.Args["zipCode"].(string)
	if !ok || zipCode == "" {
		return nil, fmt.Errorf("zipCode is required")
	}

	rows, err := r.DB.GetCrimeTimeline(p.Context, city, zipCode)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rows))
	for i, m := range rows {
		result[i] = map[string]interface{}{
			"year":  m.Year,
			"month": m.Month,
			"count": m.Count,
		}
	}
	return result, nil
}

// crimeTypes resolves Query.crimeTypes
func (r *Resolver) crimeTypes(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	zipCode, ok := p.Args["zipCode"].(string)
	if !ok || zipCode == "" {
		return nil, fmt.Errorf("zipCode is required")
	}
	year, _ := p.Args["year"].(int)

	rows, err := r.DB.GetCrimeTypes(p.Context, city, zipCode, year)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(rows))
	for i, c := range rows {
		result[i] = map[string]interface{}{
			"crimeType": c.CrimeType,
			"count":     c.Count,
		}
	}
	return result, nil
}

// availableYears resolves Query.availableYears
func (r *Resolver) availableYears(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	years, err := r.DB.GetAvailableYears(p.Context, city)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(years))
	for i, y := range years {
		result[i] = y
	}
	return result, nil
}

// yearRange resolves Query.yearRange
func (r *Resolver) yearRange(p graphql.ResolveParams) (interface{}, error) {
	city := cityArg(p)
	yr, err := r.DB.GetYearRange(p.Context, city)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"minYear": yr.MinYear,
		"maxYear": yr.MaxYear,
	}, nil
}

// availableCities resolves Query.availableCities
func (r *Resolver) availableCities(p graphql.ResolveParams) (interface{}, error) {
	cities, err := r.DB.GetAvailableCities(p.Context)
	if err != nil {
		return nil, err
	}
	result := make([]interface{}, len(cities))
	for i, c := range cities {
		result[i] = c
	}
	return result, nil
}
