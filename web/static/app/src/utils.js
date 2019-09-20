import toastr from 'toastr'
import ws from './services/messagesocket_service'

export function listenForBalanceUpdate (_this) {
  ws.registerEvtHandler('updateBalance', function (data) {
    if (_this.sourceAccountTarget.options) {
      data.accounts.forEach(account => {
        for (let i = 0; i < _this.sourceAccountTarget.options.length; i++) {
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

export const isHidden = (el) => {
  return el.classList.contains('d-none') || el.classList.contains('d-hide')
}

export const truncate = (input, maxLength) => {
  if (input.length <= maxLength) {
    return input
  }
  return input.substring(0, maxLength - 1)
}

export const splitAmountIntoParts = (amountStr) => {
  if (amountStr.indexOf('.') === -1) {
    const splitBalance = amountStr.split(' ')
    return [
      splitBalance[0], '', splitBalance[1]
    ]
  }
  let balanceParts = ['', '', '']
  const splitBalance = amountStr.split('.')
  balanceParts[0] = splitBalance[0]
  balanceParts[1] = splitBalance[1].substring(0, 1)
  const decimalPart = splitBalance[1].split(' ')[0]
  if (decimalPart.length > 2) {
    balanceParts[2] = splitBalance[1].substring(2)
  }
  return balanceParts
}
