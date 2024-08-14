import daisyui from "daisyui"
import tailwindTypography from "@tailwindcss/typography"
import tailwindContainerQueries from "@tailwindcss/container-queries"

/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {},
  },
  plugins: [
    tailwindContainerQueries,
    tailwindTypography,
    daisyui,
  ],
}

