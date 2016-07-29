pragma Singleton
import QtQuick 2.5
import QuickFlux 1.0
import "../actions"
import "controlservice.js" as Control

AppListener {

    property ListModel tracks: ListModel { }
    property TS ts: TS {
        onRefresh: {
            Control.get_recent_tracks(tracks);
        }
    }

}
