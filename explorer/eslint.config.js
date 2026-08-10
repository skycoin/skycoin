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
      // default strategy to Eager. These components declare Eager explicitly so
      // they keep the behaviour they had before Angular 22, which this rule
      // reports as opting out of the default. Adopting OnPush is worth doing,
      // but it is a behavioural change that has to be reviewed and tested
      // component by component.
      "@angular-eslint/prefer-on-push-component-change-detection": "off",
      "no-constant-binary-expression": "off",
      "@typescript-eslint/consistent-type-definitions": "error",
      "@typescript-eslint/dot-notation": "off",
      "@typescript-eslint/explicit-member-accessibility": [
        "off",
        {
          accessibility: "explicit",
        },
      ],
      "max-len": [
        "error",
        {
          code: 200,
        },
      ],
      "no-underscore-dangle": "off",
      "valid-typeof": "error",
      "@typescript-eslint/member-ordering": "off",
      "object-shorthand": "off",
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-unused-vars": "off",
      "@typescript-eslint/no-empty-function": "off",
    },
  },
  {
    files: ["**/*.html"],
    extends: [
      ...angular.configs.templateRecommended,
    ],
    rules: {},
  }
);
