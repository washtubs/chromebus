#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var success;

var id;
if (process.argv.length <= 2) {
    console.error('Usage: new-tab.js <url>');
    success = false;
    process.exit();
} else {
    url = process.argv[2];
}

Chrome.New({url: url}, function(err) {
    if (!err) {
        success = true;
    } else {
        console.error('error creating new tab with url ' + url);
        success = false;
    }
});

process.on('exit', function() { process.reallyExit(success ? 0 : 1) })

