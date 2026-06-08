/** @type {import('tailwindcss').Config} */
import forms from "@tailwindcss/forms";
import containerQueries from "@tailwindcss/container-queries";
import type { Config } from "tailwindcss";
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        "surface-container-high": "#e7e8e9",
        "on-primary-container": "#554300",
        "primary-fixed-dim": "#e9c349",
        "primary-fixed": "#ffe088",
        "on-tertiary": "#ffffff",
        "surface-container": "#edeeef",
        "on-background": "#191c1d",
        "on-tertiary-container": "#254188",
        "on-surface-variant": "#4d4635",
        "surface-container-highest": "#e1e3e4",
        "secondary-fixed-dim": "#c1c8ca",
        "tertiary-fixed-dim": "#b4c5ff",
        "tertiary-container": "#97b0ff",
        "on-primary-fixed": "#241a00",
        "on-secondary": "#ffffff",
        "tertiary": "#415ba4",
        "primary-container": "#d4af37",
        "inverse-on-surface": "#f0f1f2",
        "tertiary-fixed": "#dbe1ff",
        "inverse-surface": "#2e3132",
        "secondary-container": "#dae1e3",
        "background": "#f8f9fa",
        "secondary": "#586062",
        "surface-bright": "#f8f9fa",
        "error-container": "#ffdad6",
        "outline-variant": "#d0c5af",
        "on-error-container": "#93000a",
        "on-primary": "#ffffff",
        "error": "#ba1a1a",
        "surface-variant": "#e1e3e4",
        "on-secondary-container": "#5d6466",
        "surface-container-lowest": "#ffffff",
        "on-surface": "#191c1d",
        "on-primary-fixed-variant": "#574500",
        "inverse-primary": "#e9c349",
        "surface-dim": "#d9dadb",
        "secondary-fixed": "#dde4e6",
        "on-tertiary-fixed": "#00174b",
        "on-error": "#ffffff",
        "surface": "#f8f9fa",
        "on-tertiary-fixed-variant": "#27438a",
        "surface-container-low": "#f3f4f5",
        "on-secondary-fixed-variant": "#41484a",
        "primary": "#735c00",
        "surface-tint": "#735c00",
        "on-secondary-fixed": "#161d1f",
        "outline": "#7f7663"
      },
      borderRadius: {
        "DEFAULT": "0.25rem",
        "lg": "0.5rem",
        "xl": "0.75rem",
        "full": "9999px"
      },
      spacing: {
        "margin-desktop": "40px",
        "gutter": "24px",
        "margin-mobile": "16px",
        "base": "8px",
        "container-max": "1440px"
      },
    },
  },
  plugins: [
    forms,
    containerQueries,
  ],
}