export const AppConfig = {

  // General settings.
  ////////////////////////////////

  /**
   * If true, the option for buying coins via the OTC service will be enabled.
   * Note: not updated in long time, most likely it won't work.
   */
  otcEnabled: false,
  /**
   * How many addresses the hw wallet can have.
   */
  maxHardwareWalletAddresses: 1,
  /**
   * Max gap of unused addresses a wallet can have between 2 used addresses before the user is
   * alerted about potential problems that could appear for restoring all addresses when loading
   * the wallet again using the seed.
   */
  maxAddressesGap: 20,
  /**
   * ID of the coin on the coin price service. If null, the wallet will not show the USD price.
   */
  priceApiId: 'sky-skycoin',
  /**
   * This wallet uses the Skycoin URI Specification (based on BIP-21) when creating QR codes and
   * requesting coins. This variable defines the prefix that will be used for creating QR codes
   * and URLs. IT MUST BE UNIQUE FOR EACH COIN.
   */
  uriSpecificatioPrefix: 'skycoin',
  /**
   * Normal size for some modal windows.
   */
  mediumModalWidth: '566px',

  // Hw wallet firmware.
  ////////////////////////////////

  /**
   * URL for checking the number of the most recent version of the Skywallet firmware.
   */
  urlForHwWalletVersionChecking: 'https://version.skycoin.com/skywallet/version.txt',
  /**
   * First part of the URL for donwnloading the lastest firmware for the Skywallet. The number of
   * the lastest version and '.bin' is added at the end of the value by the code.
   */
  hwWalletDownloadUrlAndPrefix: 'https://downloads.skycoin.com/skywallet/skywallet-firmware-v',
  /**
   * URL were the user can download the lastest version of the hw wallet daemon.
   */
  hwWalletDaemonDownloadUrl: 'https://www.skycoin.com/downloads/',

  // Wallet update.
  ////////////////////////////////

  /**
   * URL for checking the number of the most recent version of the wallet software.
   */
  urlForVersionChecking: 'https://version.skycoin.com/skycoin/version.txt',
  /**
   * URL were the user can download the lastest version of the wallet software.
   */
  walletDownloadUrl: 'https://www.skycoin.com/downloads/',

  // Translations.
  ////////////////////////////////

  /**
   * Array with the available translations. For more info check the readme file in the
   * folder with the translation files.
   */
  languages: [{
      code: 'en',
      name: 'English',
      iconName: 'en.png',
    },
    {
      code: 'zh',
      name: '中文',
      iconName: 'zh.png',
    },
    {
      code: 'es',
      name: 'Español',
      iconName: 'es.png',
    },
  ],

  /**
   * Default language used by the software.
   */
  defaultLanguage: 'en',
};
