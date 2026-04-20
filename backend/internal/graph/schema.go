// Package graph defines the GraphQL schema and wires resolvers.
// Uses github.com/graphql-go/graphql — no code-generation required.
package graph

import (
	"github.com/graphql-go/graphql"
)

// ── Scalar types ─────────────────────────────────────────────────────────────

var monthlyCountType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "MonthlyCount",
	Description: "Total crime count for a specific year/month.",
	Fields: graphql.Fields{
		"year":  &graphql.Field{Type: graphql.Int},
		"month": &graphql.Field{Type: graphql.Int},
		"count": &graphql.Field{Type: graphql.Int},
	},
})

var crimeTypeCountType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "CrimeTypeCount",
	Description: "Total count for one crime type.",
	Fields: graphql.Fields{
		"crimeType": &graphql.Field{Type: graphql.String},
		"count":     &graphql.Field{Type: graphql.Int},
	},
})

var zipStatsType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ZipStats",
	Description: "Aggregated crime statistics for a single ZIP code.",
	Fields: graphql.Fields{
		"city": &graphql.Field{
			Type:        graphql.String,
			Description: "City name.",
		},
		"zipCode": &graphql.Field{
			Type:        graphql.String,
			Description: "5-digit ZIP code.",
		},
		"totalCrimes": &graphql.Field{
			Type:        graphql.Int,
			Description: "Total crime count (all types, filtered by year if provided).",
		},
		"topCrimeType": &graphql.Field{
			Type:        graphql.String,
			Description: "Most frequent crime type.",
		},
	},
})

var yearRangeType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "YearRange",
	Description: "Min and max year of data available for a city.",
	Fields: graphql.Fields{
		"minYear": &graphql.Field{Type: graphql.Int},
		"maxYear": &graphql.Field{Type: graphql.Int},
	},
})

// ── Root query ────────────────────────────────────────────────────────────────

// NewSchema builds and returns the executable GraphQL schema.
func NewSchema(r *Resolver) (graphql.Schema, error) {
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{

			// allZipStats(?city, ?year) → [ZipStats]
			"allZipStats": &graphql.Field{
				Type:        graphql.NewList(zipStatsType),
				Description: "Return aggregated crime stats for every ZIP code. Pass `city` and/or `year` to filter.",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
						Description:  "City name (e.g. \"Chicago\"). Defaults to \"Chicago\".",
					},
					"year": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
						Description:  "Calendar year (e.g. 2023). Omit or pass 0 for all years.",
					},
				},
				Resolve: r.allZipStats,
			},

			// zipStats(city, zipCode!) → ZipStats
			"zipStats": &graphql.Field{
				Type:        zipStatsType,
				Description: "Return stats for a single ZIP code (all years combined).",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
						Description:  "City name. Defaults to \"Chicago\".",
					},
					"zipCode": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: r.zipStats,
			},

			// crimeTimeline(city, zipCode!) → [MonthlyCount]
			"crimeTimeline": &graphql.Field{
				Type:        graphql.NewList(monthlyCountType),
				Description: "Monthly crime totals for a ZIP code (all time).",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
						Description:  "City name. Defaults to \"Chicago\".",
					},
					"zipCode": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: r.crimeTimeline,
			},

			// crimeTypes(city, zipCode!, ?year) → [CrimeTypeCount]
			"crimeTypes": &graphql.Field{
				Type:        graphql.NewList(crimeTypeCountType),
				Description: "Top crime types for a ZIP code.",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
						Description:  "City name. Defaults to \"Chicago\".",
					},
					"zipCode": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"year": &graphql.ArgumentConfig{
						Type:         graphql.Int,
						DefaultValue: 0,
					},
				},
				Resolve: r.crimeTypes,
			},

			// availableYears(city) → [Int]
			"availableYears": &graphql.Field{
				Type:        graphql.NewList(graphql.Int),
				Description: "List of years for which data is available for the given city.",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
					},
				},
				Resolve: r.availableYears,
			},

			// availableCities → [String]
			"availableCities": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "List of cities for which data is available.",
				Resolve:     r.availableCities,
			},

			// yearRange(city) → YearRange
			"yearRange": &graphql.Field{
				Type:        yearRangeType,
				Description: "Min and max year of data available for a city.",
				Args: graphql.FieldConfigArgument{
					"city": &graphql.ArgumentConfig{
						Type:         graphql.String,
						DefaultValue: "Chicago",
					},
				},
				Resolve: r.yearRange,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: rootQuery})
}
