/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Inter", "system-ui", "sans-serif"],
      },
      colors: {
        chicago: {
          blue: "#003087",
          red: "#CC0000",
        },
      },
    },
  },
  plugins: [],
};
