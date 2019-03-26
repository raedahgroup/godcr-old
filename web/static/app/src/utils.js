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

export const setErrorMessage = (ctl, message) => {
  ctl.errorMessageTarget.innerHTML = message
  hide(ctl.successMessageTarget)
  show(ctl.errorMessageTarget)
}

export const setSuccessMessage = (ctl, message) => {
  ctl.successMessageTarget.innerHTML = message
  hide(ctl.errorMessageTarget)
  show(ctl.successMessageTarget)
}

export const clearMessages = (ctl) => {
  hide(ctl.errorMessageTarget)
  hide(ctl.successMessageTarget)
  ctl.errorsTarget.innerHTML = ''
  ctl.successMessageTarget.innerHTML = ''
}

export const hide = (el) => {
  el.classList.add('d-none')
}

export const show = (el) => {
  el.classList.remove('d-none')
}
