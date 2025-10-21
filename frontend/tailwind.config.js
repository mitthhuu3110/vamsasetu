/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./src/**/*.{js,jsx,ts,tsx}",
  ],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        // Indian Theme Colors
        'warm-beige': '#F5F1EB',
        'soft-gold': '#D4AF37',
        'deep-gold': '#B8860B',
        'soft-green': '#90EE90',
        'forest-green': '#228B22',
        'warm-brown': '#8B4513',
        'cream': '#FFF8DC',
        'saffron': '#FF9933',
        'maroon': '#800000',
        'navy': '#000080',
        // Dark mode colors
        'dark-bg': '#1a1a1a',
        'dark-card': '#2d2d2d',
        'dark-text': '#e5e5e5',
        'dark-accent': '#4a5568',
        // Primary colors for light mode
        primary: {
          50: '#F5F1EB',
          100: '#E8DCC0',
          200: '#D4AF37',
          300: '#B8860B',
          400: '#9A7209',
          500: '#7C5A07',
          600: '#5E4305',
          700: '#402C03',
          800: '#221501',
          900: '#0A0A00',
        },
        secondary: {
          50: '#F0FFF0',
          100: '#DCFFDC',
          200: '#90EE90',
          300: '#228B22',
          400: '#1E7B1E',
          500: '#1A6B1A',
          600: '#165B16',
          700: '#124B12',
          800: '#0E3B0E',
          900: '#0A2B0A',
        }
      },
      fontFamily: {
        'sans': ['Inter', 'system-ui', 'sans-serif'],
        'display': ['Playfair Display', 'serif'],
        'telugu': ['Noto Sans Telugu', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'glow': 'glow 2s ease-in-out infinite alternate',
        'pulse-slow': 'pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        glow: {
          '0%': { boxShadow: '0 0 5px #D4AF37' },
          '100%': { boxShadow: '0 0 20px #D4AF37, 0 0 30px #D4AF37' },
        },
      },
      backgroundImage: {
        'gradient-radial': 'radial-gradient(var(--tw-gradient-stops))',
        'family-pattern': "url('data:image/svg+xml,%3Csvg width=\"60\" height=\"60\" viewBox=\"0 0 60 60\" xmlns=\"http://www.w3.org/2000/svg\"%3E%3Cg fill=\"none\" fill-rule=\"evenodd\"%3E%3Cg fill=\"%23D4AF37\" fill-opacity=\"0.1\"%3E%3Ccircle cx=\"30\" cy=\"30\" r=\"2\"/%3E%3C/g%3E%3C/g%3E%3C/svg%3E')",
      },
    },
  },
  plugins: [],
}
