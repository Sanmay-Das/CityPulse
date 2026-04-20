import { useQuery } from "@apollo/client";
import { GET_CRIME_TIMELINE, GET_CRIME_TYPES } from "../../graphql/queries";
import { useMapStore } from "../../store/useMapStore";
import TrendChart from "../Charts/TrendChart";
import CrimeTypeChart from "../Charts/CrimeTypeChart";

interface Props {
  zipCode: string;
  onClose: () => void;
}

export default function DrilldownPanel({ zipCode, onClose }: Props) {
  const { selectedCity, zipStatsMap, selectedYear } = useMapStore();
  const stats = zipStatsMap[zipCode];

  const { data: timelineData, loading: timelineLoading } = useQuery(
    GET_CRIME_TIMELINE,
    { variables: { city: selectedCity, zipCode } }
  );

  const { data: typesData, loading: typesLoading } = useQuery(GET_CRIME_TYPES, {
    variables: { city: selectedCity, zipCode, year: selectedYear || null },
  });

  const timeline = timelineData?.crimeTimeline ?? [];
  const crimeTypes = typesData?.crimeTypes ?? [];

  // Compute a simple trend: compare last 12 months vs previous 12 months.
  const trend = (() => {
    if (timeline.length < 24) return null;
    const recent = timeline.slice(-12).reduce((s: number, m: any) => s + m.count, 0);
    const prior = timeline.slice(-24, -12).reduce((s: number, m: any) => s + m.count, 0);
    if (prior === 0) return null;
    const pct = ((recent - prior) / prior) * 100;
    return { pct, improving: pct < 0 };
  })();

  return (
    <div className="absolute top-12 right-0 bottom-0 w-[30rem] z-10 flex flex-col bg-slate-900/95 backdrop-blur border-l border-slate-700 overflow-y-auto">
      {/* Header */}
      <div className="flex items-start justify-between p-5 border-b border-slate-700">
        <div>
          <p className="text-slate-400 text-xs uppercase tracking-wider mb-1">ZIP Code</p>
          <h2 className="text-white text-2xl font-bold">{zipCode}</h2>
          {stats && (
            <p className="text-slate-400 text-sm mt-1">
              Top crime: <span className="text-slate-300">{stats.topCrimeType}</span>
            </p>
          )}
        </div>
        <button
          onClick={onClose}
          className="text-slate-400 hover:text-white transition-colors p-1"
          aria-label="Close panel"
        >
          <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      {/* Stat cards */}
      {stats && (
        <div className="grid grid-cols-2 gap-3 p-5">
          <StatCard
            label="Total crimes"
            value={stats.totalCrimes.toLocaleString()}
            sub={selectedYear ? String(selectedYear) : "All years"}
          />
          {trend && (
            <StatCard
              label="1-yr trend"
              value={`${trend.improving ? "▼" : "▲"} ${Math.abs(trend.pct).toFixed(1)}%`}
              sub="vs prior year"
              highlight={trend.improving ? "green" : "red"}
            />
          )}
        </div>
      )}

      {/* Crime trend over time */}
      <Section title="Crime trend over time" loading={timelineLoading}>
        {timeline.length > 0 ? (
          <TrendChart data={timeline} />
        ) : (
          <Empty message="No timeline data" />
        )}
      </Section>

      {/* Crime type breakdown */}
      <Section
        title={`Crime types${selectedYear ? ` (${selectedYear})` : ""}`}
        loading={typesLoading}
      >
        {crimeTypes.length > 0 ? (
          <CrimeTypeChart data={crimeTypes} />
        ) : (
          <Empty message="No crime type data" />
        )}
      </Section>

      {/* Export link */}
      <div className="p-5 border-t border-slate-700 mt-auto">
        <a
          href={`/api/export/csv?zip=${zipCode}${selectedYear ? `&year=${selectedYear}` : ""}`}
          download
          className="block w-full text-center text-sm py-2.5 px-4 bg-slate-700 hover:bg-slate-600 text-slate-300 rounded-lg transition-colors"
        >
          Export ZIP data (CSV)
        </a>
      </div>
    </div>
  );
}

// ── Small helper components ───────────────────────────────────────────────────

function StatCard({
  label,
  value,
  sub,
  highlight,
}: {
  label: string;
  value: string;
  sub?: string;
  highlight?: "green" | "red";
}) {
  const valueColor =
    highlight === "green"
      ? "text-emerald-400"
      : highlight === "red"
      ? "text-red-400"
      : "text-white";

  return (
    <div className="bg-slate-800 rounded-xl p-4">
      <p className="text-slate-400 text-xs uppercase tracking-wide mb-1">{label}</p>
      <p className={`text-xl font-bold ${valueColor}`}>{value}</p>
      {sub && <p className="text-slate-500 text-xs mt-0.5">{sub}</p>}
    </div>
  );
}

function Section({
  title,
  loading,
  children,
}: {
  title: string;
  loading: boolean;
  children: React.ReactNode;
}) {
  return (
    <div className="p-5 border-t border-slate-700">
      <h3 className="text-slate-300 text-sm font-semibold mb-3">{title}</h3>
      {loading ? (
        <div className="h-12 flex items-center justify-center">
          <div className="text-slate-500 text-sm animate-pulse">Loading…</div>
        </div>
      ) : (
        children
      )}
    </div>
  );
}

function Empty({ message }: { message: string }) {
  return (
    <div className="h-12 flex items-center justify-center text-slate-500 text-sm">
      {message}
    </div>
  );
}
