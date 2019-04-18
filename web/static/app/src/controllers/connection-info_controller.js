import { Controller } from 'stimulus'
import { hide, show } from '../utils'
import ws from '../services/messagesocket_service'

export default class extends Controller {
  static get targets () {
    return [
      'totalBalance',
      'peersConnected',
      'latestBlock',
      'networkType',
      'blockScanProgress',
    ]
  }

  connect () {
    const blocksRescanProgress = localStorage.getItem('blocks-rescan-progress')
    if (blocksRescanProgress) {
      this.blockScanProgressTarget.textContent = blocksRescanProgress
    }

    ws.registerEvtHandler('updateConnInfo', data => {
      this.peersConnectedTarget.textContent = data.peersConnected
      this.totalBalanceTarget.textContent = data.totalBalance
      this.latestBlockTarget.textContent = data.latestBlock
      this.networkTypeTarget.textContent = data.networkType
    })

    ws.registerEvtHandler('updateBalance', data => {
      this.totalBalanceTarget.textContent = data.total
    })

    ws.registerEvtHandler('updateSyncProgress', syncInfo => {
      // hide the persistent blocks rescan progress section if this is the initial sync on server start (i.e. !syncInfo.done)
      // or if block headers rescan has not started or has completed (i.e. syncInfo.rescanProgress <= 0 || syncInfo.rescanProgress >= 100)
      if (!syncInfo.done || syncInfo.rescanProgress <= 0 || syncInfo.rescanProgress >= 100) {
        hide(this.blockScanProgressTarget)
        localStorage.removeItem('blocks-rescan-progress')
        return
      }

      this.blockScanProgressTarget.textContent = `Rescanning block headers ${syncInfo.rescanProgress}%. `
      this.blockScanProgressTarget.textContent += `Scanning ${syncInfo.currentRescanHeight} of ${syncInfo.totalHeadersToFetch} block headers.`
      show(this.blockScanProgressTarget)

      localStorage.setItem('blocks-rescan-progress', this.blockScanProgressTarget.textContent)
    })
  }
}
