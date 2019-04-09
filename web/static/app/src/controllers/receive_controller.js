import { Controller } from 'stimulus'
import axios from 'axios'
import { clearMessages, copyToClipboard, setErrorMessage, showSuccessNotification } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'account', 'address', 'image', 'copyButtonText', 'generateNewAddressButton'
    ]
  }

  copyAddressToClipboard () {
    copyToClipboard(this.addressTarget.textContent)
    showSuccessNotification('Copied to clipboard')
  }

  generateNewAddress () {
    this.generateNewAddressButtonTarget.textContent = 'Generating...'
    this.generateNewAddressButtonTarget.setAttribute('disabled', 'disabled')

    let _this = this
    this.generate(function () {
      _this.generateNewAddressButtonTarget.textContent = 'Generate New Address'
      _this.generateNewAddressButtonTarget.removeAttribute('disabled')
    })
  }

  generate (onComplete) {
    let _this = this
    clearMessages(this)
    axios.get('/generate-address/' + this.accountTarget.value)
      .then((response) => {
        let result = response.data
        if (result.success) {
          _this.addressTarget.textContent = result.address
          _this.imageTarget.setAttribute('src', result.imageData)
        } else {
          setErrorMessage(_this, result.message)
        }
      })
      .catch(() => {
        setErrorMessage(_this, 'Unable to generate address. Something went wrong.')
      })
      .then(function () {
        if (onComplete) {
          onComplete()
        }
      })
  }
}
