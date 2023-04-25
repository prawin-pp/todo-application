/**
 * @type {import('prettier').Config}
 **/
module.exports = {
  printWidth: 100,
  tabWidth: 2,
  semi: true,
  singleQuote: true,
  trailingComma: "es5",
  bracketSpacing: true,
  bracketSameLine: false,
  arrowParens: "always",
  endOfLine: "lf",
  plugins: [require("prettier-plugin-svelte"), require("prettier-plugin-tailwindcss")],
  tailwindConfig: "./tailwind.config.cjs",
  svelteSortOrder: "options-scripts-markup-styles",
  svelteStrictMode: false,
  svelteAllowShorthand: false,
  svelteBracketNewLine: true,
  svelteIndentScriptAndStyle: true,
  overrides: [
    {
      files: "*.svelte",
      options: { parser: "svelte" },
    },
  ],
};
