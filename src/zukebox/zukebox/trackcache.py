"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""

import os
import glob


class TrackCache:
    def __init__(self, base_path: str = '/tmp/zukebox', cache_size: int = 1000):
        self.base_path = base_path
        self.cache_size = cache_size

    def info_path(self, track_id: str):
        return self._get_path(track_id, 'json')

    def track_path(self, track_id: str):
        return self._get_path(track_id, 'mp3')

    def clean_up(self):
        def mb(byte):
            return byte / 1000 / 1000

        try:
            current_size = mb(
                sum(os.path.getsize(os.path.join(self.base_path, f)) for f in os.listdir(self.base_path) if
                    os.path.join(self.base_path, f)))

            clean_up_size = max(0, current_size - self.cache_size)
            if clean_up_size > 0:
                least_accessed = sorted(glob.iglob(self.base_path + '/*.mp3'), key=os.path.getatime)
                for file in least_accessed:
                    song_path = os.path.join(self.base_path, file)
                    info_path = song_path[:-3] + 'json'
                    print(song_path + " - " + info_path)
                    clean_up_size -= mb(os.path.getsize(song_path))
                    clean_up_size -= mb(os.path.getsize(info_path))
                    os.remove(song_path)
                    os.remove(info_path)
                    if clean_up_size <= 0:
                        break
        except:
            pass

    def is_cached(self, track_id: str):
        return os.path.isfile(self.info_path(track_id)) and os.path.isfile(self.track_path(track_id))

    def _get_path(self, track_id: str, ext: str):
        return os.path.join(self.base_path, '{id}.{ext}'.format(id=track_id, ext=ext))
