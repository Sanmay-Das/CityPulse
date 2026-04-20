# Chicago Safety Index — Setup Guide

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.22+ | https://go.dev/dl |
| Node.js | 18+ | https://nodejs.org |
| Docker Desktop | any | https://www.docker.com/products/docker-desktop |

---

## Step 1 — Download the datasets

### Chicago Crimes CSV
1. Go to https://star.cs.ucr.edu/?Chicago%20Crimes
2. Click **Download** → save as `data/Chicago_Crimes.csv`

### ZIP Code Boundaries (TIGER 2018 ZCTA5)
1. Go to https://star.cs.ucr.edu/?TIGER2018/ZCTA5
2. Zoom to Chicago area, click **Download** → extract the `.geojson` file
3. Save as `data/ZCTA5_Chicago.geojson`

Create the data directory:
```
mkdir data
```

---

## Step 2 — Start PostgreSQL (PostGIS)

```bash
docker compose up -d
```

Wait ~10 seconds for the container to be healthy, then create the schema:

```bash
docker exec -i chicago_safety_db psql -U chicago -d chicago_safety \
  < backend/internal/db/schema.sql
```

---

## Step 3 — Run the ETL (one-time)

```bash
cd backend
go mod tidy
go run ./cmd/etl \
  -crimes  ../data/Chicago_Crimes.csv \
  -zipgeo  ../data/ZCTA5_Chicago.geojson
```

This reads ~8 million rows, does the spatial join in memory, and bulk-inserts
aggregated counts. Expect 3-8 minutes depending on your machine.

---

## Step 4 — Start the Go API server

```bash
# still inside backend/
go run ./cmd/server
```

Server starts on http://localhost:8080
GraphiQL playground: http://localhost:8080/graphql

---

## Step 5 — Start the React frontend

```bash
cd frontend
npm install
npm run dev
```

App opens at http://localhost:5173

---

## Step 6 — Connect PowerBI

### Option A — CSV (recommended)
1. Open PowerBI Desktop
2. **Get Data → Web**
3. URL: `http://localhost:8080/api/export/csv`  
   (add `?year=2023` to filter by year)
4. PowerBI parses the CSV automatically.

### Option B — JSON
Same steps, URL: `http://localhost:8080/api/export/json`

### Suggested PowerBI visuals
- **Clustered bar chart**: ZIP code (X) vs crime_count (Y), filtered by crime_type
- **Line chart**: year/month vs crime_count for trend analysis
- **Map visual**: ZIP code with bubble size = total crimes
- **Slicer**: year and crime_type filters

---

## Project structure

```
chicago-safety-index/
├── docker-compose.yml          # PostgreSQL + PostGIS
├── data/                       # Downloaded datasets (gitignored)
│   ├── Chicago_Crimes.csv
│   └── ZCTA5_Chicago.geojson
├── backend/
│   ├── go.mod
│   ├── cmd/
│   │   ├── etl/main.go         # One-time ETL runner
│   │   └── server/main.go      # GraphQL + REST API server
│   └── internal/
│       ├── db/
│       │   ├── db.go           # DB connection + all queries
│       │   └── schema.sql      # Table definitions
│       ├── etl/
│       │   └── loader.go       # Spatial join + bulk insert
│       └── graph/
│           ├── schema.go       # GraphQL type definitions
│           └── resolver.go     # GraphQL resolvers
└── frontend/
    ├── package.json
    ├── vite.config.ts          # Dev proxy → Go backend
    └── src/
        ├── App.tsx             # Root: year filter, header, layout
        ├── store/useMapStore.ts        # Zustand global state
        ├── graphql/queries.ts          # Apollo gql queries
        └── components/
            ├── Map/MapView.tsx         # DeckGL choropleth map
            ├── Panel/DrilldownPanel.tsx # ZIP drilldown sidebar
            └── Charts/
                ├── TrendChart.tsx      # Monthly line chart
                └── CrimeTypeChart.tsx  # Crime type bar chart
```

---

## API reference

### GraphQL (POST /graphql)

```graphql
# All ZIP stats for a year (or all years if year omitted)
query {
  allZipStats(year: 2023) {
    zipCode
    geometry     # GeoJSON geometry string
    totalCrimes
    topCrimeType
  }
}

# Monthly trend for one ZIP
query {
  crimeTimeline(zipCode: "60601") {
    year month count
  }
}

# Crime type breakdown
query {
  crimeTypes(zipCode: "60601", year: 2023) {
    crimeType count
  }
}

# Available years in dataset
query {
  availableYears
}
```

### REST
| Endpoint | Description |
|----------|-------------|
| `GET /api/health` | Health check |
| `GET /api/export/csv?year=2023` | CSV export for PowerBI |
| `GET /api/export/json?year=2023` | JSON export |
