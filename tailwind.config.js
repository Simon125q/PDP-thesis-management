/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.html", "./**/*.templ", "./**/*.go"],
  safelist: [],
  theme: {
    extend: {
      colors: {
        "tul-dark-red": "#8b0002",
        "tul-medium-red": "#832941",
        "tul-light-red": "#c47979",
        "tul-dark-gray": "#1d1d1d",
      },
    },
  },
};
