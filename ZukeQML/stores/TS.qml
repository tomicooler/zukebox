import QtQuick 2.5
import QuickFlux 1.0
import "../actions"
import "controlservice.js" as Control

Item {
    signal refresh();

    property Timer timer: Timer {
        interval: 1000
        repeat: true

        onTriggered: refresh();
    }

    Component.onCompleted: {
        timer.start();
        refresh();
    }
}
