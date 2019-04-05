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
  hide(controller.successMessageTarget)
  show(controller.errorMessageTarget)
}

export const setSuccessMessage = (controller, message) => {
  controller.successMessageTarget.innerHTML = message
  hide(controller.errorMessageTarget)
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

export const hide = (el) => {
  el.classList.add('d-none')
}

export const show = (el) => {
  el.classList.remove('d-none')
}
