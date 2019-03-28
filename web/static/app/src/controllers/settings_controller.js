import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'oldPassword', 'oldPasswordError', 'newPassword', 'newPasswordError', 'confirmPassword', 'confirmPasswordError'
    ]
  }

  changePassword (e) {
    e.preventDefault()

    if (!this.validateChangePasswordValidationField()) {
      return
    }
    let submitBtn = e.currentTarget
    submitBtn.textContent = 'Busy...'
    submitBtn.setAttribute('disabled', true)

    let postData = $('#change-password-form').serialize()

    this.clearMessages()

    let _this = this

    axios.post('/change-password', postData).then((response) => {
      let result = response.data
      if (result.error !== undefined) {
        _this.setErrorMessage(result.error)
      } else {
        this.oldPasswordTarget.value = ''
        this.newPasswordTarget.value = ''
        this.confirmPasswordTarget.value = ''

        _this.setSuccessMessage('Password changed')
      }
    }).catch(() => {
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      submitBtn.innerHTML = 'Change Password'
      submitBtn.removeAttribute('disabled')
      $('#change-password-modal').modal('hide')
    })
  }

  validateChangePasswordValidationField () {
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
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }

  showError (error) {
    this.errorsTarget.innerHTML += `<div class="error">${error}</div>`
    this.show(this.errorsTarget)
  }
}
