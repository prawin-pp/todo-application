/**
 * @type {import('prettier').Config}
 **/
module.exports = {
  printWidth: 100,
  tabWidth: 2,
  semi: true,
  singleQuote: true,
  trailingComma: 'es5',
  bracketSpacing: true,
  bracketSameLine: false,
  arrowParens: 'always',
  endOfLine: 'lf',
  plugins: ['prettier-plugin-svelte', 'prettier-plugin-tailwindcss'],
  pluginSearchDirs: ['.'],
  overrides: [{ files: '*.svelte', options: { parser: 'svelte' } }],
};
