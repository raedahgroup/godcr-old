import { Controller } from 'stimulus'
import axios from 'axios'
import {hide} from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'sourceAccount', 'numberOfTickets', 'spendUnconfirmed', 'errors', 'submitButton',
      // from wallet passphrase modal (utils.html)
      'walletPassphrase', 'passwordError'
    ]
  }

  validateForm () {
    this.errorsTarget.innerHTML = ''
    let formValid = true

    if (this.sourceAccountTarget.value === '') {
      this.showError('The source account is required')
      formValid = false
    }

    if (this.numberOfTicketsTarget.value === '') {
      this.showError('The number of tickets is required')
      formValid = false
    }

    return formValid
  }

  submitForm () {
    if (!this.validatePassphrase()) {
      return
    }

    $('#passphrase-modal').modal('hide')

    this.submitButtonTarget.innerHTML = 'Purchasing...'
    this.submitButtonTarget.setAttribute('disabled', 'disabled')

    let postData = $('#purchase-tickets-form').serialize()
    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      postData += '&source-account=' + this.sourceAccountTarget.value
    }

    // clear password input
    this.walletPassphraseTarget.value = ''

    let _this = this
    axios.post('/purchase-tickets', postData).then((response) => {
      let result = response.data
      if (!result.success) {
        _this.setErrorMessage(result.message)
      } else {
        var successMsg = ['<p>You have purchased ' + result.message.length + ' ticket(s)</p>']
        var ticketHashes = result.message.map(ticketHash => '<p><strong>' + ticketHash + '</strong></p>')
        successMsg.push(...ticketHashes)
        _this.setSuccessMessage(successMsg.join(''))
        _this.submitButtonTarget.innerHTML = 'Purchase again'
      }
    }).catch(() => {
      _this.submitButtonTarget.innerHTML = 'Purchase'
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      _this.submitButtonTarget.removeAttribute('disabled')
    })
  }

  validatePassphrase () {
    if (this.walletPassphraseTarget.value === '') {
      this.passwordErrorTarget.innerHTML = '<div class="error">Your wallet passphrase is required</div>'
      return false
    }

    return true
  }

  getWalletPassphraseAndSubmit () {
    this.clearMessages()
    if (!this.validateForm()) {
      return
    }
    $('#passphrase-modal').modal()
  }

  setErrorMessage (message) {
    hide(this.successMessageTarget)
    show(this.errorMessageTarget)
    this.errorMessageTarget.innerHTML = message
  }

  setSuccessMessage (message) {
    hide(this.errorMessageTarget)
    show(this.successMessageTarget)
    this.successMessageTarget.innerHTML = message
  }

  clearMessages () {
    hide(this.errorMessageTarget)
    hide(this.successMessageTarget)
    hide(this.errorsTarget)
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }

  showError (error) {
    this.errorsTarget.innerHTML += `<div class="error">${error}</div>`
    show(this.errorsTarget)
  }
}
