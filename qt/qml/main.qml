import QtQuick 2.6
import QtQuick.Controls 2.1

ApplicationWindow {
	id: window
	visible: true
	title: "GoDCR"
    minimumWidth: 600
    minimumHeight: 400

    Column {
        anchors.left: parent
        topPadding: 35
        leftPadding: 15

        Row {
            spacing: 40
            Button {
                text: "Check Balance"
            }
        }
        Row {
            spacing: 25
            Button {
                text: "Receive"
            }
        }
        Row {
            spacing: 25
            Button {
                text: "Send"
            }
        }

        Component.onCompleted: {
            window.x = (Screen.width - window.width) / 2
            window.y = (Screen.height - window.height) / 2
        }
    }
}
