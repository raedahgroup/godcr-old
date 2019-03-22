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
  var protocol = (loc.protocol === 'https:') ? 'wss' : 'ws'
  return protocol + '://' + loc.host + '/ws'
}

function sleep (ms) {
  return new Promise(resolve => setTimeout(resolve, ms))
}

async function createWebSocket (loc) {
  // wait a bit to prevent websocket churn from drive by page loads
  var uri = getSocketURI(loc)
  await sleep(1000)
  ws.connect(uri)

/*var updateBlockData = function (event) {
    console.log('Received newblock message', event)
    var newBlock = JSON.parse(event)
    newBlock.block.unixStamp = new Date(newBlock.block.time).getTime() / 1000
    globalEventBus.publish('BLOCK_RECEIVED', newBlock)
  }
  ws.registerEvtHandler('newblock', updateBlockData)*/
}

createWebSocket(window.location)
