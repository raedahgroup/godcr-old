import { Controller } from 'stimulus'
import { hide, show } from '../utils'

export default class extends Controller {
  static get targets () {
    return ['syncDetails']
  }

  showDetails (e) {
    hide(e.target)
    show(this.syncDetailsTarget)
  }
}
