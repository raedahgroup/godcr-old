import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'stickyTableHeader', 'historyTable',
      'txRowTemplate',
      'errorMessage',
      'previousPageButton', 'nextPageButton',
      'pageReport',
      'loadingIndicator'
    ]
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
    hide(this.previousPageButtonTarget)
    hide(this.nextPageButtonTarget)
    hide(this.pageReportTarget)

    this.nextPage = this.nextPageButtonTarget.getAttribute('data-next-page')
    if (this.nextPage) {
      // check if there is space at the bottom to load more now
      this.checkScrollPos()
    }
  }

  checkScrollPos () {
    const scrollElement = document.body
    const scrollTop = scrollElement.scrollTop
    this.makeTableHeaderSticky(scrollTop)

    if (this.isLoading || !this.nextPage) {
      return
    }

    const scrollContentHeight = document.documentElement.scrollHeight
    const delta = scrollContentHeight - window.outerHeight - 10
    if (window.scrollY >= delta * 0.9) {
      this.isLoading = true
      this.fetchMoreTxs()
    }
  }

  windowScrolled () {
    this.checkScrollPos()
  }

  makeTableHeaderSticky (scrollTop) {
    const historyTableOffset = this.historyTableTarget.parentElement.offsetTop
    if (this.stickyTableHeaderTarget.classList.contains('d-none') && scrollTop >= historyTableOffset) {
      show(this.stickyTableHeaderTarget)
    } else if (scrollTop < historyTableOffset) {
      hide(this.stickyTableHeaderTarget)
    }
  }

  fetchMoreTxs () {
    show(this.loadingIndicatorTarget)

    const _this = this
    axios.get(`/next-history-page?page=${this.nextPage}`)
      .then(function (response) {
        let result = response.data
        if (result.success) {
          hide(_this.errorMessageTarget)
          _this.nextPage = result.nextPage
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
        hide(_this.loadingIndicatorTarget)
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
    const amountDcr = (amount) => {
      return `${amount / 100000000} DCR`
    }

    const _this = this
    txs.forEach(tx => {
      const txRow = document.importNode(_this.txRowTemplateTarget.content, true)
      const fields = txRow.querySelectorAll('td')

      fields[0].innerText = tx.long_time
      fields[1].innerText = tx.type
      fields[2].innerText = txDirection(tx.direction)
      fields[3].innerText = amountDcr(tx.amount)
      fields[4].innerText = amountDcr(tx.fee)
      fields[5].innerText = tx.status
      fields[6].innerHTML = `<a href="/transaction-details/${tx.hash}">${tx.hash}</a>`

      _this.historyTableTarget.appendChild(txRow)
    })
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    show(this.errorMessageTarget)
  }
}
