#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var success;

var id;
if (process.argv.length <= 2) {
    console.error('Usage: activate-tab.js <tab-ID>');
    success = false;
    process.exit();
} else {
    id = process.argv[2];
}
Chrome.Activate({id: id}, function (err, tab) {
    if (!err) {
        success = true;
    } else {
        console.error('error activating tab ' + id);
        success = false;
    }
});
process.on('exit', function() { process.reallyExit(success ? 0 : 1) })
