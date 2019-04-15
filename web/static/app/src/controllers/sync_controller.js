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
      this.progressbarTarget.style.width = `${syncInfo.TotalSyncProgress}%`

      this.totalSyncProgressTarget.textContent = `${syncInfo.TotalSyncProgress}% completed`
      if (syncInfo.TotalTimeRemaining !== '') {
        this.totalSyncProgressTarget.textContent += `, ${syncInfo.TotalTimeRemaining} remaining`
      }
      this.totalSyncProgressTarget.textContent += '.'

      switch (syncInfo.CurrentStep) {
        case 1:
          this.fetchedHeadersCountTarget.textContent = syncInfo.FetchedHeadersCount
          this.totalHeadersToFetchTargets.forEach(totalHeadersToFetchTarget => {
            totalHeadersToFetchTarget.textContent = syncInfo.TotalHeadersToFetch
          })
          this.headersFetchProgressTarget.textContent = syncInfo.HeadersFetchProgress

          if (syncInfo.DaysBehind !== '') {
            this.daysBehindTarget.textContent = `Your wallet is ${syncInfo.DaysBehind} behind.`
            show(this.daysBehindTarget)
          } else {
            hide(this.daysBehindTarget)
          }

          show(this.step1Target)
          hide(this.step2Target)
          hide(this.step3Target)
          break

        case 2:
          this.addressDiscoveryProgressTarget.textContent = syncInfo.AddressDiscoveryProgress

          hide(this.step1Target)
          show(this.step2Target)
          hide(this.step3Target)
          break

        case 3:
          this.currentRescanHeightTarget.textContent = syncInfo.CurrentRescanHeight
          this.totalHeadersToFetchTargets.forEach(totalHeadersToFetchTarget => {
            totalHeadersToFetchTarget.textContent = syncInfo.TotalHeadersToFetch
          })
          this.rescanProgressTarget.textContent = syncInfo.RescanProgress
          hide(this.step1Target)
          hide(this.step2Target)
          show(this.step3Target)
          break

        default:
          hide(this.step1Target)
          hide(this.step2Target)
          hide(this.step3Target)
      }

      this.connectedPeersTarget.textContent = syncInfo.ConnectedPeers
      this.networkTypeTarget.textContent = syncInfo.NetworkType

      if (syncInfo.Done) {
        window.location.reload(true)
      }
    })
  }

  showDetails () {
    hide(this.showDetailsButtonTarget)
    show(this.syncDetailsTarget)
  }
}
