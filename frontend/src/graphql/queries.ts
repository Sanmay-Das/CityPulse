import { gql } from "@apollo/client";

export const GET_ALL_ZIP_STATS = gql`
  query GetAllZipStats($city: String, $year: Int) {
    allZipStats(city: $city, year: $year) {
      zipCode
      totalCrimes
      topCrimeType
    }
  }
`;

export const GET_AVAILABLE_CITIES = gql`
  query GetAvailableCities {
    availableCities
  }
`;

export const GET_CRIME_TIMELINE = gql`
  query GetCrimeTimeline($city: String, $zipCode: String!) {
    crimeTimeline(city: $city, zipCode: $zipCode) {
      year
      month
      count
    }
  }
`;

export const GET_CRIME_TYPES = gql`
  query GetCrimeTypes($city: String, $zipCode: String!, $year: Int) {
    crimeTypes(city: $city, zipCode: $zipCode, year: $year) {
      crimeType
      count
    }
  }
`;

export const GET_AVAILABLE_YEARS = gql`
  query GetAvailableYears($city: String) {
    availableYears(city: $city)
  }
`;

export const GET_YEAR_RANGE = gql`
  query GetYearRange($city: String) {
    yearRange(city: $city) {
      minYear
      maxYear
    }
  }
`;
