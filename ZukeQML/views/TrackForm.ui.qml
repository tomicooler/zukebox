import QtQuick 2.4
import QtQuick.Controls 2.0
import QtQuick.Layouts 1.0

Item {
    property alias title: title
    property alias user: user
    property alias thumbnail: thumbnail
    property alias rowLayout1: rowLayout1
    property alias toolButton1: toolButton1

    RowLayout {
        id: rowLayout1
        Image {
            id: thumbnail
            sourceSize.height: 80
            sourceSize.width: 120
            fillMode: Image.PreserveAspectFit
            Layout.preferredHeight: 80
            Layout.preferredWidth: 100
            Layout.fillWidth: true
            source: "qrc:/qtquickplugin/images/template_image.png"
        }

        Label {
            id: title
            text: qsTr("Label")
            Layout.maximumWidth: 200
            Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter
            elide: Text.ElideRight
        }

        Label {
            id: user
            text: qsTr("Label")
            Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter
            elide: Text.ElideRight
        }

        ToolButton {
            id: toolButton1
            text: qsTr("Tool Button")
        }
    }
}
