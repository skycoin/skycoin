/**
 * This file includes polyfills needed by Angular and is loaded before the app.
 * You can add your own extra polyfills to this file.
 *
 * This file is divided into 2 sections:
 *   1. Browser polyfills. These are applied before loading ZoneJS and are sorted by browsers.
 *   2. Application imports. Files imported after ZoneJS that should be loaded before your main
 *      file.
 *
 * The current setup is for so-called "evergreen" browsers; the last versions of browsers that
 * automatically update themselves. This includes recent versions of Safari, Chrome, Firefox,
 * and Edge.
 */

/***************************************************************************************************
 * BROWSER POLYFILLS
 */

// Modern browsers don't need most polyfills

/***************************************************************************************************
 * Zone JS is required by Angular itself.
 */
import 'zone.js';

/***************************************************************************************************
 * APPLICATION IMPORTS
 */

// Global shims for Node.js modules used by crypto libraries
(window as any).global = window;

// Buffer polyfill for libraries that expect the Node global.
//
// bip39 calls Buffer.from() inside mnemonicToSeed(), so generating a new wallet
// seed throws "Buffer is not defined" without this. The webpack build supplied
// it through ProvidePlugin; esbuild has no equivalent, so it is declared here.
import { Buffer } from 'buffer';
(window as any).Buffer = Buffer;
