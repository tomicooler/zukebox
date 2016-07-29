import QtQuick 2.4
import "../actions"

TrackForm {
    thumbnail.source: thumbnailA
    title.text: titleA
    user.text: userA

    toolButton1.onClicked: AppActions.lofasz(title.text);
}
