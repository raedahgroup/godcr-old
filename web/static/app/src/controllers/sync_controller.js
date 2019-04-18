import { Controller } from 'stimulus'
import { hide, show } from '../utils'
import ws from '../services/messagesocket_service'

export default class extends Controller {
  static get targets () {
    return [
      'progressbar', 'totalSyncProgress', 'totalTimeRemaining',
      'showDetailsButton', 'syncDetails',
      'step1', 'fetchedHeadersCount', 'totalHeadersToFetch', 'headersFetchProgress', 'daysBehind',
      'step2', 'addressDiscoveryProgress',
      'step3', 'currentRescanHeight', 'rescanProgress',
      'connectedPeers', 'networkType'
    ]
  }

  connect () {
    ws.registerEvtHandler('updateSyncProgress', syncInfo => {
      this.progressbarTarget.style.width = `${syncInfo.totalSyncProgress}%`

      this.totalSyncProgressTarget.textContent = `${syncInfo.totalSyncProgress}% completed`
      if (syncInfo.totalTimeRemaining !== '') {
        this.totalSyncProgressTarget.textContent += `, ${syncInfo.totalTimeRemaining} remaining`
      }
      this.totalSyncProgressTarget.textContent += '.'

      switch (syncInfo.currentStep) {
        case 0: // fetching headers
          this.fetchedHeadersCountTarget.textContent = syncInfo.fetchedHeadersCount
          this.totalHeadersToFetchTargets.forEach(totalHeadersToFetchTarget => {
            totalHeadersToFetchTarget.textContent = syncInfo.totalHeadersToFetch
          })
          this.headersFetchProgressTarget.textContent = syncInfo.headersFetchProgress

          if (syncInfo.DaysBehind !== '') {
            this.daysBehindTarget.textContent = `Your wallet is ${syncInfo.daysBehind} behind.`
            show(this.daysBehindTarget)
          } else {
            hide(this.daysBehindTarget)
          }

          show(this.step1Target)
          hide(this.step2Target)
          hide(this.step3Target)
          break

        case 1: // discoverign used addresses
          this.addressDiscoveryProgressTarget.textContent = syncInfo.addressDiscoveryProgress

          hide(this.step1Target)
          show(this.step2Target)
          hide(this.step3Target)
          break

        case 2: // scanning block headers
          this.currentRescanHeightTarget.textContent = syncInfo.currentRescanHeight
          this.totalHeadersToFetchTargets.forEach(totalHeadersToFetchTarget => {
            totalHeadersToFetchTarget.textContent = syncInfo.totalHeadersToFetch
          })
          this.rescanProgressTarget.textContent = syncInfo.rescanProgress
          hide(this.step1Target)
          hide(this.step2Target)
          show(this.step3Target)
          break

        default:
          hide(this.step1Target)
          hide(this.step2Target)
          hide(this.step3Target)
      }

      this.connectedPeersTarget.textContent = syncInfo.connectedPeers
      this.networkTypeTarget.textContent = syncInfo.networkType

      if (syncInfo.done) {
        window.location.reload(true)
      }
    })
  }

  showDetails () {
    hide(this.showDetailsButtonTarget)
    show(this.syncDetailsTarget)
  }
}
