"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""
import threading

import vlc


class Player:
    def __init__(self, on_track_finished_callback=None):
        self.on_track_finished_callback = on_track_finished_callback
        self.instance = vlc.Instance()
        self.mediaplayer = self.instance.media_player_new()
        self.vlc_events = self.mediaplayer.event_manager()
        self.vlc_events.event_attach(vlc.EventType.MediaPlayerEndReached, self.track_finished, 1)
        self.vlc_events.event_attach(vlc.EventType.MediaPlayerEncounteredError, self.track_finished, 1)
        # mediaplayer.is_playing() is not "synchronous"
        # e.g: calling it immediately after pause() results in True
        self.is_playing = False

    def open(self, filename: str):
        media = self.instance.media_new(filename)
        self.mediaplayer.set_media(media)

    @property
    def volume(self) -> int:
        return self.mediaplayer.audio_get_volume()

    @volume.setter
    def volume(self, volume: int):
        self.mediaplayer.audio_set_volume(volume)

    @property
    def position(self) -> float:
        return self.mediaplayer.get_position()

    @position.setter
    def position(self, position: float):
        self.mediaplayer.set_position(position)

    @property
    def playing(self) -> bool:
        return self.is_playing

    @playing.setter
    def playing(self, play: bool):
        if play:
            started = self.mediaplayer.play()
            self.is_playing = True if started == 0 else False
        else:
            if self.is_playing:  # Pause would restart it otherwise
                self.mediaplayer.pause()
                self.is_playing = False

    @vlc.callbackmethod
    def track_finished(self, *args, **kwargs):
        self.is_playing = False
        if self.on_track_finished_callback:
            threading.Timer(1, self.on_track_finished_callback).start()
