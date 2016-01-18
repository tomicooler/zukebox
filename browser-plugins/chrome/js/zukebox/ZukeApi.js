function ZukeApi() {
    function ajax(method, url, data) {
        return $.ajax({
            method: method,
            contentType: 'application/json',
            url: url,
            dataType: 'json',
            data: data || '',
            error: function (jqXHR, textStatus) {
                console.log(textStatus);
            }
        });
    }

    this.get = function(url) {
        return ajax('GET', url || '');
    };

    this.post = function(url, data) {
        return ajax('POST', url, JSON.stringify(data));
    };

    this.patch = function (url, data) {
        return ajax('PATCH', url, JSON.stringify(data));
    }
}
