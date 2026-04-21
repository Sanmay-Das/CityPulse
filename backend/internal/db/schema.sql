-- ZIP code list (no geometry — stored as static GeoJSON file instead)
CREATE TABLE IF NOT EXISTS zip_codes (
    zip_code VARCHAR(10)  NOT NULL,
    city     VARCHAR(100) NOT NULL,
    PRIMARY KEY (zip_code, city)
);

-- Aggregated crime counts by ZIP + city + year + month + crime type
CREATE TABLE IF NOT EXISTS zip_crime_stats (
    zip_code    VARCHAR(10)  NOT NULL,
    city        VARCHAR(100) NOT NULL,
    year        INT          NOT NULL,
    month       INT          NOT NULL,
    crime_type  VARCHAR(100) NOT NULL,
    crime_count INT          NOT NULL DEFAULT 0,
    PRIMARY KEY (zip_code, city, year, month, crime_type)
);

CREATE INDEX IF NOT EXISTS idx_zcs_city         ON zip_crime_stats(city);
CREATE INDEX IF NOT EXISTS idx_zcs_year         ON zip_crime_stats(year);
CREATE INDEX IF NOT EXISTS idx_zcs_city_year    ON zip_crime_stats(city, year);
CREATE INDEX IF NOT EXISTS idx_zcs_zip_city     ON zip_crime_stats(zip_code, city);
