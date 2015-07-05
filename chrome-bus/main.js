#!/usr/bin/env node

var Chrome = require('chrome-remote-interface');

var tabCache = {};
var delim='||'

function event(name, id, oldTab, newTab) {
    // new
    // closed
    // urlchanged
    // focuschanged
    if (name === 'new') {
        console.log('new' + delim + id + delim + tabAttributes(newTab));
    } else if (name === 'closed') {
        console.log('closed' + delim + id + delim + tabAttributes(oldTab));
    } else if (name === 'urlchanged') {
        console.log('urlchanged' + delim + id + delim + oldTab.url + delim + newTab.url);
    } else if (name === 'focuschanged') {
        console.log('focuschanged' + delim + id + delim + oldTab.focused + delim + newTab.focused);
    }
}

function tabAttributes(tab) {
    return tab.url + delim + tab.type + delim + tab.focused;
}

function main() {
    Chrome.listTabs(function (err, tabs) {
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
        }
    });
}

main();
setInterval(function() {
    main();
}, 4000);
