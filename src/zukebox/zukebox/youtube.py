"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""

import youtube_dl
import re


class DownloadError(Exception):
    pass


class Logger(object):
    def debug(self, msg):
        pass

    def warning(self, msg):
        pass

    def error(self, msg):
        pass


class Youtube:

    options = {
        'format': 'bestaudio/best',
        'postprocessors': [{
            'key': 'FFmpegExtractAudio',
            'preferredcodec': 'mp3',
            'preferredquality': '192',
        }],
        'logger': Logger(),
    }

    def extract_info(self, youtube_url: str):
        try:
            with youtube_dl.YoutubeDL(self.options) as ydl:
                info = ydl.extract_info(youtube_url, download=False)

                if not info or 'title' not in info:
                    raise DownloadError("Could not extract song details; response='{response}'".format(response=str(info)))

                return {
                    'title': info.get('title'),
                    'duration': info.get('duration'),
                    'thumbnail': info.get('thumbnails', [{}])[0].get('url', ''),
                }
        except youtube_dl.DownloadError as e:
            raise DownloadError("Could not download song details; error='{error}'".format(error=str(e)))
        except:
            raise DownloadError("Could not download song details; error='unknown'")

    def download_audio(self, youtube_url: str, output: str):
        try:
            options = self.options.copy()
            options['outtmpl'] = output
            with youtube_dl.YoutubeDL(options) as ydl:
                ydl.download([youtube_url])
        except youtube_dl.DownloadError as e:
            raise DownloadError("Could not download song; error='{error}'".format(error=str(e)))
        except:
            raise DownloadError("Could not download song; error='unknown'")

    @classmethod
    def get_id(cls, youtube_url: str):
        result = re.match("^(?:http(?:s)?://)?(?:www\.)?(?:m\.)?(?:youtu\.be/|youtube\.com"
                          "/(?:(?:watch)?\?(?:.*&)?v(?:i)?=|(?:embed|v|vi|user)/))([^\?&\"'>]+)", youtube_url)
        return result.group(1) if result else None
