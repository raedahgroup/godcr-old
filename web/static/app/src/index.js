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

function getSocketURI (loc) {
  let protocol = (loc.protocol === 'https:') ? 'wss' : 'ws'
  return protocol + '://' + loc.host + '/ws'
}

function createWebSocket (loc) {
  setTimeout(() => {
    // wait a bit to prevent websocket churn from drive by page loads
    let uri = getSocketURI(loc)
    ws.connect(uri)
  }, 1000)
}

createWebSocket(window.location)

const application = Application.start()
const context = require.context('./controllers', true, /\.js$/)
application.load(definitionsFromContext(context))
