var gulp = require("gulp"),
	fs = require("fs"),
	dateformat = require("dateformat"),
	sync = require("gulp-sync")(gulp),
	concat = require("gulp-concat"),
	rename = require("gulp-rename"),
	header = require("gulp-header"),
	base64 = require("gulp-base64"),
	templatecache = require("gulp-angular-templatecache"),
	usemin = require("gulp-usemin"),
	uglify = require("gulp-uglify"),
	minifyCss = require("gulp-minify-css"),
	minifyHtml = require("gulp-minify-html"),
	replace = require("gulp-replace"),
	connect = require("gulp-connect"),
	karma = require("gulp-karma");

var pkg = require("./package.json");
pkg.version = (""+fs.readFileSync("../../VERSION")).trim();

var banner = [
		"/* " + pkg.name + " v" + pkg.version + " " + dateformat(new Date(), "yyyy-mm-dd"),
		" * " + pkg.homepage,
		" * License: " + pkg.license,
		" */\n\n"
	].join("\n"),
	root = "files",
	paths = {
		"index.html": [root+"/index.html"],
		"index.src.html": [root+"/index.src.html"],
		"tmpl": [root+"/directives/*.html"],
		"js": ["!**/*.tmp.js", "!**/*.test.js",
			root+"/config/*.js", root+"/directives/*.js", root+"/services/*.js", root+"/index/*.js"],
		"css": [root+"/index/*.css"],
		"test": ["!**/*.tmp.js",
			root+"/lib/angular/angular.js",
			root+"/lib/angular-mocks/angular-mocks.js",
			root+"/lib/angular-ui-bootstrap-bower/ui-bootstrap-tpls.js",
			root+"/lib/angular-animate/angular-animate.js",
			root+"/lib/angular-local-storage/angular-local-storage.js",
			root+"/lib/angular-translate/angular-translate.js",
			root+"/config/*.js", root+"/directives/*.js", root+"/services/*.js", root+"/index/*.js"]
	};

gulp.task("build", sync.sync([
	["css", "js", "tmpl", "bower.json"],
	["index.src.html"],
	["index.html"]
]));

gulp.task("default", ["build", "watch", "connect"]);

gulp.task("css", function(done) {
	gulp.src(paths.css)
		.pipe(connect.reload())
		.on("end", done);
});

gulp.task("js", function(done) {
	gulp.src(paths.js)
		.pipe(connect.reload())
		.on("end", done);
});

gulp.task("tmpl", function(done) {
	gulp.src(paths.tmpl)
		.pipe(templatecache("angular-template.tmp.js", {
			module: "sysd",
			root: "directives"
		}))
		.pipe(gulp.dest(root+"/index/"))
		.pipe(connect.reload())
		.on("end", done);
});

gulp.task("index.src.html", function(done) {
	gulp.src(paths["index.src.html"])
		.pipe(rename({ basename: "index" }))
		.pipe(gulp.dest(root))
		.pipe(connect.reload())
		.on("end", done);
});

gulp.task("index.html", function(done) {
	gulp.src(paths["index.html"])
		.pipe(usemin({
			css: [
				minifyCss(),
				header(banner)
			],
			js: [
				replace(/\.version = \"0\";/, ".version = \"" + pkg.version + "\""),
				uglify(),
				header(banner)
			],
			html: [
				minifyHtml({ empty: true })
			]
		}))
		.pipe(gulp.dest(root))
		.pipe(connect.reload())
		.on("end", done);
});

gulp.task("bower.json", function(done) {
	gulp.src(["bower.json"])
		.pipe(replace(/"version": "[^"]*"/, "\"version\": \"" + pkg.version + "\""))
		.pipe(gulp.dest("./"))
		.on("end", done);
});

gulp.task("watch", ["build"], function() {
	["tmpl", "css", "js", "index.src.html"].forEach(function(i) {
		gulp.watch(paths[i], [i]);
	});
});

gulp.task("connect", ["build"], function() {
	connect.server({
		root: root,
		port: 9000,
		livereload: true
	});
});

gulp.task("test", function(done) {
	gulp.src(paths.test)
		.pipe(karma({
			configFile: "karma.config.js",
			action: "run"
		}))
		.on("error", function(err) {
			console.log(err);
			this.emit("end");
		})
		.on("end", done);
});
