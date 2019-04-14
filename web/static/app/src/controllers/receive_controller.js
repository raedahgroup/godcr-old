import { Controller } from 'stimulus'
import axios from 'axios'
import { clearMessages, copyToClipboard, hide, setErrorMessage, show, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'account', 'addressContainer', 'address', 'image', 'copyButtonText', 'generateNewAddressButton'
    ]
  }

  copyAddressToClipboard () {
    copyToClipboard(this.addressTarget.textContent)
    showSuccessNotification('Copied to clipboard')
  }

  generateNewAddress () {
    hide(this.addressContainerTarget)
    this.generateNewAddressButtonTarget.textContent = 'Generating...'
    this.generateNewAddressButtonTarget.setAttribute('disabled', 'disabled')

    const _this = this
    clearMessages(this)
    axios.get('/generate-address/' + this.accountTarget.value)
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
        _this.generateNewAddressButtonTarget.textContent = 'Generate New Address'
        _this.generateNewAddressButtonTarget.removeAttribute('disabled')
      })
  }
}
