import { Controller } from 'stimulus'
import axios from 'axios'
import { clearMessages, copyToClipboard, hide, setErrorMessage, show, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage',
      'account',
      'generatedAddressContainer', 'generatedAddress', 'qrCodeImage',
      'generateNewAddressButton'
    ]
  }

  getCurrentAddress () {
    this.generateAddress(false)
  }

  copyAddressToClipboard () {
    copyToClipboard(this.generatedAddressTarget.textContent)
    showSuccessNotification('Copied to clipboard')
  }

  generateNewAddress () {
    this.generateAddress(true)
  }

  generateAddress (newAddress) {
    hide(this.generatedAddressContainerTarget)

    this.generateNewAddressButtonTarget.textContent = 'generating...'
    this.generateNewAddressButtonTarget.setAttribute('disabled', 'disabled')
    clearMessages(this)

    let url = '/generate-address/' + this.accountTarget.value
    if (newAddress) {
      url += '?new=yes'
    }

    const _this = this
    axios.get(url)
      .then((response) => {
        let result = response.data
        if (result.success) {
          _this.generatedAddressTarget.textContent = result.generatedAddress
          _this.qrCodeImageTarget.setAttribute('src', `data:image/png;base64,${result.qrCodeBase64Image}`)
          show(_this.generatedAddressContainerTarget)
        } else {
          setErrorMessage(_this, result.errorMessage)
        }
      })
      .catch(() => {
        setErrorMessage(_this, 'Unable to generate address. Something went wrong.')
      })
      .then(function () {
        _this.generateNewAddressButtonTarget.textContent = 'generate new address'
        _this.generateNewAddressButtonTarget.removeAttribute('disabled')
      })
  }
}
