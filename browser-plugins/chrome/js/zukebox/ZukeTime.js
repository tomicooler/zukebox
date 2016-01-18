function ZukeTime(dom) {

    function format(sec) {
        if (!sec) return '00:00:00';

        var date = new Date(sec * 1000);

        var hh = date.getUTCHours();
        var mm = date.getUTCMinutes();
        var ss = date.getUTCSeconds();

        if (hh < 10) hh = '0' + hh;
        if (mm < 10) mm = '0' + mm;
        if (ss < 10) ss = '0' + ss;

        return hh + ':' + mm + ':' + ss;
    }

    this.set = function (duration, time) {
        dom.$trackTime.text(format(time) + ' / ' + format(duration));
    };
}
