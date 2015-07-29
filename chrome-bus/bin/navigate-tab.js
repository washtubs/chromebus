#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var success;

var index, url;
if (process.argv.length <= 3) {
    console.error('Usage: navigate-tab.js <tab-INDEX> <url>');
    success = false;
    process.exit();
} else {
    index = process.argv[2];
    url = process.argv[3];
}

fn = function (chrome) {
    chrome.Page.loadEventFired(chrome.close);
    chrome.Page.enable();
    chrome.once('ready', function () {
        chrome.send(
            'Page.navigate',
            {'url': url},
            function(err, response) {
                console.error(response)
                chrome.close()
                if (err) {
                    console.error(response)
                    success = false
                } else {
                    success = true
                }
            }
        );
    });
};

Chrome({chooseTab: function() {return index;}}, fn).on('error', function () {
    console.error('Cannot connect to Chrome');
    success = false
});

process.on('exit', function() { process.reallyExit(success ? 0 : 1) })
