import { Link } from "react-router-dom";

export default function Insights() {
  return (
    <div className="h-full overflow-y-auto bg-slate-950 text-slate-300">
      {/* Header */}
      <header className="bg-slate-900/90 border-b border-slate-700 px-6 py-3 flex items-center gap-3">
        <div className="w-2 h-6 bg-chicago-red rounded-full" />
        <span className="text-white font-bold text-lg tracking-tight">CityCrimes</span>
        <span className="text-slate-500 text-sm">/ Insights</span>
      </header>

      <main className="max-w-3xl mx-auto px-6 py-12">
        <h1 className="text-white text-3xl font-bold mb-2">Crime Insights</h1>
        <p className="text-slate-400 text-sm mb-6">
          Data-driven analysis of crime trends across Chicago, Los Angeles, and San Francisco —
          sourced from 10.8M+ official crime records (2001–2026).
        </p>

        {/* CTA */}
        <div className="bg-slate-900 border border-slate-700 rounded-xl p-6 text-center mb-12">
          <h3 className="text-white font-bold text-lg mb-2">Explore the Data Yourself</h3>
          <p className="text-slate-400 text-sm mb-4">
            CityCrimes gives you interactive access to 10.8M+ crime records across 261 ZIP codes.
            Click any ZIP code on the map for a detailed breakdown, or use the AI query box to ask
            natural language questions about the data.
          </p>
          <Link
            to="/map"
            className="inline-block bg-chicago-red hover:bg-red-700 text-white text-sm font-semibold px-6 py-2 rounded-lg transition-colors"
          >
            Open the Map
          </Link>
        </div>

        {/* Article 1 */}
        <article className="mb-14">
          <span className="text-chicago-red text-xs font-semibold uppercase tracking-widest">Chicago</span>
          <h2 className="text-white text-2xl font-bold mt-1 mb-3">
            Top 10 Lowest-Crime ZIP Codes in Chicago
          </h2>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            Not all of Chicago is created equal when it comes to crime. Using official Chicago Police
            Department data spanning over two decades, we identified the ZIP codes with the lowest
            total recorded crime counts. These areas — many of which border suburban Cook County —
            consistently report far fewer incidents than the city average.
          </p>

          <div className="bg-slate-900 border border-slate-700 rounded-xl overflow-hidden mb-4">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-700">
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">Rank</th>
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">ZIP Code</th>
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">Total Crimes (All Years)</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { zip: "60007", total: 5 },
                  { zip: "60076", total: 14 },
                  { zip: "60077", total: 18 },
                  { zip: "46320", total: 30 },
                  { zip: "60419", total: 38 },
                  { zip: "60176", total: 39 },
                  { zip: "60171", total: 54 },
                  { zip: "60106", total: 100 },
                  { zip: "60406", total: 104 },
                  { zip: "60453", total: 171 },
                ].map((row, i) => (
                  <tr key={row.zip} className="border-b border-slate-800 hover:bg-slate-800/50 transition-colors">
                    <td className="px-4 py-3 text-slate-500">#{i + 1}</td>
                    <td className="px-4 py-3 text-white font-medium">{row.zip}</td>
                    <td className="px-4 py-3 text-slate-300">{row.total.toLocaleString()}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            ZIP code 60007 (Elk Grove Village) recorded just 5 crimes across the entire dataset —
            the lowest of any ZIP in the Chicago area. ZIP codes 60076 and 60077, both covering
            parts of Skokie, similarly show very low crime counts, reflecting the relatively safer
            nature of Chicago's northern suburbs compared to downtown or south side neighborhoods.
          </p>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            It's worth noting that lower crime counts in border ZIP codes can also reflect partial
            data coverage — some of these ZIPs straddle city and county lines, meaning not all
            incidents are captured in the city's open data portal.
          </p>
          <Link
            to="/map?city=Chicago"
            className="inline-block text-chicago-red hover:underline text-sm"
          >
            → Explore Chicago crime map
          </Link>
        </article>

        {/* Article 2 */}
        <article className="mb-14">
          <span className="text-chicago-red text-xs font-semibold uppercase tracking-widest">Chicago & Los Angeles</span>
          <h2 className="text-white text-2xl font-bold mt-1 mb-3">
            Most Common Crime Types: Chicago vs. Los Angeles
          </h2>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            The nature of crime differs significantly between Chicago and Los Angeles. Chicago's
            crime landscape is dominated by property crime and violent offenses, while Los Angeles
            sees a distinct pattern driven by vehicle-related crimes — a reflection of the city's
            car-dependent culture.
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
            {/* Chicago */}
            <div className="bg-slate-900 border border-slate-700 rounded-xl p-4">
              <h3 className="text-white font-semibold text-sm mb-3">Chicago — Top Crime Types</h3>
              <div className="space-y-2">
                {[
                  { type: "Theft", total: 1458346 },
                  { type: "Battery", total: 1235673 },
                  { type: "Criminal Damage", total: 796155 },
                  { type: "Narcotics", total: 709130 },
                  { type: "Other Offense", total: 436170 },
                ].map((item) => (
                  <div key={item.type} className="flex justify-between items-center">
                    <span className="text-slate-400 text-xs">{item.type}</span>
                    <span className="text-white text-xs font-medium">{item.total.toLocaleString()}</span>
                  </div>
                ))}
              </div>
            </div>

            {/* LA */}
            <div className="bg-slate-900 border border-slate-700 rounded-xl p-4">
              <h3 className="text-white font-semibold text-sm mb-3">Los Angeles — Top Crime Types</h3>
              <div className="space-y-2">
                {[
                  { type: "Vehicle Stolen", total: 115120 },
                  { type: "Battery - Simple Assault", total: 74464 },
                  { type: "Burglary From Vehicle", total: 63435 },
                  { type: "Theft of Identity", total: 62511 },
                  { type: "Vandalism - Felony", total: 61025 },
                ].map((item) => (
                  <div key={item.type} className="flex justify-between items-center">
                    <span className="text-slate-400 text-xs">{item.type}</span>
                    <span className="text-white text-xs font-medium">{item.total.toLocaleString()}</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            In Chicago, <strong className="text-slate-300">Theft</strong> leads with over 1.45 million
            recorded incidents, followed closely by <strong className="text-slate-300">Battery</strong>
            at 1.23 million. The high narcotics count (709,130) reflects years of targeted enforcement
            in specific neighborhoods. Together, the top 5 crime types account for the majority of
            all recorded crime in the city.
          </p>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            Los Angeles tells a different story. <strong className="text-slate-300">Vehicle theft</strong>
            {" "}tops the list at 115,120 incidents — not surprising given that LA is one of the most
            car-dependent cities in the United States. <strong className="text-slate-300">Burglary
            from vehicle</strong> (63,435) and <strong className="text-slate-300">theft of
            identity</strong> (62,511) round out a crime profile that is heavily skewed toward
            property and financial crimes rather than violent offenses.
          </p>
          <Link
            to="/map?city=Los Angeles"
            className="inline-block text-chicago-red hover:underline text-sm"
          >
            → Explore Los Angeles crime map
          </Link>
        </article>

        {/* Article 3 */}
        <article className="mb-14">
          <span className="text-chicago-red text-xs font-semibold uppercase tracking-widest">San Francisco</span>
          <h2 className="text-white text-2xl font-bold mt-1 mb-3">
            San Francisco Crime Trends: 2003–2025
          </h2>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            San Francisco has seen dramatic swings in crime over the past two decades. From a
            relatively stable period in the early 2000s to a sharp drop during the COVID-19
            pandemic, and a gradual decline in recent years, the data reveals how public health,
            policy, and economic conditions shape crime patterns.
          </p>

          <div className="bg-slate-900 border border-slate-700 rounded-xl overflow-hidden mb-4">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-slate-700">
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">Year</th>
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">Total Crimes</th>
                  <th className="text-left px-4 py-3 text-slate-400 font-medium">Notable</th>
                </tr>
              </thead>
              <tbody>
                {[
                  { year: 2003, total: 140239, note: "Baseline" },
                  { year: 2005, total: 133938, note: "" },
                  { year: 2010, total: 123817, note: "Decade low" },
                  { year: 2015, total: 146646, note: "Rising" },
                  { year: 2018, total: 183279, note: "Peak year" },
                  { year: 2019, total: 134689, note: "" },
                  { year: 2020, total: 109231, note: "COVID-19 drop" },
                  { year: 2021, total: 118697, note: "Recovery" },
                  { year: 2022, total: 124735, note: "" },
                  { year: 2023, total: 122820, note: "" },
                  { year: 2024, total: 103626, note: "Declining" },
                  { year: 2025, total: 90555, note: "Recent low" },
                ].map((row) => (
                  <tr key={row.year} className="border-b border-slate-800 hover:bg-slate-800/50 transition-colors">
                    <td className="px-4 py-3 text-white font-medium">{row.year}</td>
                    <td className="px-4 py-3 text-slate-300">{row.total.toLocaleString()}</td>
                    <td className="px-4 py-3 text-slate-500 text-xs">{row.note}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            San Francisco recorded its highest crime count in <strong className="text-slate-300">2018</strong>
            {" "}with 183,279 incidents. This was followed by a sharp decline in
            {" "}<strong className="text-slate-300">2020</strong> — dropping to 109,231 — a pattern
            seen across most major US cities during the COVID-19 pandemic, when lockdowns
            significantly reduced foot traffic and public activity.
          </p>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            Since 2022, San Francisco has seen a steady downward trend, reaching a recent low of
            {" "}<strong className="text-slate-300">90,555 in 2025</strong> — a 36% reduction from
            the 2018 peak. Whether this reflects improved policing, demographic shifts, or changes
            in reporting practices remains a subject of ongoing debate among city officials and
            researchers.
          </p>
          <p className="text-slate-400 text-sm leading-relaxed mb-4">
            The decade low of 123,817 in 2010 was part of a broader national trend of declining
            crime rates following the 2008 financial crisis, as cities increased community
            policing efforts and expanded social services.
          </p>
          <Link
            to="/map?city=San Francisco"
            className="inline-block text-chicago-red hover:underline text-sm"
          >
            → Explore San Francisco crime map
          </Link>
        </article>

        <div className="mt-12 pt-6 border-t border-slate-800 flex justify-between items-center">
          <Link to="/map" className="text-slate-400 hover:text-white text-sm transition-colors">
            ← Back to Map
          </Link>
          <Link to="/privacy" className="text-slate-400 hover:text-white text-sm transition-colors">
            Privacy Policy
          </Link>
        </div>
      </main>
    </div>
  );
}
