#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var success;

var id;
if (process.argv.length <= 2) {
    console.error('Usage: close-tab.js <tab-ID>');
    success = false;
    process.exit();
} else {
    id = process.argv[2];
}

Chrome.Close({id: id}, function(err, tab) {
    if (!err) {
        success = true;
    } else {
        console.error('error closing tab ' + id);
        success = false;
    }
});

process.on('exit', function() { process.reallyExit(success ? 0 : 1) })
