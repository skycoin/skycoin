'use strict'

const fs = require('fs');

/////////////////////////////////////////////
// Initial configuration
/////////////////////////////////////////////

console.log('Starting to check the language files.', '\n');

// Load the current English file.
if (!fs.existsSync('en.json')) {
  exitWithError('Unable to find the English language file.');
}
let currentData = JSON.parse(fs.readFileSync('en.json', 'utf8'));

// 2 charaters code of the languages that will be checked.
const langs = [];
// If false, the code will only verify the differences in the elements (ignoring its contents) of the
// base files and the files with the translations. If not, the code will also verify the differences
// in the elements and contents of the base files and the current English file.
let checkFull = false;

// If a param was send, it must be "all" or the 2 charaters code of the language that must be checked.
// If a param is provided, checkFull is set to true.
if (process.argv.length > 2) {
  if (process.argv.length > 3) {
    exitWithError('Invalid number of parameters.');
  }

  if (process.argv[2] != 'all') {
    if (process.argv[2].length !== 2) {
      exitWithError('You can only send as parameter to this script the 2-letter code of one of the language files in this folder, or "all".');
    }
    langs.push(process.argv[2]);
  }

  checkFull = true;
}

// If no language code was send as param, the code will check all languages.
if (langs.length === 0) {
  let localFiles = fs.readdirSync('./');

  const langFiles = [];
  const langFilesMap = new Map();
  const baseLangFilesMap = new Map();
  localFiles.forEach(file => {
    if (file.length === 12 && file.endsWith('_base.json')) {
      langs.push(file.substring(0, 2));
      baseLangFilesMap.set(file.substring(0, 2), true);
    }
    if (file !== 'en.json' && file.length === 7 && file.endsWith('.json')) {
      langFiles.push(file.substring(0, 2));
      langFilesMap.set(file.substring(0, 2), true);
    }
  });

  langs.forEach(lang => {
    if (!langFilesMap.has(lang)) {
      exitWithError('The \"' + lang + '_base.json\" base file does not have its corresponding language file.');
    }
  });

  langFiles.forEach(lang => {
    if (!baseLangFilesMap.has(lang)) {
      exitWithError('The \"' + lang + '.json\" file does not have its corresponding base file.');
    }
  });

  if (langs.length === 0) {
    exitWithError('No language files to check.');
  }
}

console.log('Checking the following languages:');
langs.forEach(lang => {
  console.log(lang);
});
console.log('');

/////////////////////////////////////////////
// Verifications
/////////////////////////////////////////////

// The following arrays will contain the list of elements with problems. Each element of the
// arrays contains a "lang" property with the language identifier and a "elements" array with
// the path of the problematic elements.

// Elements that are present in a base file but not in its corresponding translation file.
const baseFileOnly = [];
// Elements that are present in a translation file but not in its corresponding base file.
const translatedFileOnly = [];
// Elements that are present in the English file but not in the currently checked base translation file.
const enOnly = [];
// Elements that are present in the currently checked base translation file but not in the English file.
const translatedOnly = [];
// Elements that have different values in the currently checked base translation file and the English file.
const different = [];

function addNewLangToArray(array, lang) {
  array.push({
    lang: lang,
    elements: []
  });
}

langs.forEach(lang => {
  addNewLangToArray(baseFileOnly, lang);
  addNewLangToArray(translatedFileOnly, lang);
  addNewLangToArray(enOnly, lang);
  addNewLangToArray(translatedOnly, lang);
  addNewLangToArray(different, lang);

  // Try to load the translation file and its corresponding base file.
  if (!fs.existsSync(lang + '.json')) {
    exitWithError('Unable to find the ' + lang +  '.json file.');
  }
  let translationData = JSON.parse(fs.readFileSync(lang + '.json', 'utf8'));

  if (!fs.existsSync(lang + '_base.json')) {
    exitWithError('Unable to find the ' + lang +  '_base.json language file.');
  }
  let baseTranslationData = JSON.parse(fs.readFileSync(lang + '_base.json', 'utf8'));

  // Check the differences in the elements of the translation file and its base file.
  checkElement('', '', baseTranslationData, translationData, baseFileOnly, true);
  checkElement('', '', translationData, baseTranslationData, translatedFileOnly, true);

  // Check the differences in the elements and content of the base translation file the English file.
  if (checkFull) {
    checkElement('', '', currentData, baseTranslationData, enOnly, false);
    checkElement('', '', baseTranslationData, currentData, translatedOnly, true);
  }
});


// Check recursively if the elements and content of two language objects are the same.
//
// path: path of the currently checked element. As this function works with nested elements,
// the path is the name of all the parents, separated by a dot.
// key: name of the current element.
// fist: first element for the comparation.
// second: second element for the comparation.
// arrayForMissingElements: array in which the list of "fist" elements that are not in "second"
// will be added.
// ignoreDifferences: if false, each time the content of an element in "fist" is different to the
// same element in "second" that element will be added to the "different" array.
function checkElement(path, key, fist, second, arrayForMissingElements, ignoreDifferences) {
  let pathPrefix = '';
  if (path.length > 0) {
    pathPrefix = '.';
  }

  // This means that, at some point, the code found an element in the "first" branch that is
  // not in the "second" branch.
  if (second === undefined || second === null) {
    arrayForMissingElements[arrayForMissingElements.length - 1].elements.push(path + pathPrefix + key);
    return;
  }

  if (typeof fist !== 'object') {
    // If the current element is a string, compare the contents, but ony if ignoreDifferences
    // is true.
    if (!ignoreDifferences && fist != second) {
      different[different.length - 1].elements.push(path + pathPrefix + key);
    }
  } else {
    // If the current element is an object, check the childs.
    Object.keys(fist).forEach(currentKey => {
      checkElement(path + pathPrefix + key, currentKey, fist[currentKey], second[currentKey], arrayForMissingElements, ignoreDifferences);
    });
  }
}

/////////////////////////////////////////////
// Results processing
/////////////////////////////////////////////

// Becomes true if any of the verifications failed.
let failedValidation = false;

// If "failedValidation" is false, writes to the console the header of the error list
// and updates the value of "failedValidation" to true.
function updateErrorSumary() {
  if (!failedValidation) {
    failedValidation = true;
    console.log('The following problems were found:', '\n');
  }
}

// Checks all arrays for errors. This loop is for the languages.
for (let i = 0; i < baseFileOnly.length; i++) {

  // This loop if for checking all the arrays.
  [baseFileOnly, translatedFileOnly, enOnly, translatedOnly, different].forEach((array, idx) => {
    // If the array has elements, it means that errors were found.
    if (array[i].elements.length > 0) {
      updateErrorSumary();

      // Show the appropriate error text according to the current array.
      if (idx === 0) {
        console.log('The \"' + baseFileOnly[i].lang + '_base.json\" base file has elements that are not present in \"' + baseFileOnly[i].lang + '.json\":');
      } else if (idx === 1) {
        console.log("\"" + translatedFileOnly[i].lang + '.json\" has elements that are not present in the \"' + baseFileOnly[i].lang + '_base.json\" base file:');
      } else if (idx === 2) {
        console.log('The \"en.json\" file has elements that are not present in the \"' + enOnly[i].lang + '_base.json\" base file:');
      } else if (idx === 3) {
        console.log('The \"' + translatedOnly[i].lang + '_base.json\" base file has elements that are not present in \"en.json\":');
      } else if (idx === 4) {
        console.log('The \"' + different[i].lang + '_base.json\" base file has values that do not match the ones in \"en.json\":');
      }
      // Show all the elements with errors.
      array[i].elements.forEach(element => console.log(element));
      console.log('');
    }
  });
}

// If no error was detected, show a success message on the console. If not, exit with an error code.
if (!failedValidation) {
  console.log('The verification passed without problems.');
} else {
  process.exit(1);
}

function exitWithError(errorMessage) {
  console.log('Error: ' + errorMessage)
  process.exit(1);
}
