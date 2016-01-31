function ZukeBox(options) {

    var that = this;
    var address = options.address || 'http://10.50.1.13:5000';
    var user = options.user || 'Stranger';
    var lang = options.lang || 'hu';
    var dom = new ZukeDom();
    var api = new ZukeApi();
    var control = new ZukeControl(dom);
    var time = new ZukeTime(dom);
    var trackProgress = new ZukeProgress(dom, dom.$trackProgress);
    var volumeProgress = new ZukeProgress(dom, dom.$volumeProgress);

    var currentDuration = null;
    trackProgress.onChange(function (percent) {
        if (currentDuration !== null) {
            var time = (currentDuration * (percent / 100)) | 0;
            api.patch(address + '/player/control', {
                time: time
            });
        }
    });

    var currentVolume = null;
    volumeProgress.onChange(function (percent) {
        if (currentVolume !== null) {
            setVolumeIcon(percent);
            api.patch(address + '/player/control', {
                volume: percent
            });
        }
    });

    function setVolumeIcon(value) {
        if (value == 0) {
            dom.$volumeIcon.attr('class', 'icon-volume-mute');
        } else if (value < 30) {
            dom.$volumeIcon.attr('class', 'icon-volume-low');
        } else if (value < 60) {
            dom.$volumeIcon.attr('class', 'icon-volume-medium');
        } else {
            dom.$volumeIcon.attr('class', 'icon-volume-high');
        }
    }

    control.onChange(function (status) {
        api.patch(address + '/player/control', {
            playing: status
        });
    });

    this.setUrl = function (url) {
        if (url.match(/.*youtube\.com\/watch\?v=.*/)) {
            dom.$sendUrl.val(url);
        }
    };

    this.send = function (url, message) {
        if (url.match(/.*youtube\.com\/watch\?v=.*/)) {
            api.post(address + '/player/tracks', {
                url: url,
                user: user,
                message: message || '',
                lang: lang || ''
            });
        }
    };

    this.refreshPlayer = function () {
        api.get(address + '/player/control').success(function (data) {
            if (!data.track) return;
            var track = data.track;

            if (track.thumbnail) dom.$trackThumbnail.attr('src', track.thumbnail);
            if (track.title) dom.$trackTitle.text(track.title);
            if (track.user) dom.$trackUser.text(track.user);
            if (track.url) dom.$trackUrl.attr('href', track.url);

            if (track.duration && data.time) {
                currentDuration = track.duration;
                trackProgress.setValue(currentDuration, data.time);

                time.set(track.duration, data.time);
            }

            if (data.volume >= 0) {
                currentVolume = data.volume;
                setVolumeIcon(currentVolume);
                volumeProgress.setPercent(currentVolume);
            }

            if (data.playing) control.setPlaying(data.playing);
        });
    };

    this.refreshQueue = function () {
        api.get(address + '/player/tracks').success(function (data) {
            var tracks = data.tracks;
            var length = tracks.length;
            var items = '';
            var itemsMaxSize = 3;

            dom.$queue.empty();

            for (var i = 0; i < length; ++i) {
                var track = tracks[i];

                items += '<a href="' + track.url + '" target="_blank" class="queue-item">'
                    + '<div class="queue-thumbnail"><img src="' + track.thumbnail + '" alt=""/></div>'
                    + '<div class="queue-title">' + track.title + '</div>'
                    + '<div class="queue-user">' + track.user + '</div>'
                    + '</a>';

                if (i === itemsMaxSize) break;
            }

            dom.$queue.append(items);

            dom.$queueInfo.text('( ' + length + ' tracks in the queue. )');
        });
    };

    dom.$sendButton.on('click', function () {
        that.send(dom.$sendUrl.val(), dom.$sendMessage.val());
    });

    dom.$trackThumbnail.attr('src', '/img/cover' + Math.floor((Math.random() * 3) + 1) + '.jpg');
}
