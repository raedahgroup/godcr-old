import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, truncate } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'selectedFilter',
      'transactionCountContainer', 'transactionCount', 'transactionTotalCount',
      'stickyTableHeader', 'historyTable',
      'txRowTemplate',
      'errorMessage',
      'previousPageButton', 'pageReport', 'nextPageButton',
      'loadingIndicator'
    ]
  }

  connect () {
    window.addEventListener('resize', this.alignTableHeaderWithStickyHeader.bind(this))
    this.alignTableHeaderWithStickyHeader()
  }

  alignTableHeaderWithStickyHeader () {
    if (this.historyTableTarget.innerHTML === '') {
      return
    }
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

  selectedFilterChanged () {
    this.historyTableTarget.innerHTML = ''
    hide(this.transactionCountContainerTarget)
    this.nextPage = 1
    this.fetchMoreTxs()
  }

  fetchMoreTxs () {
    show(this.loadingIndicatorTarget)

    const filter = this.selectedFilterTarget.value

    const _this = this
    axios.get(`/next-history-page?page=${this.nextPage}&filter=${filter}`)
      .then(function (response) {
        // since results are appended to the table, discard this response
        // if the user has changed the filter before the result is gotten
        if (_this.selectedFilterTarget.value !== filter) {
          return
        }
        let result = response.data
        if (result.success) {
          _this.transactionTotalCountTarget.textContent = result.transactionTotalCount
          show(_this.transactionCountContainerTarget)

          hide(_this.errorMessageTarget)
          _this.nextPage = result.nextPage
          _this.displayTxs(result.txs)

          _this.isLoading = false
          _this.checkScrollPos()
        } else {
          _this.setErrorMessage(result.message)
        }
      }).catch(function (e) {
        console.log(e)
        _this.setErrorMessage('A server error occurred')
      }).then(function () {
        _this.isLoading = false
        hide(_this.loadingIndicatorTarget)
      })
  }

  displayTxs (txs) {
    const directions = ['Sent', 'Received', 'Yourself']
    const txDirection = (direction) => {
      if (direction >= 0 && direction < directions.length) {
        return directions[direction]
      }
      return 'Unclear'
    }
    const amountDcr = (amount) => {
      return `${amount / 100000000} DCR`
    }

    const accountName = (tx) => {
      let accountNames = new Set()
      if (tx.direction === 1) {
        tx.outputs.forEach(output => {
          if (parseInt(output.previous_account) !== -1) {
            accountNames.add(output.account_name)
          }
        })
      } else {
        tx.inputs.forEach(input => {
          if (parseInt(input.previous_account) !== -1) {
            accountNames.add(input.account_name)
          }
        })
      }

      return Array.from(accountNames).join(', ')
    }

    const directionImage = (tx) => {
      switch (tx.direction) {
        case 0:
          return 'ic_send.svg'
        case 1:
          return 'ic_receive.svg'
      }
      if (tx.type === 'Ticket') {
        return 'live_ticket.svg'
      }
      return 'ic_tx_transferred.svg'
    }

    const _this = this

    txs.forEach(tx => {
      const txRow = document.importNode(_this.txRowTemplateTarget.content, true)
      const fields = txRow.querySelectorAll('td')

      fields[0].innerText = accountName(tx)
      fields[1].innerText = tx.long_time
      fields[2].innerText = tx.type

      const direction = txDirection(tx.direction)
      const image = directionImage(tx).toString()
      fields[3].innerHTML = '<img style="width: 15px" src="/static/images/' + image + '"> ' + direction

      fields[4].innerHTML = amountDcr(tx.amount)
      fields[5].innerHTML = amountDcr(tx.fee)

      fields[6].innerText = tx.status
      fields[7].innerHTML = `<a href="/transaction-details/${tx.hash}">${truncate(tx.hash, 10)}</a>`

      _this.historyTableTarget.appendChild(txRow)
    })

    _this.transactionCountTarget.textContent = _this.historyTableTarget.childElementCount
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    show(this.errorMessageTarget)
  }
}
