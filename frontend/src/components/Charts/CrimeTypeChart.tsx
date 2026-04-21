import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Cell,
} from "recharts";

interface CrimeTypeCount {
  crimeType: string;
  count: number;
}

interface Props {
  data: CrimeTypeCount[];
}

// Colour ramp for bars — subtle palette on dark background.
const BAR_COLORS = [
  "#ef4444", "#f97316", "#eab308", "#22c55e",
  "#06b6d4", "#6366f1", "#a855f7", "#ec4899",
  "#14b8a6", "#f59e0b", "#84cc16", "#10b981",
];

// Wraps long crime type names across two lines instead of truncating.
function CustomYTick({ x, y, payload }: any) {
  const text: string = payload.value;
  const MAX = 18;

  let lines: string[];
  if (text.length <= MAX) {
    lines = [text];
  } else {
    const breakAt = text.lastIndexOf(" ", MAX);
    if (breakAt > 4) {
      const second = text.slice(breakAt + 1);
      lines = [
        text.slice(0, breakAt),
        second.length > MAX ? second.slice(0, MAX - 1) + "…" : second,
      ];
    } else {
      lines = [text.slice(0, MAX - 1) + "…"];
    }
  }

  const lineH = 11;
  const offsetY = lines.length > 1 ? -(lineH / 2) : 0;

  return (
    <g transform={`translate(${x},${y})`}>
      {lines.map((line, i) => (
        <text
          key={i}
          x={-4}
          y={offsetY + i * lineH}
          textAnchor="end"
          dominantBaseline="middle"
          fill="#94a3b8"
          fontSize={9}
        >
          {line}
        </text>
      ))}
    </g>
  );
}

export default function CrimeTypeChart({ data }: Props) {
  const chartHeight = Math.max(300, data.length * 32);

  return (
    <ResponsiveContainer width="100%" height={chartHeight}>
      <BarChart
        data={data}
        layout="vertical"
        margin={{ top: 0, right: 8, left: 4, bottom: 0 }}
        barCategoryGap="20%"
      >
        <XAxis
          type="number"
          tick={{ fill: "#64748b", fontSize: 10 }}
          tickLine={false}
          axisLine={false}
        />
        <YAxis
          type="category"
          dataKey="crimeType"
          width={155}
          tick={<CustomYTick />}
          tickLine={false}
          axisLine={false}
        />
        <Tooltip
          contentStyle={{
            background: "#0f172a",
            border: "1px solid #334155",
            borderRadius: 8,
            fontSize: 12,
          }}
          labelStyle={{ color: "#f1f5f9", fontWeight: 600 }}
          itemStyle={{ color: "#f1f5f9" }}
          formatter={(v: number) => [v.toLocaleString(), "Crimes"]}
        />
        <Bar dataKey="count" radius={[0, 4, 4, 0]}>
          {data.map((_, i) => (
            <Cell key={i} fill={BAR_COLORS[i % BAR_COLORS.length]} />
          ))}
        </Bar>
      </BarChart>
    </ResponsiveContainer>
  );
}
