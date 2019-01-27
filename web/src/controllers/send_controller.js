import { Controller } from 'stimulus'

export default class extends Controller {
  addAnotherAddress () {
    console.log('added')
    window.alert('hi ')
  }
}
