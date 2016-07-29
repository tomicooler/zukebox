pragma Singleton
import QtQuick 2.5
import QtAV 1.6
import QuickFlux 1.0
import "../actions"

AppListener {

    property alias player: mediaPlayer

    MediaPlayer {
        id: mediaPlayer
        autoPlay: true
        property bool isPlaying: playbackState === MediaPlayer.PlayingState
    }

    Filter {
        type: ActionTypes.playPauseVideo
        onDispatched: {
            if (player.isPlaying)
                player.pause();
            else
                player.play();
        }
    }

    Filter {
        type: ActionTypes.seekInVideo
        onDispatched: {
            player.seek(message.position);
        }
    }
}
