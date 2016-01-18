function ZukeControl(dom) {
    var playing = false;
    var callback;

    dom.$controlButton.on('click', toggle);

    function setIcon() {
        dom.$controlButton.attr('class', playing ? 'icon-play' : 'icon-pause');
    }

    function toggle() {
        playing = !playing;
        setIcon();
        callback(playing);
    }

    this.setPlaying = function (value) {
        playing = !!value;
        setIcon();
    };

    this.onChange = function (call) {
        callback = call;
    };
}
