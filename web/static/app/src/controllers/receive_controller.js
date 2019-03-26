import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return [
      'account', 'address', 'image'
    ]
  }

  generate (e) {
    e.preventDefault()
    let _this = this
    axios.get('/generate-address/' + this.accountTarget.value)
      .then((response) => {
        let result = response.data
        if (result.success) {
          _this.addressTarget.textContent = result.address
          _this.imageTarget.setAttribute('src', result.imageData)
        } else {
          window.alert(result.message)
        }
      })
      .catch((error) => {
        console.log(error)
        window.alert('Unable to generate address. Something went wrong.')
      })
  }
}
