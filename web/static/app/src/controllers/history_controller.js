import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return ['stickyTableHeader', 'historyTable', 'nextPageButton', 'loadingIndicator', 'errorMessage']
  }

  connect () {
    window.addEventListener('resize', this.alignTableHeaderWithStickyHeader.bind(this))
    this.alignTableHeaderWithStickyHeader()
  }

  alignTableHeaderWithStickyHeader () {
    this.stickyTableHeaderTarget.style.width = `${this.historyTableTarget.clientWidth}px`

    // set column width on sticky header to match real table
    const tableColumns = this.historyTableTarget.querySelector('tr').querySelectorAll('td')
    const staticHeaderColumns = this.stickyTableHeaderTarget.querySelector('tr').querySelectorAll('th')
    staticHeaderColumns.forEach((headerColumn, index) => {
      headerColumn.style.width = `${tableColumns[index].clientWidth}px`
    })
  }

  initialize () {
    // hide next page button to use infinite scroll
    this.hide(this.nextPageButtonTarget)
    this.nextBlockHeight = this.nextPageButtonTarget.getAttribute('data-next-block-height')
    this.checkScrollPos()
  }

  checkScrollPos () {
    // check if there is space at the bottom to load more now
    this.windowScrolled({ target: document })
  }

  windowScrolled (e) {
    const element = e.target.scrollingElement
    const scrollTop = element.scrollTop
    this.makeTableHeaderSticky(scrollTop)

    if (this.isLoading || !this.nextBlockHeight) {
      return
    }

    const scrollPos = scrollTop + element.clientHeight
    if (scrollPos >= element.scrollHeight * 0.95) {
      this.isLoading = true
      this.fetchMoreTxs()
    }
  }

  makeTableHeaderSticky (scrollTop) {
    const historyTableOffset = this.historyTableTarget.parentElement.offsetTop
    if (this.stickyTableHeaderTarget.classList.contains('d-none') && scrollTop >= historyTableOffset) {
      this.show(this.stickyTableHeaderTarget)
    } else if (scrollTop < historyTableOffset) {
      this.hide(this.stickyTableHeaderTarget)
    }
  }

  fetchMoreTxs () {
    this.show(this.loadingIndicatorTarget)

    const _this = this
    axios.get(`/next-history-page?start=${this.nextBlockHeight}`)
      .then(function (response) {
        let result = response.data
        if (result.success) {
          _this.hide(_this.errorMessageTarget)
          _this.nextBlockHeight = result.nextBlockHeight
          _this.displayTxs(result.txs)

          _this.isLoading = false
          _this.checkScrollPos()
        } else {
          _this.setErrorMessage(result.message)
        }
      }).catch(function () {
        _this.setErrorMessage('A server error occurred')
      }).then(function () {
        _this.isLoading = false
        _this.hide(_this.loadingIndicatorTarget)
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
                  <td style="text-align: right">${tx.amount}</td>
                  <td style="text-align: right">${tx.fee}</td>
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
