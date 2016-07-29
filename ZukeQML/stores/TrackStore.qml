pragma Singleton
import QtQuick 2.5
import QuickFlux 1.0
import "../actions"
import "controlservice.js" as Control

AppListener {

    property ListModel tracks: ListModel { }
    property TS ts: TS {
        onRefresh: {
            Control.get_tracks(tracks);
        }
    }

    Filter {
        type: ActionTypes.lofasz
        onDispatched: {
            for (var i = 0 ; i < tracks.count ; i++) {
                var item  = tracks.get(i);
                if (item.titleA === message.title) {
                    Control.delete_track(i);
                    tracks.remove(item);
                }
            }
        }
    }

}
