"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""

from multiprocessing import Pool

import json

from gtts import gTTS
from tempfile import NamedTemporaryFile
import os

from zukebox.player import Player
from zukebox.trackcache import TrackCache
from zukebox.youtube import Youtube, DownloadError


tracks = []
recent_tracks = []
track_ids_being_downloaded = []
current_track = {}
cache = TrackCache()
youtube = Youtube()
_pool = None


def on_track_finished_callback():
    recent_tracks.insert(0, current_track.copy())
    if len(recent_tracks) > 100:
        del recent_tracks[len(recent_tracks) - 1]
    current_track.clear()
    play_next_track()


# Seriously, I hate python
player = Player(on_track_finished_callback)


def pool():
    # Workaround for pickle error
    global _pool
    if not _pool:
        _pool = Pool(4)
    return _pool


def shutdown():
    if _pool:
        _pool.close()
        _pool.join()


def play_next_track():
    if len(tracks) == 0 or player.playing:
        return

    track = tracks[0]
    track_id = track['id']
    if not cache.is_cached(track_id):
        return

    current_track.clear()
    current_track.update(track)
    del tracks[0]

    if 'message' in track and 'lang' in track:
        try:
            tts = gTTS(text=track['message'], lang=track['lang'])
            f = NamedTemporaryFile()
            tts.write_to_fp(f)
            f.file.flush()
            # TODO: find a good python library for playing sounds, VLC is not good enough
            os.system(' '.join(['playsound', '--volume', str(float(player.volume) / float(100.0)),
                                f.name, '&>', '/dev/null']))
            f.close()
        except Exception as e:
            print("Error: {}".format(e))

    path = cache.track_path(track_id)
    player.open(path)
    player.playing = True

    def touch(name, times=None):
        with open(name, 'a'):
            os.utime(name, times)

    touch(path)
    cache.clean_up()


def async_download_track(track: dict):
    url = track['url']
    track_id = track['id']

    track_path = cache.track_path(track_id)
    info_path = cache.info_path(track_id)

    try:
        youtube.download_audio(url, track_path)
        with open(info_path, 'w') as outfile:
            json.dump(track, outfile)
    except DownloadError as e:
        print(str(e))

    return track


def track_downloaded(track: dict):
    track_ids_being_downloaded.remove(track['id'])
    play_next_track()


def create_track(youtube_url: str, user: str, message: str, lang: str) -> dict:
    youtube_id = youtube.get_id(youtube_url)
    track_id = 'YOUTUBE-{id}'.format(id=youtube_id)

    if cache.is_cached(track_id):
        cache.info_path(track_id)
        with open(cache.info_path(track_id)) as track_info_file:
            track = json.load(track_info_file)
        track['user'] = user
        if len(message) > 0 and len(lang) > 0:
            track['message'] = message
            track['lang'] = lang
        tracks.append(track)
        play_next_track()
    else:
        url = 'https://www.youtube.com/watch?v={id}'.format(id=youtube_id)
        track = youtube.extract_info(url)
        track['id'] = track_id
        track['url'] = url

        if track_id not in track_ids_being_downloaded:
            track_ids_being_downloaded.append(track_id)
            pool().apply_async(async_download_track, args=[track], callback=track_downloaded)

        track['user'] = user
        if len(message) > 0 and len(lang) > 0:
            track['message'] = message
            track['lang'] = lang

        tracks.append(track)

    return track


def is_item_exist(items: [], index: int):
    return 0 <= index < len(items)
