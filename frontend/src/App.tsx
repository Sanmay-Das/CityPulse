import { useEffect } from "react";
import { useQuery } from "@apollo/client";
import { GET_ALL_ZIP_STATS, GET_AVAILABLE_YEARS, GET_AVAILABLE_CITIES, GET_YEAR_RANGE } from "./graphql/queries";
import { useMapStore } from "./store/useMapStore";
import MapView from "./components/Map/MapView";
import DrilldownPanel from "./components/Panel/DrilldownPanel";

export default function App() {
  const { selectedCity, selectedYear, selectedZip, setSelectedCity, setSelectedYear, setZipStats, setSelectedZip } =
    useMapStore();

  const { data: citiesData } = useQuery(GET_AVAILABLE_CITIES);
  const { data: yearsData } = useQuery(GET_AVAILABLE_YEARS, {
    variables: { city: selectedCity },
  });
  const { data: yearRangeData } = useQuery(GET_YEAR_RANGE, {
    variables: { city: selectedCity },
  });

  const { data, loading, error } = useQuery(GET_ALL_ZIP_STATS, {
    variables: { city: selectedCity, year: selectedYear || null },
  });

  useEffect(() => {
    if (data?.allZipStats) {
      setZipStats(data.allZipStats);
    }
  }, [data, setZipStats]);

  const cities: string[] = citiesData?.availableCities ?? [];
  const years: number[] = yearsData?.availableYears ?? [];
  const minYear: number = yearRangeData?.yearRange?.minYear ?? 0;
  const maxYear: number = yearRangeData?.yearRange?.maxYear ?? 0;
  const yearRangeLabel = minYear && maxYear ? `${minYear}–${maxYear}` : "…";

  return (
    <div className="relative w-full h-full flex flex-col">
      {/* ── Header ─────────────────────────────────────────────────────────── */}
      <header className="absolute top-0 left-0 right-0 z-10 flex items-center justify-between px-5 py-3 bg-slate-900/90 backdrop-blur border-b border-slate-700">
        <div className="flex items-center gap-3">
          <div className="w-2 h-6 bg-chicago-red rounded-full" />
          <h1 className="text-white font-bold text-lg tracking-tight">CityPulse</h1>
          <span className="text-slate-400 text-sm hidden sm:block">
            Crime by ZIP Code · {yearRangeLabel}
          </span>
        </div>

        <div className="flex items-center gap-3">
          {/* City switcher */}
          <select
            value={selectedCity}
            onChange={(e) => setSelectedCity(e.target.value)}
            className="bg-slate-800 border border-slate-600 text-slate-200 text-sm rounded-lg px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-chicago-red"
          >
            {cities.length > 0 ? (
              cities.map((c) => (
                <option key={c} value={c}>{c}</option>
              ))
            ) : (
              <option value="Chicago">Chicago</option>
            )}
          </select>

          {/* Year filter */}
          <select
            value={selectedYear}
            onChange={(e) => setSelectedYear(Number(e.target.value))}
            className="bg-slate-800 border border-slate-600 text-slate-200 text-sm rounded-lg px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-chicago-red"
          >
            <option value={0}>Available years</option>
            {years.map((y) => (
              <option key={y} value={y}>{y}</option>
            ))}
          </select>

          {/* Export */}
          <a
            href={`/api/export/csv?city=${selectedCity}${selectedYear ? `&year=${selectedYear}` : ""}`}
            download
            className="text-xs px-3 py-1.5 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg transition-colors"
          >
            Export CSV
          </a>
        </div>
      </header>

      {loading && (
        <div className="absolute inset-0 z-20 flex items-center justify-center bg-slate-900/60 backdrop-blur-sm">
          <div className="text-slate-300 text-sm animate-pulse">Loading crime data…</div>
        </div>
      )}
      {error && (
        <div className="absolute inset-0 z-20 flex items-center justify-center bg-slate-900/80">
          <div className="bg-red-900/80 border border-red-700 rounded-xl p-6 max-w-sm text-center">
            <p className="text-red-300 font-semibold">Failed to load data</p>
            <p className="text-red-400 text-sm mt-1">{error.message}</p>
          </div>
        </div>
      )}

      <MapView />

      {selectedZip && (
        <DrilldownPanel zipCode={selectedZip} onClose={() => setSelectedZip(null)} />
      )}

      {/* Legend */}
      <div className="absolute bottom-6 left-5 z-10 bg-slate-900/90 backdrop-blur border border-slate-700 rounded-xl p-4">
        <p className="text-slate-400 text-xs font-medium mb-2 uppercase tracking-wide">
          Crime density
        </p>
        <div
          className="h-3 w-32 rounded"
          style={{ background: "linear-gradient(to right, #fef9c3, #fbbf24, #ef4444, #7f1d1d)" }}
        />
        <div className="flex justify-between text-xs text-slate-500 mt-1">
          <span>Low</span>
          <span>High</span>
        </div>
      </div>
    </div>
  );
}
