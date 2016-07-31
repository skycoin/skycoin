"use strict";

const gulp = require("gulp");
const del = require("del");
const tsc = require("gulp-typescript");
const sourcemaps = require('gulp-sourcemaps');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
const tsProject = tsc.createProject("tsconfig.json");
var replace = require('gulp-replace');

/**
 * Remove dist directory.
 */
gulp.task('clean', (cb) => {
    return del(["dist"], cb);
});

gulp.task('clean_dev', (cb) => {
    return del(["dev"], cb);
});

/**
 * Compile TypeScript sources and create sourcemaps in build directory.
 */
gulp.task("compile", () => {
    let tsResult = gulp.src("src/**/*.ts")
        .pipe(sourcemaps.init())
        .pipe(tsc(tsProject));
return tsResult.js
    .pipe(replace('.ts', '.js'))
    .pipe(sourcemaps.write("."))
    .pipe(gulp.dest("dist"));
});

gulp.task("compile_dev", () => {
    let tsResult = gulp.src("src/**/*.ts")
        .pipe(sourcemaps.init())
        .pipe(tsc(tsProject));
return tsResult.js
    .pipe(replace('.ts', '.js'))
    .pipe(sourcemaps.write("."))
    .pipe(gulp.dest("dev"));
});

/**
 * Copy all resources that are not TypeScript files into build directory.
 */
gulp.task("resources", () => {
    return gulp.src(["src/**/*", "!**/*.ts"])
        .pipe(gulp.dest("dist"));
});

gulp.task("resources_dev", () => {
    return gulp.src(["src/**/*", "!**/*.ts"])
        .pipe(gulp.dest("dev"));
});

/**
 * Build the project.
 */
gulp.task("build", ['compile_dev', 'resources_dev'], () => {
    console.log("Building the project to dev directory...");
});

/**
 * Build the project.
 */
gulp.task("dist", ['compile', 'resources'], () => {
    console.log("Building the project to dist directory...");
});
