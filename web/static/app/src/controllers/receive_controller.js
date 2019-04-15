import { Controller } from 'stimulus'
import axios from 'axios'
import { clearMessages, copyToClipboard, hide, setErrorMessage, show, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage',
      'account',
      'addressContainer', 'address', 'image', 'copyButtonText',
      'generateNewAddressButton'
    ]
  }

  getReceiveAddress () {
    this.generateAddress(false)
  }

  copyAddressToClipboard () {
    copyToClipboard(this.addressTarget.textContent)
    showSuccessNotification('Copied to clipboard')
  }

  generateNewAddress () {
    this.generateAddress(true)
  }

  generateAddress (newAddress) {
    hide(this.addressContainerTarget)
    this.generateNewAddressButtonTarget.textContent = 'generating...'
    this.generateNewAddressButtonTarget.setAttribute('disabled', 'disabled')

    const _this = this
    clearMessages(this)
    let url = '/generate-address/' + this.accountTarget.value
    if (newAddress) {
      url += '?new=1'
    }
    axios.get(url)
      .then((response) => {
        let result = response.data
        if (result.success) {
          _this.addressTarget.textContent = result.address
          _this.imageTarget.setAttribute('src', result.imageData)
          show(_this.addressContainerTarget)
        } else {
          setErrorMessage(_this, result.message)
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
