'use strict';

var gulp = require('gulp');
var electron = require('gulp-electron');
var exec = require('child_process').exec;
var packageJson = require('./src/package.json');

gulp.task('electron', () => {
    var platforms = [
        'win32-x64',
        'win32-ia32',
        'darwin-x64',
        'linux-x64',
        'linux-arm'
    ];
    if (process.argv.length > 4) {
        platforms = [process.argv[4]];
    }
    console.log(platforms);
    gulp.src("")
        .pipe(electron({
            src: './src',
            packageJson: packageJson,
            release: './.electron_output',
            cache: './.electron_cache',
            version: 'v1.4.13', // electron version
            packaging: false, // zip/tar results; we do this manually since
            // we need to copy our skycoin binaries in
            // due to liimitations of electron-gulp
            // token: 'abc123...',  // GITHUB_TOKEN if there is ratelimit issue
            platforms: platforms,
            platformResources: {
                darwin: {
                    CFBundleDisplayName: packageJson.productName,
                    CFBundleIdentifier: 'org.skycoin.Skycoin',
                    CFBundleName: packageJson.productName,
                    CFBundleVersion: packageJson.version,
                    CFBundleURLTypes: [{
                        CFBundleURLName: 'Skycoin',
                        CFBundleURLSchemes: ['skycoin'],
                    }],
                    icon: './assets/osx/appIcon.icns'
                },
                win: {
                    "version-string": packageJson.version,
                    "file-version": packageJson.version,
                    "product-version": packageJson.version,
                    "icon": './assets/windows/favicon.ico'
                }
            }
        }))
        .pipe(gulp.dest(""));
});