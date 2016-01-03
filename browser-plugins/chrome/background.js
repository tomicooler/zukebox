chrome.runtime.onInstalled.addListener(function() {
  chrome.declarativeContent.onPageChanged.removeRules(undefined, function() {
    chrome.declarativeContent.onPageChanged.addRules([
    {
      conditions: [
        new chrome.declarativeContent.PageStateMatcher({
          pageUrl: { urlContains: 'youtube.com/watch?v=' },
        })
      ],
      actions: [ new chrome.declarativeContent.ShowPageAction() ]
    }
    ]);
  });
});

chrome.pageAction.onClicked.addListener(function(tab) {
  sendTrack(tab);
});

function sendTrack(tab) {
    chrome.storage.sync.get({
        address: 'http://10.50.1.13:5000',
        user: "Stranger"
    }, function(items) {
        postTrack(tab.url, tab.id, items.address, items.user);
    });
}

function postTrack(url, tab_id, address, user) {
    var http = new XMLHttpRequest();
    http.open("POST", address + "/player/tracks", true);
    http.setRequestHeader("Content-type", "application/json");

    chrome.pageAction.setIcon({path:"progress.png", tabId:tab_id});

    http.onreadystatechange = function() {
        if (http.readyState == 4) {
            if (http.status != 201) {
                chrome.pageAction.setIcon({path:"error.png", tabId:tab_id});
                alert(http.statusText);
            } else {
                chrome.pageAction.setIcon({path:"done.png", tabId:tab_id});
            }
        }
    };

    var data = {
        url: url,
        user: user
    };

    http.send(JSON.stringify(data));
}
