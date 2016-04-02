"""
zukebox: Main module

Copyright 2015, Tamas Domok
Licensed under MIT.
"""

from flask import Flask, jsonify, request, make_response, abort

import zukebox.zukebox as zb

app = Flask(__name__)


def _get_track(index: int):
    if not zb.is_item_exist(zb.tracks, index):
        abort(404, {'error': "Track does not exists; index='{index}'".format(index=index)})
    return zb.tracks[index]


def _get_recent_track(index: int):
    if not zb.is_item_exist(zb.recent_tracks, index):
        abort(404, {'error': "History track does not exists; index='{index}'".format(index=index)})
    return zb.recent_tracks[index]


def _ensure_json_contains_a_string_key_value_pair(key: str):
    if key not in request.json:
        abort(422, {'error': "Missing parameter from request; parameter='{parameter}'".format(parameter=key)})
    if not isinstance(request.json[key], str):
        abort(422, {'error': "Parameter is not a string; parameter='{parameter}'".format(parameter=key)})


def _ensure_json_value_is_integer(key: str):
    try:
        int(request.json[key])
    except ValueError:
        abort(422, {'error': "Invalid parameter value for '{key}', must be an integer".format(key=key)})


@app.route('/player/tracks', methods=['GET'])
def get_tracks():
    return jsonify({"tracks": zb.tracks})


@app.route('/player/tracks/<int:index>', methods=['GET'])
def get_track(index):
    return jsonify(_get_track(index))


@app.route('/player/tracks/<int:index>', methods=['DELETE'])
def delete_task(index):
    zb.tracks.remove(_get_track(index))
    return make_response('', 204)


@app.route('/player/tracks', methods=['POST'])
def create_track():
    if not request.json:
        abort(400, {'error': 'Not a json'})
    _ensure_json_contains_a_string_key_value_pair('url')
    _ensure_json_contains_a_string_key_value_pair('user')

    message = request.json.get('message', '')
    lang = request.json.get('lang', '')

    track = zb.create_track(request.json['url'],
                            request.json['user'],
                            message, lang)

    return make_response(jsonify(track), 201)


@app.route('/player/recent-tracks', methods=['GET'])
def get_recent_tracks():
    return jsonify({"tracks": zb.recent_tracks})


@app.route('/player/recent-tracks/<int:index>', methods=['GET'])
def get_recent_track(index):
    return jsonify(_get_recent_track(index))


@app.route('/player/control', methods=['GET'])
def get_control():
    return jsonify({
        'playing': zb.player.playing,
        'volume': zb.player.volume,
        'time': int(zb.player.position * zb.current_track[
            'duration']) if zb.player.playing and 'duration' in zb.current_track else 0,
        'track': zb.current_track,
    })


@app.route('/player/control', methods=['PATCH'])
def patch_control():
    if not request.json:
        abort(400, {'error': 'Not a json'})

    if 'playing' in request.json:
        if request.json['playing'] not in (True, False):
            abort(422, {'error': "Invalid parameter value for 'playing', must be true or false"})

        playing = request.json['playing'] == True
        zb.player.playing = playing

    if 'volume' in request.json:
        _ensure_json_value_is_integer('volume')

        volume = max(0, min(int(request.json['volume']), 100))
        zb.player.volume = volume

    if 'time' in request.json:
        _ensure_json_value_is_integer('time')

        if not zb.player.playing:
            abort(422, {'error': "ZukeBox is not playing, seeking is impossible"})

        duration = float(zb.current_track['duration'])
        time = float(max(0, min(float(request.json['time']), duration)))

        zb.player.position = float(time / duration)

    return get_control()


@app.errorhandler(400)
def bad_request(error):
    return make_response(jsonify(error.description), 400)


@app.errorhandler(404)
def not_found(error):
    return make_response(jsonify(error.description), 404)


@app.errorhandler(422)
def unprocessable_entity(error):
    return make_response(jsonify(error.description), 422)


def main():
    """
    Main function of the boilerplate code is the entry point of the 'zukebox' executable script (defined in setup.py).
    
    Use doctests, those are very helpful.
    
    >>> main()
    Hello
    >>> 2 + 2
    4
    """
    try:
        app.run(host='0.0.0.0', debug=True)
    except KeyboardInterrupt:
        zb.shutdown()
