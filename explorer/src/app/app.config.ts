/* eslint-disable @typescript-eslint/naming-convention */

/**
 * General configuration file.
 */


/**
 * Names to identify the coins and hours in case it is not possible to obtain the values
 * from the node.
 */
export const CoinIdentifiers = {
  fullName: 'Skycoin',
  coinName: 'SKY',
  HoursName: 'Coin Hours',
  HoursNameSingular: 'Coin Hour',
};

/**
 * Configuration for the QR codes.
 */
export const QrConfig = {
  /**
   * Prefix that will be added to the addresses in the QR codes, to identify what coin the address
   * is for. Corresponds to the BIP-21 specification.
   */
  prefix: 'skycoin:',
};

/**
 * Configuration for the generic header. Read the customization docs for more info.
 */
export const HeaderConfig = {
  // Set to true for using the generic header, instead of the Skycoin one.
  useGenericHeader: false,
  genericHeaderUrl: 'https://www.skycoin.com/',
};

/**
 * Configuration for the generic footer. Read the customization docs for more info.
 */
export const FooterConfig = {
  // Set to true for using the generic footer, instead of the Skycoin one.
  useGenericFooter: false,
  contactLinks: [
    {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-github"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-telegram"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-twitter"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-youtube"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-discord"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-instagram"></i>',
    } , {
      url: 'https://www.skycoin.com/',
      content: '<i class="fab fa-reddit"></i>',
    }
  ],
};

/**
 * Laguage configuration.
 */
export const languageConfig = {
  /**
   * List of available languages. See the documentation in the assets/i18n folder and the
   * LanguageData class (inside the LanguageService file) for more information.
   */
  languages: [{
      code: 'en',
      name: 'English',
      iconName: 'en.png'
    },
    {
      code: 'es',
      name: 'Español',
      iconName: 'es.png'
    },
    {
      code: 'zh',
      name: '中文',
      iconName: 'zh.png'
    }
  ],
  /**
   * Code of the default language.
   */
  defaultLanguage: 'en'
};

/**
 * List with named addresses.
 */
export const namedAddresses = [
  {
    address: '2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8',
    name: 'Binance hot wallet'
  }
];

/**
 * If the user navigates back or forward with the browser buttons, the app will restore the
 * previous data. After that, if the data was saved too long ago (which means that it is older
 * than the time defined on this var, in ms), the page will request the data from the server
 * again on the background, to refresh the info displayed to the user.
 */
export const dataValidityTime = 60000;
