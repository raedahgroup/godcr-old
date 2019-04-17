import toastr from 'toastr'
import ws from './services/messagesocket_service'

export function listenForBalanceUpdate (_this) {
  ws.registerEvtHandler('updateBalance', function (data) {
    if (_this.sourceAccountTarget.options) {
      data.accounts.forEach(account => {
        for (let i = 0; i < _this.sourceAccountTarget.length; i++) {
          const opt = _this.sourceAccountTarget.options[i]
          if (parseInt(opt.value) === account.number) {
            opt.textContent = account.info
          }
        }
      })
    } else {
      if (_this.sourceAccountSpanTarget) {
        _this.sourceAccountSpanTarget.textContent = data.total
      }
    }
  })
}

export const copyToClipboard = str => {
  const el = document.createElement('textarea')
  el.value = str
  el.setAttribute('readonly', '')
  el.style.position = 'absolute'
  el.style.left = '-9999px'
  document.body.appendChild(el)
  const selected =
      document.getSelection().rangeCount > 0
        ? document.getSelection().getRangeAt(0)
        : false
  el.select()
  document.execCommand('copy')
  document.body.removeChild(el)
  if (selected) {
    document.getSelection().removeAllRanges()
    document.getSelection().addRange(selected)
  }
}

export const setErrorMessage = (controller, message) => {
  controller.errorMessageTarget.innerHTML = message
  if (controller.successMessageTarget) {
    hide(controller.successMessageTarget)
  }
  show(controller.errorMessageTarget)
}

export const setSuccessMessage = (controller, message) => {
  controller.successMessageTarget.innerHTML = message
  if (controller.errorMessageTarget) {
    hide(controller.errorMessageTarget)
  }
  show(controller.successMessageTarget)
}

export const clearMessages = (controller) => {
  if (controller.errorMessageTarget) {
    hide(controller.errorMessageTarget)
    controller.errorMessageTarget.innerHTML = ''
  }
  if (controller.successMessageTarget) {
    hide(controller.successMessageTarget)
    controller.successMessageTarget.innerHTML = ''
  }
}

export const showErrorNotification = (message) => {
  toastr.error(message)
}

export const showSuccessNotification = (message) => {
  toastr.success(message)
}

export const hide = (el) => {
  el.classList.add('d-none')
}

export const show = (el) => {
  el.classList.remove('d-none')
}
