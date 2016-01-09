function isEmpty(ob){
   for(var i in ob){ return false;}
  return true;
}

function refreshPlayer() {
    chrome.storage.sync.get({
        address: 'http://10.50.1.13:5000',
        user: "Stranger"
    }, function(items) {
        var http = new XMLHttpRequest();
        http.open("GET", items.address + "/player/control", true);
        http.setRequestHeader("Content-type", "application/json");

        http.onreadystatechange = function() {
            if (http.readyState == 4) {
                if (http.status != 200) {
                    alert(http.statusText);
                } else {
                    response = JSON.parse(http.response);

                    var title = "Not playing";
                    var progress = 0;
                    var thumbnail = "music.png";
                    var progressLabel = "";

                    if (!isEmpty(response.track)) {
                        progress = (response.time / response.track.duration * 100);
                        thumbnail = response.track.thumbnail;
                        progressLabel = (new Date).clearTime().addSeconds(response.time).toString('HH:mm:ss');
                        title = response.track.title;
                    }

                    $('#current_track').text(title);
                    $('#thumbnail').attr("src", thumbnail);
                    $('#progress').width(progress + '%');
                    $('#progress-label').text(progressLabel);
                    $('#volume').text(response.volume);

                    $('#playing').attr('class', response.playing ? 'glyphicon glyphicon-play' : 'glyphicon glyphicon-pause');
                }
            }
        };

        http.send();
    });
}

function refreshPlayerLoop() {
    refreshPlayer();
    setTimeout(refreshPlayer, 1000 * 1);
}

function refreshQueue() {
    chrome.storage.sync.get({
        address: 'http://10.50.1.13:5000',
        user: "Stranger"
    }, function(items) {
        var http = new XMLHttpRequest();
        http.open("GET", items.address + "/player/tracks", true);
        http.setRequestHeader("Content-type", "application/json");

        http.onreadystatechange = function() {
            if (http.readyState == 4) {
                if (http.status != 200) {
                    alert(http.statusText);
                } else {
                    response = JSON.parse(http.response);

                    $("#queue").empty();
                    jQuery.each(response.tracks, function(index, track) {
                         $("#queue").append("<tr>" +
                         "<td><img width=\"64\" height=\"64\" src=\"music.png\" class=\"img-thumbnail\"></img></td>" +
                          "<td class=\"title\"></td>" +
                          "<td class=\"user\"></td>" +
                          "</tr>");

                          // For escaping, believe me I have no idea what I'm doing
                          $("tr:last .img-thumbnail").attr('src', track.thumbnail);
                          $("tr:last .title").text(track.title);
                          $("tr:last .user").text(track.user);
                         if (index > 3) return false;
                    });
                    $("#queue-length").text(response.tracks.length + " tracks in the queue.");
                }
            }
        };

        http.send();
    });
}

function refreshQueueLoop() {
    refreshQueue();
    setTimeout(refreshQueue, 1000 * 15);
}

function refreshUrl() {
    chrome.tabs.getSelected(null, function(tab) {
        if (tab.url.indexOf("youtube.com/watch?v=") > -1) {
            $("#send-div").show();
            $("#url").val(tab.url);
        } else {
            $("#send-div").hide();
        }
    });
}

function main() {
    refreshPlayerLoop();
    refreshQueueLoop();
    refreshUrl();
}

function sendTrack() {
    chrome.storage.sync.get({
        address: 'http://10.50.1.13:5000',
        user: "Stranger"
    }, function(items) {
        var http = new XMLHttpRequest();
        http.open("POST", items.address + "/player/tracks", true);
        http.setRequestHeader("Content-type", "application/json");

        http.onreadystatechange = function() {
            if (http.readyState == 4) {
                if (http.status != 201) {
                    alert(http.responseText);
                }
            }
        };

        var data = {
            url: $("#url").val(),
            user: items.user
        };

        http.send(JSON.stringify(data));

        refreshPlayer();
        refreshQueue();
    });
}

document.addEventListener('DOMContentLoaded', main);
document.getElementById('send').addEventListener('click', sendTrack);
