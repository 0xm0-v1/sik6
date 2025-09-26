/**
 * ESLint Reference Configuration for sik6 (Angular + TypeScript)
 * --------------------------------------------------------------
 * This file is for documentation and team onboarding.
 * The actual running config is `.eslintrc.json`.
 *
 * Key ideas:
 * - Keep the executable config minimal at start.
 * - Prefer adding strict rules progressively to avoid noise.
 */

module.exports = {
  // Root marks this config as the top-most; ESLint won't merge with parent folders.
  root: true,

  // Environments define global variables and parsing expectations.
  env: {
    browser: true, // Window, document, etc.
    es2022: true
  },

  // Parser options: keep it simple first. If you adopt type-aware rules later,
  // you'll set `parserOptions.project` to your tsconfig path.
  parserOptions: {
    ecmaVersion: 2022,
    sourceType: "module",
    project: false // set a tsconfig path when enabling type-aware rules
  },

  // Base recommendations + Angular ESLint plugin recommendations.
  extends: [
    "eslint:recommended",
    "plugin:@angular-eslint/recommended"
    // You may add "plugin:@typescript-eslint/recommended" when you enable the TS plugin explicitly.
  ],

  // Per-file-type overrides: TS source files and Angular HTML templates.
  overrides: [
    {
      files: ["*.ts"],
      // TypeScript parser enables TS syntax; type-aware mode requires parserOptions.project.
      parser: "@typescript-eslint/parser",

      // Angular-specific linting; you can add @typescript-eslint plugin later if needed.
      plugins: ["@angular-eslint"],

      rules: {
        // Start minimal: enforce unused variables as errors (keeps code clean early).
        "no-unused-vars": "error",

        // Allow console in dev, warn to discourage accidental commits.
        "no-console": "warn"

        // Examples to enable later:
        // "@angular-eslint/directive-selector": ["error", { "type": "attribute", "prefix": "app", "style": "camelCase" }],
        // "@angular-eslint/component-selector": ["error", { "type": "element", "prefix": "app", "style": "kebab-case" }]
      }
    },
    {
      files: ["*.html"],
      // Lint Angular templates (accessibility, best practices).
      plugins: ["@angular-eslint/template"],
      extends: ["plugin:@angular-eslint/template/recommended"],
      rules: {
        // Examples to enable later:
        // "@angular-eslint/template/click-events-have-key-events": "warn",
        // "@angular-eslint/template/interactive-supports-focus": "warn"
      }
    }
  ]

  // Notes:
  // - Prettier/formatting lives in your formatter ticket (kept separate from linting).
  // - When your project grows, consider:
  //   * Adding @typescript-eslint/eslint-plugin rules
  //   * Enabling type-aware linting (set parserOptions.project)
  //   * Adding import ordering rules, accessibility rules, etc.
};
