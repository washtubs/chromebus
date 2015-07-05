#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var id;
if (process.argv.length <= 2) {
    console.log('Missing id argument');
    return;
} else {
    id = process.argv[2];
}

Chrome.closeTab({id: id}, function(err, tab) {
    if (!err) {
        console.log('closing tab ' + id);
    } else {
        console.error('error closing tab ' + id);
    }
});
