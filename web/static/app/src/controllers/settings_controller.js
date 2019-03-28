import { Controller } from 'stimulus'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage'
    ]
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
    this.hide(this.errorsTarget)
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
