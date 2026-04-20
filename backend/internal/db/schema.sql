-- Enable PostGIS
CREATE EXTENSION IF NOT EXISTS postgis;

-- ZIP code list (no geometry — stored as static GeoJSON file instead)
CREATE TABLE IF NOT EXISTS zip_codes (
    zip_code VARCHAR(10) PRIMARY KEY
);

-- Aggregated crime counts by ZIP + year + month + crime type
CREATE TABLE IF NOT EXISTS zip_crime_stats (
    zip_code    VARCHAR(10) NOT NULL REFERENCES zip_codes(zip_code),
    year        INT         NOT NULL,
    month       INT         NOT NULL,
    crime_type  VARCHAR(100) NOT NULL,
    crime_count INT         NOT NULL DEFAULT 0,
    PRIMARY KEY (zip_code, year, month, crime_type)
);

CREATE INDEX IF NOT EXISTS idx_zcs_zip      ON zip_crime_stats(zip_code);
CREATE INDEX IF NOT EXISTS idx_zcs_year     ON zip_crime_stats(year);
CREATE INDEX IF NOT EXISTS idx_zcs_zip_year ON zip_crime_stats(zip_code, year);
