var gulp = require('gulp');
var concat = require('gulp-concat');
var minifyCSS = require('gulp-minify-css');
var rename = require('gulp-rename');
var uglifyJS = require('gulp-uglify');

// Copy HTML
gulp.task('html', function() {
  return gulp.src('assets/html/*.html')
    .pipe(gulp.dest('public/'));
});

// Copy CSS
gulp.task('css', function() {
  return gulp.src(['assets/css/*.css', 'node_modules/bootstrap/dist/css/bootstrap.css'])
    .pipe(minifyCSS())
    .pipe(concat('app.css'))
    .pipe(gulp.dest('public/'));
});

// Copy JS
gulp.task('js', function() {
  return gulp.src(['assets/js/*.js', 'node_modules/bootstrap/dist/js/bootstrap.js'])
    .pipe(uglifyJS())
    .pipe(concat('app.js'))
    .pipe(gulp.dest('public/'));
});

// Default
gulp.task('default', ['html', 'css', 'js']);
