function ZukeProgress(dom, $container) {

    var $bar = $container.find('.progress-bar');
    var $line = $container.find('.progress-line');
    var move = false;

    var left = 0;
    var maxWidth = 0;
    var callback;
    var lineWidthPercent = 0;

    function pos() {
        var lineWidth = intervalCorrection((event.pageX - left), 0, maxWidth);
        lineWidthPercent = valueToPercent(maxWidth, lineWidth);
        $line.width(lineWidthPercent + '%'); // % -> responsive!
    }

    function valueToPercent(max, current) {
        return ((current / max) * 100) | 0; // | 0 -> bitwise floor a number
    }

    function intervalCorrection(value, min, max) {
        return value < min ? min : value > max ? max : value;
    }

    $container.on('mousedown', function (event) {
        move = true;
        left = $bar.offset().left;
        maxWidth = $bar.width();
        pos(event);

        dom.$window.on('mousemove', function (event) {
            pos(event);
        });

        dom.$window.on('mouseup', function () {
            move = false;
            if (typeof callback === 'function') callback(lineWidthPercent);
            dom.$window.off('mousemove');
            dom.$window.off('mouseup');
        });
    });

    this.setPercent = function (percent) {
        if (move) return;
        percent = intervalCorrection(percent, 0, 100);
        $line.width(percent + '%'); // % -> responsive!
    };

    this.setValue = function (max, current) {
        if (move) return;
        var percent = valueToPercent(max, current);
        this.setPercent(percent);
    };

    this.onChange = function (call) {
        callback = call;
    };
}
