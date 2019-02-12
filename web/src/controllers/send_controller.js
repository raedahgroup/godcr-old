import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return ['sourceAccount', 'address', 'amount', 'destinations', 'destinationTemplate', 'changeAddress', 'changeAmount', 'errors', 'customInput', 'customTxRow', 'customInputContent', 'changeOutputsTarget', 'submitButton', 'nextButton', 'removeDestinationButton', 'form', 'walletPassphrase', 'passwordError', 'useCustom', 'spendUnconfirmed', 'errorMessage', 'successMessage', 'progressBar']
  }

  initialize () {
    this.newDestination()
  }

  validateSendForm () {
    this.errorsTarget.innerHTML = ''
    let errors = []
    let isClean = this.validateDestinationsField()

    if (this.sourceAccountTarget.value === '') {
      errors.push('The source account is required')
    }

    if (this.useCustomTarget.value && this.getSelectedInputsSum() < this.getTotalSendAmount()) {
      errors.push('The sum of selected inputs is less than send amount')
    }

    if (errors.length) {
      for (let i in errors) {
        this.showError(errors[i])
      }
      isClean = false
    }

    return isClean
  }

  validateAndRefreshPercentage () {
    if (this.validateDestinationsField()) {
      this.calculateSelectedInputPercentage()
    }
  }

  validateDestinationsField () {
    this.clearMessages()
    let isClean = true
    this.addressTargets.forEach((el, i) => {
      if (el.value === '') {
        this.showError('The destination address and amount are required')
        isClean = false
      }
    })
    this.amountTargets.forEach((el, i) => {
      let amount = parseFloat(el.value)
      if (!(amount > 0)) {
        this.showError('Amount must be a non-zero positive number')
        isClean = false
      }
    })
    return isClean
  }

  getSelectedInputsSum () {
    var sum = 0
    this.customTxRowTarget.querySelectorAll('input.custom-input:checked').forEach((el, i) => {
      sum += parseFloat(el.dataset.amount)
    })
    return sum
  }

  getTotalSendAmount () {
    let amount = 0
    this.amountTargets.forEach((el, i) => {
      amount += parseFloat(el.value)
    })
    return amount
  }

  calculateSelectedInputPercentage () {
    var sendAmount = this.getTotalSendAmount()
    var selectedInputSum = this.getSelectedInputsSum()
    var percentage = 0

    if (selectedInputSum >= sendAmount) {
      percentage = 100
    } else {
      percentage = (selectedInputSum / sendAmount) * 100
    }
    this.progressBarTarget.style.width = percentage + '%'
  }

  resetCustomizePanel () {
    this.customTxRowTarget.querySelectorAll('tbody').forEach(el => {
      el.innerHTML = ''
    })
    this.customTxRowTarget.querySelectorAll('.status').forEach(el => {
      this.show(el)
    })
    this.customTxRowTarget.querySelectorAll('.alert-danger').forEach(el => {
      el.parentNode.removeChild(el)
    })
  }

  openCustomizePanel () {
    let _this = this
    this.resetCustomizePanel()
    $('#custom-tx-row').slideDown()

    let accountNumber = _this.sourceAccountTarget.value
    let callback = function (txs) {
      // populate outputs
      let utxoHtml = txs.map(tx => {
        let receiveDateTime = new Date(tx.receive_time * 1000).toString().split(' ').slice(0, 5).join(' ')
        let dcrAmount = tx.amount / 100000000
        return `<tr>
                  <td width='5%'><input data-action='click->send#calculateSelectedInputPercentage' type='checkbox' class='custom-input' name='utxo' value='${tx.key}' data-amount='${dcrAmount}' /></td>
                  <td width='40%'>${tx.address}</td>
                  <td width='15%'>${dcrAmount} DCR</td>
                  <td width='25%'>${receiveDateTime}</td>
                  <td width='15%'>${tx.confirmations} confirmation(s)</td>
                </tr>`
      })
      _this.customInputContentTarget.innerHTML = utxoHtml
      $('#custom-tx-row .status').hide()
    }
    this.getUnspentOutputs(accountNumber, callback)
  }

  getUnspentOutputs (accountNumber, successCallback) {
    this.submitButtonTarget.innerHTML = 'Loading...'
    this.submitButtonTarget.setAttribute('disabled', 'disabled')
    let _this = this

    let url = '/unspent-outputs/' + accountNumber
    if (_this.spendUnconfirmedTarget.checked) {
      url += '?getUnconfirmed=true'
    }

    axios.get(url).then(function (response) {
      let result = response.data
      if (result.success) {
        successCallback(result.message)
      } else {
        _this.setErrorMessage(result.message)
      }
    }).catch(function () {
      _this.setErrorMessage('A server error occurred')
    }).then(function () {
      _this.submitButtonTarget.innerHTML = 'Next'
      _this.submitButtonTarget.removeAttribute('disabled')
    })
  }

  submitSendForm () {
    if (!this.walletPassphraseTarget.value) {
      this.setErrorMessage('Your wallet passphrase is required')
      return
    }
    $('#passphrase-modal').modal('hide')

    this.submitButtonTarget.innerHTML = 'Sending...'
    this.submitButtonTarget.setAttribute('disabled', 'disabled')

    var postData = $('#send-form').serialize()
    postData += '&totalSelectedInputAmountDcr=' + this.getSelectedInputsSum()

    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      postData += '&source-account=' + this.sourceAccountTarget.value
    }

    console.log(postData)

    let _this = this
    axios.post('/send', postData).then((response) => {
      let result = response.data
      if (result.error !== undefined) {
        _this.setErrorMessage(result.error)
      } else {
        let txHash = 'The transaction was published successfully. Hash: <strong>' + result.txHash + '</strong>'
        _this.setSuccessMessage(txHash)
      }
    }).catch((error) => {
      console.log(error)
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      _this.submitButtonTarget.innerHTML = 'Send'
      _this.submitButtonTarget.removeAttribute('disabled')
    })
  }

  newDestination () {
    if (!this.validateDestinationsField()) {
      return
    }
    let destinationTemplate = document.importNode(this.destinationTemplateTarget.content, true)
    this.destinationsTarget.appendChild(destinationTemplate)
    if (this.destinationsTarget.querySelectorAll('div.row').length > 1) {
      this.show(this.removeDestinationButtonTarget)
    }
  }

  removeDestination () {
    let count = this.destinationsTarget.querySelectorAll('div.row').length
    if (!(count > 1)) {
      return
    }
    this.destinationsTarget.removeChild(this.destinationsTarget.querySelector('div.row:last-child'))
    if (!(this.destinationsTarget.querySelectorAll('div.row').length > 1)) {
      this.hide(this.removeDestinationButtonTarget)
    }
  }

  validatePassphrase () {
    if (this.walletPassphraseTarget.value === '') {
      this.passwordErrorTarget.innerHTML = '<div class="error">Your wallet passphrase is required</div>'
      return false
    }

    return true
  }

  getWalletPassphraseAndSubmit () {
    if (!this.validateDestinationsField()) {
      return
    }
    $('#passphrase-modal').modal()
  }

  toggleUseCustom () {
    if (!this.useCustomTarget.checked) {
      $('#custom-tx-row').slideUp()
      this.resetCustomizePanel()
      return
    }
    if (!this.validateDestinationsField()) {
      this.useCustomTarget.checked = false
      return
    }
    this.clearMessages()
    this.openCustomizePanel()
  }

  toggleSpendUnconfirmed () {
    if (this.useCustomTarget.checked) {
      this.resetCustomizePanel()
      this.openCustomizePanel()
    }
  }

  setErrorMessage (message) {
    this.hide(this.successMessageTarget)
    this.show(this.errorMessageTarget)
    this.errorMessageTarget.innerHTML = message
  }

  setSuccessMessage (message) {
    this.hide(this.errorMessageTarget)
    this.show(this.successMessageTarget)
    this.successMessageTarget.innerHTML = message
  }

  clearMessages () {
    this.hide(this.errorMessageTarget)
    this.hide(this.successMessageTarget)
    this.errorsTarget.innerHTML = ''
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }

  showError (error) {
    this.errorsTarget.innerHTML += `<div class="error">${error}</div>`
  }
}
