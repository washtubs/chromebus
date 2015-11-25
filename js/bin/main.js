#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var tabCache = {};
var delim='||'
var tabDelim='|+|'

function event(name, id, oldTab, newTab) {
    // new
    // closed
    // urlchanged
    // focuschanged
    console.log(name + delim + id + delim + (oldTab ? tabAttributes(oldTab) : 'nil') + delim + (newTab ? tabAttributes(newTab) : 'nil'));
}

function tabAttributes(tab) {
    return tab.url + tabDelim + tab.type + tabDelim + tab.index + tabDelim + tab.focused;
}

function main() {
    Chrome.List(function (err, tabs) {
        if (!err) {
            for (var id in tabCache) {
                var found = false;
                for (i = 0; i < tabs.length; i++) {
                    if (id === tabs[i].id) {
                        found = true;
                        break;
                    }
                }
                if (!found) {
                    event('closed', id, tabCache[id], undefined);
                    delete tabCache[id];
                }
            }
            for (i = 0; i < tabs.length; i++ ) {
                tab = tabs[i];
                id = tab.id
                // construct tabc object
                newtabc = {
                    url: tab.url,
                    type: tab.type,
                    index: i,
                    focused: !!(i === 0)
                }
                if (!tabCache.hasOwnProperty(id)) {
                    event('new', id, undefined, newtabc);
                } else {
                    prevtabc = tabCache[tab.id];
                    if (prevtabc.url !== newtabc.url) {
                        event('urlchanged', id, prevtabc, newtabc);
                    }
                    if (newtabc.focused !== prevtabc.focused ) {
                        event('focuschanged', id, prevtabc, newtabc);
                    }
                }
                tabCache[id] = newtabc;
            }
        } else {
            console.error("Something went wrong. Failed to connect? Maybe chrome isn't running with the debug port");
            console.error(err);
            process.exit(1)
        }
    });
}

main();
setInterval(function() {
    main();
}, 4000);
