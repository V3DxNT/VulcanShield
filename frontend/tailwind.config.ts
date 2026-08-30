import type { Config } from 'tailwindcss';

const config: Config = {
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', '-apple-system', 'BlinkMacSystemFont', 'Helvetica Neue', 'Helvetica', 'sans-serif'],
      },
      colors: {
        apple: {
          blue: '#0071e3',
          gray: '#f5f5f7',
          dark: '#1d1d1f',
          mid: '#6e6e73',
          light: '#a1a1a6',
          border: '#d2d2d7',
          green: '#30d158',
          amber: '#ff9f0a',
          red: '#ff3b30',
        }
      },
    },
  },
  plugins: [],
};
export default config;
