# CityCrimes

A multi-city crime analytics platform that visualizes crime density by ZIP code on an interactive map.

<img width="1919" height="963" alt="Screenshot 2026-04-21 045616" src="https://github.com/user-attachments/assets/298766f9-feaf-4b04-8c87-966be1ef9e0a" />



Currently covers **Chicago**, **Los Angeles** and **San Francisco**.

---

## Stack

| Layer | Technology |
|---|---|
| Frontend | React + Vite + TypeScript |
| Map | react-map-gl + MapLibre GL + Deck.gl |
| Charts | Recharts + D3 color scale |
| State | Zustand + Apollo Client |
| API | GraphQL |
| Database | PostgreSQL |
| ETL | Go - spatial join via R-tree (paulmach/orb + tidwall/rtree) |
| Styling | Tailwind CSS |
| Reporting | PowerBI Desktop via CSV export endpoint |

---

## Features

- **Choropleth map** - ZIP codes colored yellow -> red by crime density
- **City switcher** - switch between Chicago and Los Angeles; map re-centers automatically
- **Year filter** - filter by any year available for the selected city
- **Drilldown panel** - click any ZIP to see:
  - Total crimes + 1-year trend indicator
  - Monthly crime trend line chart
  - Top crime types bar chart
- **Dynamic year range** - header shows the actual data range for the selected city
- **CSV export** - export raw aggregated data per city/year for PowerBI or Excel
- **PowerBI integration** - connect via the `/api/export/csv` REST endpoint

---

## Project Structure

```
citypulse/
├── backend/
│   ├── cmd/
│   │   ├── server/        # GraphQL + REST API server
│   │   └── etl/           # Load crimes CSV -> PostgreSQL
│   └── internal/
│       ├── db/            # PostgreSQL queries
│       ├── etl/           # Spatial join, aggregation, bulk insert
│       └── graph/         # GraphQL schema
├── frontend/
│   ├── public/            # Static GeoJSON boundary files (per city)
│   └── src/
│       ├── components/
│       │   ├── Map/       # MapView (choropleth)
│       │   ├── Panel/     # Drilldown panel
│       │   └── Charts/    # TrendChart, CrimeTypeChart
│       ├── graphql/       # Apollo queries
│       └── store/         # Zustand global state
├── data/                  # Raw CSVs and GeoJSON
└── docker-compose.yml     # PostgreSQL database
```

---

## Getting Started

### Prerequisites

- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)

### 1. Start the database

```bash
docker compose up -d
```

### 2. Create the schema

```bash
Get-Content backend/internal/db/schema.sql | docker exec -i chicago_safety_db psql -U chicago -d chicago_safety
```

### 3. Run the ETL

**Chicago:**
```bash
cd backend
go run ./cmd/etl -crimes ../data/Chicago_Crimes.csv -zipgeo ../data/TIGER2018_ZCTA5.geojson -city Chicago -format chicago
```

**Los Angeles:**
```bash
cd backend
go run ./cmd/etl -crimes ../data/LA_Crimes.csv -zipgeo ../data/TIGER2018_ZCTA5_LA.geojson -city "Los Angeles" -format la
```

### 4. Start the backend

```bash
cd backend
go run ./cmd/server
# Runs on http://localhost:8080
# GraphiQL playground: http://localhost:8080/graphql
```

### 5. Start the frontend

```bash
cd frontend
npm install
npm run dev
# Runs on http://localhost:5173
```

---

## API

### GraphQL - `POST /graphql`

```graphql
# All ZIP stats for a city/year
query {
  allZipStats(city: "Chicago", year: 2019) {
    zipCode
    totalCrimes
    topCrimeType
  }
}

# Monthly trend for a ZIP code
query {
  crimeTimeline(city: "Chicago", zipCode: "60614") {
    year
    month
    count
  }
}

# Available years for a city
query {
  availableYears(city: "Los Angeles")
}

# Data year range for a city
query {
  yearRange(city: "Chicago") {
    minYear
    maxYear
  }
}
```

### REST

| Endpoint | Description |
|---|---|
| `GET /api/export/csv?city=Chicago&year=2019` | Download aggregated data as CSV |
| `GET /api/export/json?city=Chicago` | Same data as JSON |

---

## PowerBI Integration

1. Open PowerBI Desktop
2. **Get Data -> Web**
3. Enter: `http://localhost:8080/api/export/csv?city=Chicago`
4. Load and build reports on top of the aggregated ZIP-level crime data
