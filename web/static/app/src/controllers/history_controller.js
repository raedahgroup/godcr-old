import { Controller } from 'stimulus'
import axios from 'axios'
import '../css/history_style.scss'

export default class extends Controller {
  static get targets () {
    return ['historyTable', 'nextPageButton', 'errorMessage']
  }

  initialize () {
    // hide next page button to use infinite scroll
    this.hide(this.nextPageButtonTarget)
    this.nextBlockHeight = this.nextPageButtonTarget.getAttribute('data-next-block-height')
  }

  didScroll (e) {
    if (this.isLoading || !this.nextBlockHeight) {
      return
    }

    const element = e.target.scrollingElement
    const scrollPos = element.scrollTop + element.clientHeight
    if (scrollPos >= element.scrollHeight * 0.9) {
      this.isLoading = true
      this.fetchMoreTxs()
    }
  }

  fetchMoreTxs () {
    const _this = this
    axios.get(`/next-history-page?start=${this.nextBlockHeight}`)
      .then(function (response) {
        let result = response.data
        if (result.success) {
          _this.hide(_this.errorMessageTarget)
          _this.nextBlockHeight = result.nextBlockHeight
          _this.displayTxs(result.txs)
        } else {
          _this.setErrorMessage(result.message)
        }
      }).catch(function (e) {
        _this.setErrorMessage('A server error occurred')
      }).then(function () {
        _this.isLoading = false
      })
  }

  displayTxs (txs) {
    const directions = ['Sent', 'Received', 'Transferred']
    const txDirection = (direction) => {
      if (direction >= 0 && direction < directions.length) {
        return directions[direction]
      }
      return 'Unclear'
    }

    let n = this.historyTableTarget.querySelectorAll('tr').length

    const txRows = txs.map(tx => {
      return `<tr>
                  <td>${++n}</td>
                  <td>${tx.formatted_time}</td>
                  <td>${txDirection(tx.direction)}</td>
                  <td>${tx.amount / 100000000} DCCR</td>
                  <td>${tx.fee / 100000000} DCR</td>
                  <td>${tx.type}</td>
                  <td><a href="/transaction-details/${tx.hash}">${tx.hash}</a></td>
              </tr>`
    })

    this.historyTableTarget.innerHTML += txRows.join('')
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    this.show(this.errorMessageTarget)
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }
}
