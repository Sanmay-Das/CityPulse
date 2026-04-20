import { useMemo, useState, useEffect } from "react";
import Map from "react-map-gl/maplibre";
import DeckGL from "@deck.gl/react";
import { GeoJsonLayer } from "@deck.gl/layers";
import type { PickingInfo } from "@deck.gl/core";
import { scaleSequential } from "d3-scale";
import { interpolateYlOrRd } from "d3-scale-chromatic";
import { useMapStore, CITY_CENTERS } from "../../store/useMapStore";

const INITIAL_VIEW = {
  longitude: -87.7173,
  latitude: 41.8337,
  zoom: 10,
  pitch: 0,
  bearing: 0,
};

interface TooltipInfo {
  x: number;
  y: number;
  zipCode: string;
  totalCrimes: number;
  topCrimeType: string;
}

export default function MapView() {
  const { selectedCity, zipStatsMap, maxCrimes, setSelectedZip } = useMapStore();
  const [tooltip, setTooltip] = useState<TooltipInfo | null>(null);
  const [geojson, setGeojson] = useState<any>(null);
  const [viewState, setViewState] = useState(INITIAL_VIEW);

  // Re-center map when city changes
  useEffect(() => {
    const center = CITY_CENTERS[selectedCity];
    if (center) {
      setViewState((v) => ({ ...v, ...center }));
    }
  }, [selectedCity]);

  // Load the city-specific static GeoJSON file produced by the ETL.
  // File names match the pattern written by cmd/etl: zip_boundaries_{city}.geojson
  // where the city name is lowercased with spaces replaced by underscores.
  useEffect(() => {
    setGeojson(null);
    const filename = selectedCity.toLowerCase().replace(/ /g, "_");
    fetch(`/zip_boundaries_${filename}.geojson`)
      .then((r) => r.json())
      .then(setGeojson)
      .catch((e) => console.error(`Failed to load zip_boundaries_${filename}.geojson:`, e));
  }, [selectedCity]);

  const colorScale = useMemo(
    () => scaleSequential(interpolateYlOrRd).domain([0, maxCrimes]),
    [maxCrimes]
  );

  const toRGBA = (cssColor: string, alpha = 200): [number, number, number, number] => {
    const m = cssColor.match(/\d+/g);
    if (!m) return [100, 100, 100, alpha];
    return [parseInt(m[0]), parseInt(m[1]), parseInt(m[2]), alpha];
  };

  const layer = geojson
    ? new GeoJsonLayer({
        id: "zip-choropleth",
        data: geojson,
        pickable: true,
        stroked: true,
        filled: true,
        getFillColor: (f: any) => {
          const zip = f.properties?.zipCode;
          const stats = zipStatsMap[zip];
          if (!stats || stats.totalCrimes === 0) return [30, 41, 59, 120];
          return toRGBA(colorScale(stats.totalCrimes));
        },
        getLineColor: [255, 255, 255, 60],
        lineWidthMinPixels: 1,
        updateTriggers: {
          getFillColor: [zipStatsMap, maxCrimes],
        },
        onClick: (info: PickingInfo) => {
          const zip = info.object?.properties?.zipCode;
          if (zip) setSelectedZip(zip);
        },
        onHover: (info: PickingInfo) => {
          if (info.object) {
            const zip = info.object.properties?.zipCode ?? "";
            const stats = zipStatsMap[zip];
            setTooltip({
              x: info.x,
              y: info.y,
              zipCode: zip,
              totalCrimes: stats?.totalCrimes ?? 0,
              topCrimeType: stats?.topCrimeType ?? "",
            });
          } else {
            setTooltip(null);
          }
        },
      })
    : null;

  return (
    <div className="absolute inset-0 pt-12">
      <DeckGL
        viewState={viewState}
        onViewStateChange={({ viewState: vs }: any) => setViewState(vs)}
        controller
        layers={layer ? [layer] : []}
        style={{ position: "relative", width: "100%", height: "100%" }}
        getCursor={({ isHovering }) => (isHovering ? "pointer" : "grab")}
      >
        <Map
          mapStyle="https://basemaps.cartocdn.com/gl/dark-matter-gl-style/style.json"
          reuseMaps
        />
      </DeckGL>

      {tooltip && (
        <div
          className="deck-tooltip pointer-events-none fixed z-50"
          style={{ left: tooltip.x + 12, top: tooltip.y + 12 }}
        >
          <div className="font-semibold text-white mb-1">ZIP {tooltip.zipCode}</div>
          <div className="text-slate-300">
            {tooltip.totalCrimes.toLocaleString()} crimes
          </div>
          <div className="text-slate-400 text-xs mt-0.5">
            Top: {tooltip.topCrimeType}
          </div>
        </div>
      )}

      {!geojson && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <p className="text-slate-500 text-sm">Loading ZIP boundaries…</p>
        </div>
      )}
    </div>
  );
}
