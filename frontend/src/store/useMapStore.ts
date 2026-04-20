import { create } from "zustand";

export interface ZipStats {
  zipCode: string;
  totalCrimes: number;
  topCrimeType: string;
}

// City centre coordinates for map re-centering
export const CITY_CENTERS: Record<string, { longitude: number; latitude: number; zoom: number }> = {
  Chicago: { longitude: -87.7173, latitude: 41.8337, zoom: 10 },
  "Los Angeles": { longitude: -118.2437, latitude: 34.0522, zoom: 10 },
  Riverside: { longitude: -117.3961, latitude: 33.9533, zoom: 11 },
};

interface MapStore {
  selectedCity: string;
  selectedYear: number;
  selectedZip: string | null;
  zipStatsMap: Record<string, ZipStats>;
  maxCrimes: number;

  setSelectedCity: (city: string) => void;
  setSelectedYear: (year: number) => void;
  setSelectedZip: (zip: string | null) => void;
  setZipStats: (stats: ZipStats[]) => void;
}

export const useMapStore = create<MapStore>((set) => ({
  selectedCity: "Chicago",
  selectedYear: 0,
  selectedZip: null,
  zipStatsMap: {},
  maxCrimes: 1,

  setSelectedCity: (city) => set({ selectedCity: city, selectedYear: 0, selectedZip: null, zipStatsMap: {}, maxCrimes: 1 }),
  setSelectedYear: (year) => set({ selectedYear: year, selectedZip: null }),
  setSelectedZip: (zip) => set({ selectedZip: zip }),

  setZipStats: (stats) => {
    const map: Record<string, ZipStats> = {};
    let max = 1;
    for (const s of stats) {
      map[s.zipCode] = s;
      if (s.totalCrimes > max) max = s.totalCrimes;
    }
    set({ zipStatsMap: map, maxCrimes: max });
  },
}));
