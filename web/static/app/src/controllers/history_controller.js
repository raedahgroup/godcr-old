import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, splitAmountIntoParts, truncate } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'stickyTableHeader', 'historyTable',
      'txRowTemplate',
      'errorMessage',
      'previousPageButton', 'nextPageButton',
      'pageReport',
      'loadingIndicator',
      'transactionType',
      'transactionCountContainer', 'transactionCount', 'transactionTotalCount',
    ]
  }

  connect () {
    window.addEventListener('resize', this.alignTableHeaderWithStickyHeader.bind(this))
    this.alignTableHeaderWithStickyHeader()
  }

  alignTableHeaderWithStickyHeader () {
    this.stickyTableHeaderTarget.style.width = `${this.historyTableTarget.clientWidth}px`

    if (this.historyTableTarget.innerHTML === '') {
      return
    }
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

  changeTransactionType () {
    this.historyTableTarget.innerHTML = ''
    this.transactionCountTarget.setAttribute('data-transactionCount', 0)
    hide(this.transactionCountContainerTarget)
    this.nextPage = 1
    this.fetchMoreTxs()
  }

  fetchMoreTxs () {
    show(this.loadingIndicatorTarget)

    let transType = this.transactionTypeTarget.value

    const _this = this
    axios.get(`/next-history-page?page=${this.nextPage}&trans-type=${transType}`)
      .then(function (response) {
        let result = response.data
        if (result.success) {
          if (!result.txs || result.txs.length === 0) {
            _this.setErrorMessage('No Transaction Found')
            return
          }

          let transCount = parseInt(_this.transactionCountTarget.getAttribute('data-transactionCount')) + result.txs.length
          _this.transactionCountTarget.textContent = transCount
          _this.transactionCountTarget.setAttribute('data-transactionCount', transCount)
          _this.transactionTotalCountTarget.textContent = result.txCount
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

    const accountName = tx => {
      let accountNames = []
      if (tx.direction === 1) {
        tx.outputs.forEach(output => {
          if (output.account_numberr === -1 || accountNames.indexOf(output.account_name) !== -1) {
            return
          }
          accountNames.push(output.account_name)
        })
      } else {
        tx.inputs.forEach(input => {
          if (inputs.account_numberr === -1 || accountNames.indexOf(input.account_name) !== -1) {
            return
          }
          accountNames.push(input.account_name)
        })
      }

      return accountNames.join(', ')
    }

    const _this = this
    txs.forEach(tx => {
      const txRow = document.importNode(_this.txRowTemplateTarget.content, true)
      const fields = txRow.querySelectorAll('td')

      fields[0].innerText = accountName(tx)
      fields[1].innerText = tx.long_time
      fields[2].innerText = tx.type
      fields[3].innerText = txDirection(tx.direction)

      const amountParts = splitAmountIntoParts(amountDcr(tx.amount))
      fields[4].innerHTML = `${amountParts[0]}<span>${amountParts[1]}${amountParts[2]}</span>`

      const feeParts = splitAmountIntoParts(amountDcr(tx.fee))
      fields[5].innerHTML = `${feeParts[0]}<span>${feeParts[1]}${feeParts[2]}</span>`

      fields[6].innerText = tx.status
      fields[7].innerHTML = `<a href="/transaction-details/${tx.hash}">${truncate(tx.hash, 10)}}</a>`

      _this.historyTableTarget.appendChild(txRow)
    })
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    show(this.errorMessageTarget)
  }
}
