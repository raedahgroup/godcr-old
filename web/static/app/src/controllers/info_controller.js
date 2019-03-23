import { Controller } from 'stimulus'
import ws from '../services/messagesocket_service'

import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return [
      'container',
      'totalBalance',
      'peersConnected',
      'latestBlock',
      'networkType',
      'dbDir'
    ]
  }

  connect () {
    let _this = this

    ws.registerEvtHandler('updateConnInfo', function (data) {
      _this.totalBalanceTarget.textContent = data.totalBalance
      _this.peersConnectedTarget.textContent = data.peersConnected
      _this.latestBlockTarget.textContent = data.latestBlock
      _this.networkTypeTarget.textContent = data.networkType
      _this.dbDirTarget.textContent = data.dbDir

      _this.show(_this.containerTarget)
    })
  }

  initialize () {
    let _this = this
    axios.get('/connection-info').then(function (response) { // TODO use the actual url
      let data = response.data
      _this.totalBalanceTarget.textContent = data.totalBalance
      _this.peersConnectedTarget.textContent = data.peersConnected
      _this.latestBlockTarget.textContent = data.latestBlock
      _this.networkTypeTarget.textContent = data.networkType
      _this.dbDirTarget.textContent = data.dbDir

      _this.show(_this.containerTarget)
    }).catch(function (e) {
      console.log('Error in updating connection info:', e)
      _this.hide(_this.containerTarget)
    })
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }
}
