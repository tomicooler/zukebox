var BASE = "http://10.30.255.175:5000";

function request(verb, endpoint, obj, cb) {
    //print('request: ' + verb + ' ' + BASE + (endpoint?'/' + endpoint:''))
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        //print('xhr: on ready state change: ' + xhr.readyState)
        if(xhr.readyState === XMLHttpRequest.DONE) {
            if(cb) {
                var res = JSON.parse(xhr.responseText.toString())
                cb(res);
            }
        }
    }
    xhr.open(verb, BASE + (endpoint?'/' + endpoint:''));
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.setRequestHeader('Accept', 'application/json');
    var data = obj?JSON.stringify(obj):''
    xhr.send(data)
}


function get_tracks(tracks) {
    _get_tracks("player/tracks", tracks);
}

function get_recent_tracks(tracks) {
    _get_tracks("player/recent-tracks", tracks);
}

function _get_tracks(path, tracks) {
    request('GET', path, null, function(resp) {
        tracks.clear();
        var zuke_tracks = resp.tracks;
        for(var i = 0; i < zuke_tracks.length; i++) {
            var tr = zuke_tracks[i];
            var track = {
                thumbnailA: tr.thumbnail,
                titleA: tr.title,
                userA: tr.user
            }

            tracks.append(track);
        }
    });
}

function delete_track(index) {
    request('DELETE', "player/tracks/" + index, null, null);
}
