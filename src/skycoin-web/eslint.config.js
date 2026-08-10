// @ts-check
const eslint = require("@eslint/js");
const tseslint = require("typescript-eslint");
const angular = require("angular-eslint");
const createAngularEslintConfig = require("../../eslint.base.config.js");

module.exports = createAngularEslintConfig({
  eslint,
  tseslint,
  angular,
  // ClipboardDirective is declared with a `clipboard` attribute selector.
  directiveSelectorPrefix: ["app", "clipboard"],
  // Pre-existing debt, not policy. Each entry is a rule this project's code
  // cannot satisfy yet; fixing the code and deleting the entry is the goal.
  legacyRules: {
    // The only front-end still on ChangeDetectionStrategy.Eager. Its components
    // declare Eager explicitly to keep their pre-Angular-22 behaviour, which
    // this rule reports as opting out of the default. Remove this entry once
    // skycoin-web is migrated to OnPush, as explorer and gui/static have been.
    "@angular-eslint/prefer-on-push-component-change-detection": "off",
    "@angular-eslint/no-output-on-prefix": "off",
    "@typescript-eslint/no-wrapper-object-types": "off",
    "@typescript-eslint/no-unused-expressions": "off",
    "@typescript-eslint/no-empty-object-type": "off",
    "@typescript-eslint/no-require-imports": "off",
    "no-useless-assignment": "off",
    "no-empty": "off",
    "no-var": "off",
    "no-self-assign": "off",
  },
  templateRules: {
    "@angular-eslint/template/eqeqeq": "off",
  },
});
