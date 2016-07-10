"use strict";

const gulp = require("gulp");
const del = require("del");
const tsc = require("gulp-typescript");
const sourcemaps = require('gulp-sourcemaps');
var concat = require('gulp-concat');

/**
 * Remove build directory.
 */
gulp.task('clean', (cb) => {
    return del(["dist"], cb);
});

/**
 * Compile TypeScript sources and create sourcemaps in build directory.
 */
gulp.task("compile", () => {
    var tsResult = gulp.src([
        'node_modules/angular2/typings/browser.d.ts',
        'src/**/*.ts'
    ]).pipe(sourcemaps.init()) // This means sourcemaps will be generated
        .pipe(tsc({
            sortOutput: true,
            experimentalDecorators:true
        }));

    return tsResult.js
        .pipe(concat('script.js')) // You can use other plugins that also support gulp-sourcemaps
        .pipe(gulp.dest('dist/app'));
});

/**
 * Copy all resources that are not TypeScript files into build directory.
 */
gulp.task("resources", () => {
    return gulp.src(["src/**/*", "!**/*.ts"])
        .pipe(gulp.dest("dist"));
});

/**
 * Build the project.
 */
gulp.task("build", ['compile', 'resources'], () => {
    console.log("Building the project ...");
});