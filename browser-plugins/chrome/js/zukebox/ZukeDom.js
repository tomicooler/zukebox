function ZukeDom() {

    this.$window = $(window);

    this.$sendMessage = $('#send-message');
    this.$sendUrl = $('#send-url');
    this.$sendButton = $('#send-btn');

    this.$trackThumbnail = $('#track-thumbnail');
    this.$trackTitle = $('#track-title');
    this.$trackUser = $('#track-user');
    this.$trackUrl = $('#track-url');

    this.$trackProgress = $('#track-progress');
    this.$volumeProgress = $('#volume-progress');

    this.$volumeButton = $('#volume-btn');
    this.$volumeIcon = this.$volumeButton.find('#volume-ico');

    this.$queue = $('#queue');
    this.$queueInfo = $('#queue-info');

    this.$controlButton = $('#control-btn');
    this.$trackTime = $('#track-time');
}
