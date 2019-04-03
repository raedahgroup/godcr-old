import { Application } from 'stimulus'
import { definitionsFromContext } from 'stimulus/webpack-helpers'
import '../node_modules/toastr/build/toastr.css'
import ws from './services/messagesocket_service'
import './css/style.scss'
import { library, dom } from '@fortawesome/fontawesome-svg-core'
import { faCog, faCopy } from '@fortawesome/free-solid-svg-icons'

library.add(faCog, faCopy)
dom.watch()

const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))

function getSocketURI () {
  let protocol = (window.location.protocol === 'https:') ? 'wss' : 'ws'
  return `${protocol}://${window.location.host}/ws`
}

function createWebSocket () {
  setTimeout(() => {
    // wait a bit to prevent websocket churn from drive by page loads
    let uri = getSocketURI()
    ws.connect(uri)
  }, 1000)
}

createWebSocket()

const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
