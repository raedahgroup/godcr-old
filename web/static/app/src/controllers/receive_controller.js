import { Controller } from 'stimulus'
import axios from 'axios'

import { copyToClipboard, setErrorMessage } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'account', 'address', 'image'
    ]
  }

  copyToClipboard () {
    copyToClipboard(this.addressTarget.textContent)
  }

  generate () {
    let _this = this
    _this.clearMessages()
    axios.get('/generate-address/' + this.accountTarget.value)
      .then((response) => {
        let result = response.data
        if (result.success) {
          _this.addressTarget.textContent = result.address
          _this.imageTarget.setAttribute('src', result.imageData)
        } else {
          setErrorMessage(result.message)
        }
      })
      .catch(() => {
        setErrorMessage('Unable to generate address. Something went wrong.')
      })
  }
}
