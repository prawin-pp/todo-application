/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,js,svelte,ts}', './node_modules/flowbite/**/*.js'],
  theme: {
    extend: {
      colors: {
        primary: '#ff0090',
        'primary-700': '#d90082',
        'primary-900': '#8c005e',
      },
    },
  },
  plugins: [require('flowbite/plugin')],
};
