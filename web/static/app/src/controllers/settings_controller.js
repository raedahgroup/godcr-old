import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, showErrorNotification, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'oldPassword', 'oldPasswordError', 'newPassword', 'newPasswordError', 'confirmPassword',
      'confirmPasswordError', 'changePasswordErrorMessage',
      'spendUnconfirmedFunds', 'showIncomingTransactionNotification', 'showNewBlockNotification',
      'changeCurrencyConverterErrorMessage', 'currencyConverterNone', 'currencyConverterBitrex', 'updateCurrencyConverterButton'
    ]
  }

  changePassword (e) {
    e.preventDefault()
    if (!this.validateChangePasswordFields()) {
      return
    }

    this.clearChangePasswordError()

    let submitBtn = e.currentTarget
    submitBtn.textContent = 'Changing Password...'
    submitBtn.setAttribute('disabled', true)

    let _this = this

    let postData = $('#change-password-form').serialize()
    axios.post('/change-password', postData).then((response) => {
      let result = response.data
      if (result.error) {
        _this.showChangePasswordError(result.error)
      } else {
        _this.oldPasswordTarget.value = ''
        _this.newPasswordTarget.value = ''
        _this.confirmPasswordTarget.value = ''

        showSuccessNotification('Password changed')
        $('#change-password-modal').modal('hide')
      }
    }).catch(() => {
      _this.showChangePasswordError('A server error occurred')
    }).then(() => {
      submitBtn.innerHTML = 'Change Password'
      submitBtn.removeAttribute('disabled')
    })
  }

  showChangePasswordError (message) {
    this.changePasswordErrorMessageTarget.textContent = message
    show(this.changePasswordErrorMessageTarget)
  }

  clearChangePasswordError () {
    this.changePasswordErrorMessageTarget.textContent = ''
    hide(this.changePasswordErrorMessageTarget)
  }

  validateChangePasswordFields () {
    this.clearChangePasswordValidationErrors()
    let isClean = true
    if (this.oldPasswordTarget.value === '') {
      this.oldPasswordErrorTarget.textContent = 'Old password is required'
      isClean = false
    }
    if (this.newPasswordTarget.value === '') {
      this.newPasswordErrorTarget.textContent = 'New password is required'
      isClean = false
    }
    if (this.confirmPasswordTarget.value === '') {
      this.confirmPasswordErrorTarget.textContent = 'Confirm password is required'
      isClean = false
    }
    if (this.confirmPasswordTarget.value !== this.newPasswordTarget.value) {
      if (this.confirmPasswordErrorTarget.textContent === '') {
        this.confirmPasswordErrorTarget.textContent = 'Confirm password doesn\'t match'
      }
      isClean = false
    }
    return isClean
  }

  clearChangePasswordValidationErrors () {
    this.oldPasswordErrorTarget.textContent = ''
    this.newPasswordErrorTarget.textContent = ''
    this.confirmPasswordErrorTarget.textContent = ''
    this.clearChangePasswordError()
  }

  updateSpendUnconfirmed () {
    const _this = this
    const postData = `spend-unconfirmed=${this.spendUnconfirmedFundsTarget.checked}`
    axios.put('/settings', postData).then((response) => {
      let result = response.data
      if (result.success) {
        showSuccessNotification('Changes saved successfully')
      } else {
        showErrorNotification(result.error ? result.error : 'Something went wrong, please try again later')
        _this.spendUnconfirmedFundsTarget.checked = !_this.spendUnconfirmedFundsTarget.checked
      }
    }).catch(() => {
      _this.spendUnconfirmedFundsTarget.checked = !_this.spendUnconfirmedFundsTarget.checked
      showErrorNotification('A server error occurred')
    })
  }

  updateShowIncomingTransactionNotification () {
    const _this = this
    const postData = `show-incoming-transaction-notification=${this.showIncomingTransactionNotificationTarget.checked}`
    axios.put('/settings', postData).then((response) => {
      let result = response.data
      if (result.success) {
        showSuccessNotification('Changes saved successfully')
      } else {
        showErrorNotification(result.error ? result.error : 'Something went wrong, please try again later')
        _this.showIncomingTransactionNotificationTarget.checked = !_this.showIncomingTransactionNotificationTarget.checked
      }
    }).catch(() => {
      _this.showIncomingTransactionNotificationTarget.checked = !_this.showIncomingTransactionNotificationTarget.checked
      showErrorNotification('A server error occurred')
    })
  }

  updateShowNewBlockNotification () {
    const _this = this
    const postData = `show-new-block-notification=${this.showNewBlockNotificationTarget.checked}`
    axios.put('/settings', postData).then((response) => {
      let result = response.data
      if (result.success) {
        showSuccessNotification('Changes saved successfully')
      } else {
        showErrorNotification(result.error ? result.error : 'Something went wrong, please try again later')
        _this.showNewBlockNotificationTarget.checked = !_this.showNewBlockNotificationTarget.checked
      }
    }).catch(() => {
      _this.showNewBlockNotificationTarget.checked = !_this.showNewBlockNotificationTarget.checked
      showErrorNotification('A server error occurred')
    })
  }

  updateCurrencyConverter () {
    const _this = this
    const postData = `currency-converter=${(this.currencyConverterBitrexTarget.checked ? 'bitrex' : 'none')}`
    axios.put('/settings', postData).then((response) => {
      let result = response.data
      if (result.error) {
        _this.changeCurrencyConverterErrorMessageTarget.textContent = result.error ? result.error : 'Something went wrong, please try again later'
        return
      }
      showSuccessNotification('Changes saved successfully')
      $('#currency-converter-modal').modal('hide')
    }).catch(() => {
      _this.changeCurrencyConverterErrorMessageTarget.textContent = 'A server error occurred'
    })
  }
}
