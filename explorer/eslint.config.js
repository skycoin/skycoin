// @ts-check
const eslint = require("@eslint/js");
const tseslint = require("typescript-eslint");
const angular = require("angular-eslint");
const createAngularEslintConfig = require("../eslint.base.config.js");

module.exports = createAngularEslintConfig({
  eslint,
  tseslint,
  angular,
  // Pre-existing debt, not policy. Each entry is a rule this project's code
  // cannot satisfy yet; fixing the code and deleting the entry is the goal.
  legacyRules: {
    "@typescript-eslint/dot-notation": "off",
    "@typescript-eslint/member-ordering": "off",
    "@typescript-eslint/explicit-member-accessibility": [
      "off",
      {
        accessibility: "explicit",
      },
    ],
    "object-shorthand": "off",
  },
});
