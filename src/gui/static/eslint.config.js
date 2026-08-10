// @ts-check
const eslint = require("@eslint/js");
const tseslint = require("typescript-eslint");
const angular = require("angular-eslint");

module.exports = tseslint.config(
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
          prefix: "app",
          style: "camelCase",
        },
      ],
      "@angular-eslint/prefer-inject": "off",
      "@angular-eslint/prefer-standalone": "off",
      // Angular 22 made OnPush the framework default and renamed the previous
      // default strategy to Eager. The v22 upgrade migration added an explicit
      // `ChangeDetectionStrategy.Eager` to every component here specifically to
      // preserve the pre-22 behaviour, which this rule then reports as opting
      // out of the default. Adopting OnPush is worth doing, but it is a
      // behavioural change that has to be reviewed and tested component by
      // component rather than carried along by a dependency bump.
      "@angular-eslint/prefer-on-push-component-change-detection": "off",
      "no-constant-binary-expression": "off",
      "no-useless-assignment": "off",
      "no-empty": "off",
      "no-var": "off",
      "no-self-assign": "off",
      "@angular-eslint/no-output-on-prefix": "off",
      "@typescript-eslint/no-wrapper-object-types": "off",
      "@typescript-eslint/no-unused-expressions": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-unused-vars": "off",
      "@typescript-eslint/no-empty-function": "off",
      "@typescript-eslint/no-empty-object-type": "off",
      "@typescript-eslint/no-require-imports": "off",
      "max-len": ["error", { code: 200 }],
      "no-underscore-dangle": "off",
      "prefer-const": "error",
      curly: "error",
      eqeqeq: ["error", "always", { null: "ignore" }],
    },
  },
  {
    files: ["**/*.html"],
    extends: [
      ...angular.configs.templateRecommended,
    ],
    rules: {
      "@angular-eslint/template/eqeqeq": "off",
    },
  }
);
