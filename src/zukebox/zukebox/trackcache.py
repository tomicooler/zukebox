"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""

import os


class TrackCache:
    def __init__(self, base_path: str='/tmp/zukebox'):
        self.base_path = base_path

    def info_path(self, track_id: str):
        return self._get_path(track_id, 'json')

    def track_path(self, track_id: str):
        return self._get_path(track_id, 'mp3')

    def is_cached(self, track_id: str):
        return os.path.isfile(self.info_path(track_id)) and os.path.isfile(self.track_path(track_id))

    def _get_path(self, track_id: str, ext: str):
        return os.path.join(self.base_path, '{id}.{ext}'.format(id=track_id, ext=ext))

