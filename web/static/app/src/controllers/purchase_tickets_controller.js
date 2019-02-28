import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return ['sourceAccount', 'errors', 'spendUnconfirmed', 'numberOfTickets', 'walletPassphrase', 'successMessage', 'errorMessage', 'passwordError', 'submitButton']
  }

  validateForm () {
    this.errorsTarget.innerHTML = ''
    let errors = []
    let isClean = true

    if (this.sourceAccountTarget.value === '') {
      errors.push('The source account is required')
    }

    if (this.numberOfTicketsTarget.value === '') {
      errors.push('The number of tickets is required')
    }

    if (errors.length) {
      for (let i in errors) {
        this.showError(errors[i])
      }
      isClean = false
    }

    return isClean
  }

  submitForm () {
    if (!this.validatePassphrase()) {
      return
    }

    $('#passphrase-modal').modal('hide')
    this.submitButtonTarget.innerHTML = 'Purchasing...'
    this.submitButtonTarget.setAttribute('disabled', 'disabled')

    var postData = $('#purchase-tickets-form').serialize()

    // clear password input
    this.walletPassphraseTarget.value = ''

    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      postData += '&source-account=' + this.sourceAccountTarget.value
    }

    let _this = this
    axios.post('/purchase-tickets', postData).then((response) => {
      let result = response.data
      if (!result.success) {
        _this.setErrorMessage(result.message)
      } else {
        var successMsg = ['<p>You have purchased ' + result.message.length + ' ticket(s)</p>']
        var ticketHashes = result.message.map(ticketHash => '<p><strong>' + ticketHash + '</strong></p>')
        successMsg.push(ticketHashes)
        _this.setSuccessMessage(successMsg.join(''))
      }
    }).catch(() => {
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      _this.submitButtonTarget.innerHTML = 'Purchase'
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
