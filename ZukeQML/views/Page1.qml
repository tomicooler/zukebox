import QtQuick 2.7
import "../stores"

Page1Form {

    property ListModel model

    listView1.delegate: trackDelegate
    listView1.model: model

    Component {
        id: trackDelegate

        Track {
            width: listView1.width
            height: rowLayout1.height
        }
    }
}
