// @ts-check
//
// Shared ESLint flat-config base for the Angular front-ends: /explorer,
// /src/gui/static and /src/skycoin-web.
//
// The three projects had drifted into three different rule sets, so the same
// code could be clean in one and rejected in another. This file holds the rules
// that should hold everywhere; each project's eslint.config.js passes in only
// what is genuinely local to it.
//
// It is a factory rather than a plain config object because ESLint plugins are
// installed per project. Each project requires its own copies of @eslint/js,
// typescript-eslint and angular-eslint and hands them in, which keeps this file
// dependency-free and avoids a repository-root node_modules.
//
// Usage, from a project's eslint.config.js:
//
//   const eslint = require("@eslint/js");
//   const tseslint = require("typescript-eslint");
//   const angular = require("angular-eslint");
//   const createConfig = require("<relative path>/eslint.base.config.js");
//
//   module.exports = createConfig({ eslint, tseslint, angular });

/**
 * @param {object} options
 * @param {any} options.eslint                     `@eslint/js`
 * @param {any} options.tseslint                   `typescript-eslint`
 * @param {any} options.angular                    `angular-eslint`
 * @param {string|string[]} [options.directiveSelectorPrefix]
 *   Attribute-selector prefix(es) accepted for directives. Defaults to "app".
 * @param {Record<string, any>} [options.legacyRules]
 *   Per-project opt-outs. These are pre-existing debt, not policy: each entry
 *   marks a rule the project's code cannot satisfy yet. Deleting entries here
 *   (and fixing the code) is the direction of travel.
 * @param {Record<string, any>} [options.templateRules]
 *   Per-project overrides for the HTML template block.
 */
module.exports = function createAngularEslintConfig({
  eslint,
  tseslint,
  angular,
  directiveSelectorPrefix = "app",
  legacyRules = {},
  templateRules = {},
}) {
  return tseslint.config(
    {
      files: ["**/*.ts"],
      extends: [
        eslint.configs.recommended,
        ...tseslint.configs.recommended,
        ...angular.configs.tsRecommended,
      ],
      processor: angular.processInlineTemplates,
      languageOptions: {
        parserOptions: {
          project: ["tsconfig.json"],
        },
      },
      rules: {
        // --- Angular conventions -------------------------------------------
        "@angular-eslint/component-selector": [
          "error",
          {
            type: "element",
            prefix: "app",
            style: "kebab-case",
          },
        ],
        "@angular-eslint/directive-selector": [
          "error",
          {
            type: "attribute",
            prefix: directiveSelectorPrefix,
            style: "camelCase",
          },
        ],

        // --- Opinionated angular-eslint rules this codebase does not follow --
        // inject() vs constructor injection, and standalone components, are
        // both migrations these NgModule-based apps have not made.
        "@angular-eslint/prefer-inject": "off",
        "@angular-eslint/prefer-standalone": "off",

        // --- Style and correctness held in common ---------------------------
        "max-len": ["error", { code: 200 }],
        "no-underscore-dangle": "off",
        "prefer-const": "error",
        curly: "error",
        eqeqeq: ["error", "always", { null: "ignore" }],
        "valid-typeof": "error",
        "@typescript-eslint/consistent-type-definitions": "error",

        // --- Deliberately relaxed everywhere --------------------------------
        // These fire throughout code that predates the current toolchain and
        // would be noise rather than signal.
        "no-constant-binary-expression": "off",
        "@typescript-eslint/no-explicit-any": "off",
        "@typescript-eslint/no-unused-vars": "off",
        "@typescript-eslint/no-empty-function": "off",

        // --- Per-project debt ------------------------------------------------
        ...legacyRules,
      },
    },
    {
      files: ["**/*.html"],
      extends: [...angular.configs.templateRecommended],
      rules: {
        ...templateRules,
      },
    }
  );
};
