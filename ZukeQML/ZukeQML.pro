QT += qml quick quickcontrols2

CONFIG += c++11

SOURCES += main.cc

RESOURCES += qml.qrc \
    resources.qrc

# Additional import path used to resolve QML modules in Qt Creator's code model
QML_IMPORT_PATH =

# Default rules for deployment.
include(deployment.pri)
include(quickflux/quickflux.pri)
