$(function () {

    function ChromePlugin() {
        function checkFunctionParameter(methodName, callback) {
            if (typeof callback !== 'function') new Error(methodName + ': The parameter is not a function.');
        }

        this.getOptions = function (callback) {
            checkFunctionParameter('getOptions', callback);
            chrome.storage.sync.get({
                address: '',
                user: ''
            }, callback);
        };

        this.getActiveTabUrl = function (callback) {
            checkFunctionParameter('getOptions', callback);
            chrome.tabs.query({
                active: true,
                lastFocusedWindow: true
            }, function (tabs) {
                var tab = tabs[0];
                callback(tab.url);
            });
        };
    }

    var chromePlugin = new ChromePlugin();

    /**
     * Init
     */
    var init = false;
    $(window).on('focus', function () {
        if (init) return;
        init = true;
        chromePlugin.getOptions(function (options) {
            var zukeBox = new ZukeBox(options);

            chromePlugin.getActiveTabUrl(function (url) {
                zukeBox.setUrl(url);
            });

            function refresh() {
                zukeBox.refreshPlayer();
                zukeBox.refreshQueue();
            }

            setInterval(refresh, 1000);
            refresh();
        });
    });
});
