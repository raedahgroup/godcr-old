import { Controller } from 'stimulus'
import axios from 'axios'
import { showErrorNotification, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return ['hideAccount', 'defaultAccount']
  }

  toggleHideAccount (e) {
    let accountElement = e.currentTarget
    let accountNumber = accountElement.getAttribute('data-account')

    const postData = (accountElement.checked) ? `hide-account=${accountNumber}` : `reveal-account=${accountNumber}`
    axios.put('/settings', postData).then((response) => {
      const result = response.data
      if (result.success) {
        showSuccessNotification('Changes saved successfully')
      } else {
        showErrorNotification(result.error ? result.error : 'Something went wrong, please try again later')
        accountElement.checked = !accountElement.checked
      }
    }).catch(() => {
      accountElement.checked = !accountElement.checked
      showErrorNotification('A server error occurred')
    })
  }

  updateDefaultAccount (e) {
    const defaultAccountEl = e.currentTarget
    const defaultAccount = defaultAccountEl.getAttribute('data-account')

    // uncheck all other accounts that were previously marked default
    this.defaultAccountTargets.forEach((el, i) => {
      if (el.checked && el.getAttribute('data-account') !== defaultAccount) {
        el.checked = false
      }
    })

    // post data
    const postData = `default-account=${defaultAccount}`
    axios.put('/settings', postData).then((response) => {
      const result = response.data
      if (result.success) {
        showSuccessNotification('Changes saved successfully')
      } else {
        showErrorNotification(result.error ? result.error : 'Something went wrong, please try again later')
        defaultAccountEl.checked = !defaultAccountEl.checked
      }
    }).catch(() => {
      defaultAccountEl.checked = !defaultAccountEl.checked
      showErrorNotification('A server error occurred')
    })
  }
}
