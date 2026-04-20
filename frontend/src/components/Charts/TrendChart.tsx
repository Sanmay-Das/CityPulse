import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ReferenceLine,
} from "recharts";

interface MonthlyCount {
  year: number;
  month: number;
  count: number;
}

interface Props {
  data: MonthlyCount[];
}

const MONTH_SHORT = ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"];

export default function TrendChart({ data }: Props) {
  // Format as "Jan 2023" labels, show every 12th tick (yearly).
  const chartData = data.map((d) => ({
    label: `${MONTH_SHORT[d.month - 1]} ${d.year}`,
    count: d.count,
    year: d.year,
    month: d.month,
  }));

  // Only show January ticks (one per year)
  const yearBoundaries = chartData
    .filter((d) => d.month === 1)
    .map((d) => d.label);

  return (
    <ResponsiveContainer width="100%" height={200}>
      <LineChart data={chartData} margin={{ top: 4, right: 8, left: -20, bottom: 0 }}>
        <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" />
        {yearBoundaries.map((label) => (
          <ReferenceLine key={label} x={label} stroke="#334155" strokeDasharray="4 2" />
        ))}
        <XAxis
          dataKey="label"
          tick={{ fill: "#64748b", fontSize: 10 }}
          tickLine={false}
          ticks={yearBoundaries.filter((_, i) => i % 2 === 0)}
          tickFormatter={(val) => val.split(" ")[1]} // show only the year number
        />
        <YAxis
          tick={{ fill: "#64748b", fontSize: 10 }}
          tickLine={false}
          axisLine={false}
        />
        <Tooltip
          contentStyle={{
            background: "#0f172a",
            border: "1px solid #334155",
            borderRadius: 8,
            color: "#f1f5f9",
            fontSize: 12,
          }}
          labelStyle={{ color: "#94a3b8" }}
          formatter={(v: number) => [v.toLocaleString(), "Crimes"]}
        />
        <Line
          type="monotone"
          dataKey="count"
          stroke="#ef4444"
          strokeWidth={2}
          dot={false}
          activeDot={{ r: 4, fill: "#ef4444" }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
}
