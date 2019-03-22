import { Controller } from 'stimulus'
import ws from '../services/messagesocket_service'

import axios from 'axios'
// customTxRow
export default class extends Controller {
  static get targets () {
    return [
        'container',
      'totalBalance',
      'peersConnected',
      'lattestBlock',
      'networkType',
      'databasePath',
    ]
  }

  connect () {
    let _this = this

    ws.registerEvtHandler('newBlock', function (data) {
      _this.lattestBlockTarget = data.height
    })

    ws.registerEvtHandler('newPeer', function (data) {
      _this.peersConnectedBlockTarget = data.peersCount
    })
  }

  initialize () {
    let _this = this
    axios.get(url).then(function (response) {
      let result = response.data
      if (result.success) {
        let data = result.data
        _this.totalBalanceTarget.textContent = data.totalBalance
        _this.peersConnectedTarget.textContent = data.peersConnected
        _this.lattestBlockTarget.textContent = data.lattestBlock
        _this.networkTypeTarget.textContent = data.networkType
        _this.databasePathTarget.textContent = data.databasePath

        _this.show(this.containerTarget)
      } else {
        // TODO this message be shown in site wide error board
        console.log(result.message)
      }
    }).catch(function () {
      console.log('A server error occurred')
    })
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }
}
