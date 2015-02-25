'use strict';
// =======================================================================
// Gulp Plugins
// =======================================================================
var gulp          = require('gulp'),
  connect         = require('gulp-connect'),
  gutil           = require('gulp-util'),
  jshint          = require('gulp-jshint'),
  stylish         = require('jshint-stylish'),
  concat          = require('gulp-concat'),
  streamify       = require('gulp-streamify'),
  uglify          = require('gulp-uglify'),
  sourcemaps      = require('gulp-sourcemaps'),
  stylus          = require('gulp-stylus'),
  nib             = require('nib'),
  prefix          = require('gulp-autoprefixer'),
  minifyCSS       = require('gulp-minify-css'),
  notify          = require('gulp-notify'),
  watchify        = require('watchify'),
  del             = require('del'),
  source          = require('vinyl-source-stream'),
  buffer          = require('vinyl-buffer'),
  runSequence     = require('run-sequence');


// =======================================================================
// File Paths
// =======================================================================
var filePath = {
  build: {
    dest: './dist'
  },
  lint: {
    src: ['./app/*.js', './app/**/*.js']
  },
  browserify: {
    src: './app/app.js',
    watch:
      [
      '!./app/assets/libs/*.js',
      '!./app/assets/libs/**/*.js',
      './app/*.js','./app/**/*.js',
      '/app/**/*.html'
    ]
  },
  styles: {
    src: './app/app.styl',
    watch: ['./app/app.styl','./app/**/*.styl']
  },
  images: {
    src: './app/assets/images/**/*',
    watch: ['./app/assets/images', './app/assets/images/**/*'],
    dest: './dist/images/'
  },
  vendorJS: {
    // These files will be bundled into a single vendor.js file that's called at the bottom of index.html
    src:
      [
      './libs/jquery/dist/jquery.js', // v2.1.1
      './libs/bootstrap/dist/js/bootstrap.js' // v3.1.1
    ]
  },
  vendorFONTS: {
    src: [
      './libs/bootstrap/fonts/*', // v3.1.1
      './libs/font-awesome/fonts/*'
    ],
    dest: './dist/fonts/'
  },
  vendorIMG: {
    src: [
      './libs/bootstrap-glyphicons/*.png'
    ],
    dest: './dist/img/'
  },
  vendorCSS: {
    src:
      [
      './libs/bootstrap/dist/css/bootstrap.css', // v3.1.1
      './libs/font-awesome/css/font-awesome.css', // v4.1.0
      './libs/bootstrap-glyphicons/css/bootstrap.icon-large.css'
    ]
  },
  copyIndex: {
    src: './app/index.html',
    watch: './app/index.html'
  },
  copyFavicon: {
    src: './app/favicon.png'
  }
};


// =======================================================================
// Error Handling
// =======================================================================
function handleError(err) {
  gutil.log(err.toString());
  /*jshint validthis:true */
  this.emit('end');
}


// =======================================================================
// Server Task
// =======================================================================
var express = require('express'),
  server  = express();

// Server settings
server.use(express.static(filePath.build.dest));
// Redirects everything back to our index.html
server.all('/*', function(req, res) {
  res.sendfile('/', { root: filePath.build.dest });
});

// uncomment the "middleware" section when you are ready to connect to an API
gulp.task('devServer', function() {
  connect.server({
    root: filePath.build.dest,
    fallback: filePath.build.dest + '/index.html',
    port: 5000,
    livereload: true
    // ,
    // middleware: function(connect, o) {
    //     return [ (function() {
    //         var url = require('url');
    //         var proxy = require('proxy-middleware');
    //         var options = url.parse('http://localhost:3000/'); // path to your dev API
    //         options.route = '/api';
    //         return proxy(options);
    //     })() ];
    // }
  });
});

gulp.task('prodServer', function() {
  connect.server({
    root: filePath.build.dest,
    fallback: filePath.build.dest + '/index.html',
    port: 5050,
    livereload: true
    // ,
    // middleware: function(connect, o) {
    //     return [ (function() {
    //         var url = require('url');
    //         var proxy = require('proxy-middleware');
    //         var options = url.parse('https://api-staging.your-domain.com/'); // path to your staging API
    //         options.route = '/api';
    //         return proxy(options);
    //     })() ];
    // }
  });
});


// =======================================================================
// Clean out dist folder contents on build
// =======================================================================
gulp.task('clean-dev', function () {
  del(['./dist/*.js', './dist/*.css', '!./dist/vendor.js', '!./dist/vendor.css', './dist/*.html', './dist/*.png', './dist/*.ico']);
});

gulp.task('clean-full', function () {
  del(['./dist/*']);
});


// =======================================================================
// JSHint
// =======================================================================
gulp.task('lint', function() {
  return gulp.src(filePath.lint.src)
  .pipe(jshint())
  .pipe(jshint.reporter(stylish));
});


// =======================================================================
// Browserify Bundle
// =======================================================================
gulp.task('bundle-dev', function() {

  var entryFile = filePath.browserify.src,
    bundler = watchify(entryFile);

  function rebundle () {
    return bundler.bundle({ debug: true })
    .pipe(source('bundle.js'))
    .on('error', handleError)
    .pipe(buffer())
    .pipe(sourcemaps.init({loadMaps: true}))
    .pipe(sourcemaps.write('./'))
    .pipe(gulp.dest(filePath.build.dest))
    .pipe(notify({ message: 'Browserify task complete' }))
    .pipe(connect.reload());
  }

  bundler.on('update', rebundle);

  return rebundle();
});

gulp.task('bundle-prod', function() {

  var entryFile = filePath.browserify.src,
    bundler = watchify(entryFile);

  function rebundle () {
    return bundler.bundle({ debug: true })
    .pipe(source('bundle.js'))
    .on('error', handleError)
    .pipe(buffer())
    .pipe(streamify(uglify({mangle: false})))
    .pipe(gulp.dest(filePath.build.dest))
    .pipe(notify({ message: 'Browserify task complete' }))
    .pipe(connect.reload());
  }

  bundler.on('update', rebundle);

  return rebundle();
});


// =======================================================================
// Styles Task
// =======================================================================
gulp.task('styles-dev', function () {
  return gulp.src(filePath.styles.src)
  .pipe(sourcemaps.init())
  .pipe(stylus({use: nib()}))
  .on('error', handleError)
  .pipe(sourcemaps.write())
  .pipe(gulp.dest(filePath.build.dest))
  .on('error', handleError)
  .pipe(notify({ message: 'Styles task complete' }))
  .pipe(connect.reload());
});

gulp.task('styles-prod', function () {
  return gulp.src(filePath.styles.src)
  .pipe(stylus({use: nib()}))
  .on('error', handleError)
  .pipe(prefix('last 1 version', '> 1%', 'ie 8', 'ie 7'), {map: true})
  .pipe(minifyCSS())
  .pipe(gulp.dest(filePath.build.dest))
  .on('error', handleError)
  .pipe(notify({ message: 'Styles task complete' }));
});


// =======================================================================
// Images Task
// =======================================================================
gulp.task('images', function() {
  return gulp.src(filePath.images.src)
  .on('error', handleError)
  .pipe(gulp.dest(filePath.images.dest))
  .pipe(notify({ message: 'Images copied' }))
  .pipe(connect.reload());
});


// =======================================================================
// Vendor JS Task
// =======================================================================
gulp.task('vendorJS', function () {
  return gulp.src(filePath.vendorJS.src)
  .pipe(concat('vendor.js'))
  .on('error', handleError)
  .pipe(uglify())
  .pipe(gulp.dest(filePath.build.dest))
  .pipe(notify({ message: 'VendorJS task complete' }));
});


// =======================================================================
// Vendor CSS Task
// =======================================================================
gulp.task('vendorCSS', function () {
  return gulp.src(filePath.vendorCSS.src)
  .pipe(concat('vendor.css'))
  .on('error', handleError)
  .pipe(minifyCSS())
  .pipe(gulp.dest(filePath.build.dest))
  .pipe(notify({ message: 'VendorCSS task complete' }))
  .pipe(connect.reload());
});


// =======================================================================
// Vendor FONTS Task
// =======================================================================
gulp.task('vendorFONTS', function () {
  return gulp.src(filePath.vendorFONTS.src)
  .on('error', handleError)
  .pipe(gulp.dest(filePath.vendorFONTS.dest))
  .pipe(notify({ message: 'Fonts copied' }))
  .pipe(connect.reload());
});

// =======================================================================
// Vendor IMG Task
// =======================================================================
gulp.task('vendorIMG', function () {
  return gulp.src(filePath.vendorIMG.src)
  .on('error', handleError)
  .pipe(gulp.dest(filePath.vendorIMG.dest))
  .pipe(notify({ message: 'img copied' }))
  .pipe(connect.reload());
});


// =======================================================================
// Copy index.html
// =======================================================================
gulp.task('copyIndex', function () {
  return gulp.src(filePath.copyIndex.src)
  .pipe(gulp.dest(filePath.build.dest))
  .pipe(notify({ message: 'index.html successfully copied' }))
  .pipe(connect.reload());
});


// =======================================================================
// Copy Favicon
// =======================================================================
gulp.task('copyFavicon', function () {
  return gulp.src(filePath.copyFavicon.src)
  .pipe(gulp.dest(filePath.build.dest))
  .pipe(notify({ message: 'favicon successfully copied' }));
});


// =======================================================================
// Watch for changes
// =======================================================================
gulp.task('watch', function () {
  gulp.watch(filePath.browserify.watch, ['bundle-dev']);
  gulp.watch(filePath.styles.watch, ['styles-dev']);
  gulp.watch(filePath.images.watch, ['images']);
  gulp.watch(filePath.vendorJS.src, ['vendorJS']);
  gulp.watch(filePath.vendorCSS.src, ['vendorCSS']);
  gulp.watch(filePath.copyIndex.watch, ['copyIndex']);
  gutil.log('Watching...');
});


// =======================================================================
// Sequential Build Rendering
// =======================================================================

// run "gulp" in terminal to build the DEV app
gulp.task('build-dev', function(callback) {
  runSequence(
    ['clean-dev', 'lint'],
    // images and vendor tasks are removed to speed up build time. Use "gulp build" to do a full re-build of the dev app.
    ['bundle-dev', 'styles-dev', 'copyIndex', 'copyFavicon'],
    ['devServer', 'watch'],
    callback
  );
});

// run "gulp prod" in terminal to build the PROD-ready app
gulp.task('build-prod', function(callback) {
  runSequence(
    ['clean-full', 'lint'],
    ['bundle-prod', 'styles-prod', 'images',
      'vendorJS', 'vendorCSS', 'vendorFONTS', 'vendorIMG',
      'copyIndex', 'copyFavicon'],
    ['prodServer'],
    callback
  );
});

// run "gulp build" in terminal for a full re-build in DEV
gulp.task('build', function(callback) {
  runSequence(
    ['clean-full', 'lint'],
    ['bundle-dev', 'styles-dev', 'images',
      'vendorJS', 'vendorCSS', 'vendorFONTS', 'vendorIMG',
      'copyIndex', 'copyFavicon'],
    ['devServer', 'watch'],
    callback
  );
});


gulp.task('default',['build-dev']);
gulp.task('prod',['build-prod']);
