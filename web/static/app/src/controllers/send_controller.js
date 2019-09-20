import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, listenForBalanceUpdate, isHidden } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'form',
      'sourceAccount', 'sourceAccountSpan',
      'spendUnconfirmed',
      'destinations', 'destinationTemplate', 'address', 'addressError', 'amount', 'amountUsd', 'amountError', 'maxSendAmountCheck', 'removeDestinationBtn',
      'destinationAccounts', 'destinationAccountTemplate', 'destinationAccount',
      'useCustom', 'toggleCustomInputPnl', 'fetchingUtxos', 'utxoSelectionProgressBar', 'customInputsTable', 'utxoCheckbox',
      'changeOutputs', 'numberOfChangeOutputs', 'useRandomChangeOutputs', 'generateOutputsButton', 'generatedChangeOutputs',
      'changeOutputTemplate', 'changeOutputPercentage', 'changeOutputAddress', 'changeOutputAmount',
      'errors',
      'nextButton',
      // from wallet passphrase modal (utils.html)
      'walletPassphrase', 'passwordError', 'transactionDetails', 'fee', 'estimateSize', 'exchangeRate', 'balanceAfter'
    ]
  }

  connect () {
    listenForBalanceUpdate(this)
  }

  initialize () {
    this.setBusy(false)
    this.initializeSendToAddress()

    let _this = this

    // bootstrap4-toggle is not triggering stimulusjs change action directly
    this.useCustomTarget.onchange = function () {
      _this.toggleUseCustom()
    }

    this.exchangeRate = parseFloat(this.sourceAccountTarget.getAttribute('data-echange-rate'))
    if (this.exchangeRate === 0) {
      this.exchangeRateTarget.textContent = 'N/A'
      this.amountUsdTargets.forEach(target => {
        hide(target.parentElement)
      })
    }
    this.customInputPnlOpnen = false
  }

  setBusy (busy) {
    this.busy = busy
    this.updateSendButtonState()
  }

  toggleCustomInputPnlClicked () {
    if (this.toggleCustomInputPnlTarget.checked) {
      $('#custom-inputs').slideDown()
    } else {
      $('#custom-inputs').slideUp()
    }
    this.customInputPnlOpnen = !this.customInputPnlOpnen
  }

  toggleSpendUnconfirmed () {
    if (this.useCustomTarget.checked) {
      this.openCustomInputsAndChangeOutputsPanel()
    }
  }

  toggleUseCustom () {
    if (!this.useCustomTarget.checked) {
      hide(this.toggleCustomInputPnlTarget.parentElement)
      this.resetCustomInputsAndChangeOutputs()
      this.updateSendButtonState()
      return
    }

    this.openCustomInputsAndChangeOutputsPanel()
    show(this.toggleCustomInputPnlTarget.parentElement)
    this.calculateCustomInputsPercentage()
    this.updateSendButtonState()
  }

  maxSendAmountCheckboxToggle (event) {
    this.setMaxAmountForDestination(event.currentTarget)
  }

  utxoSelectedOrDeselected () {
    this.calculateCustomInputsPercentage()
    this.updateMaxAmountFieldIfSet()
    this.updateSendButtonState()
  }

  updateSendButtonState () {
    if (this.busy || !this.validateSendForm(true)) {
      this.nextButtonTarget.disabled = true
      this.nextButtonTarget.classList.add('disabledBtn')
    } else {
      this.nextButtonTarget.removeAttribute('disabled')
      this.nextButtonTarget.classList.remove('disabledBtn')
    }
  }

  toggleCustomChangeOutputsVisibility () {
    this.clearMessages()
    this.useCustomChangeOutput = !this.useCustomChangeOutput
    this.resetChangeOutput()
  }

  initializeSendToAddress () {
    if (this.sendingToAddress) {
      return
    }
    this.destinationsTarget.innerHTML = ''
    this.destinationAccountsTarget.innerHTML = ''
    this.destinationCount = 0
    this.destinationIndex = 0
    this.newDestination()

    this.sendingToAddress = true
    this.sendingToAccount = false

    this.maxSendAmountCheckTargets.forEach(checkbox => {
      checkbox.removeAttribute('readonly')
      checkbox.parentElement.classList.remove('disabled')
    })
    this.updateSendButtonState()
  }

  initializeSendToAccount () {
    if (this.sendingToAccount) {
      return
    }
    this.destinationsTarget.innerHTML = ''
    this.destinationAccountsTarget.innerHTML = ''
    this.destinationCount = 0
    this.destinationIndex = 0
    this.newDestinationAccount()

    this.sendingToAccount = true
    this.sendingToAddress = false

    this.maxSendAmountCheckTargets.forEach(checkbox => {
      checkbox.removeAttribute('readonly')
      checkbox.parentElement.classList.remove('disabled')
    })
    this.updateSendButtonState()
  }

  newDestination () {
    if (!this.destinationFieldsValid()) {
      return
    }

    const destinationTemplate = document.importNode(this.destinationTemplateTarget.content, true)

    const destinationNode = destinationTemplate.firstElementChild
    const addressInput = destinationNode.querySelector('input[name="destination-address"]')
    const addressErrorDiv = destinationNode.querySelector('div.address-error')
    const amountInput = destinationNode.querySelector('input.amount')
    const amountUsdInput = destinationNode.querySelector('input.amount-usd')
    const sendMaxCheckbox = destinationNode.querySelector('input[type="checkbox"]')
    const removeDestinationButton = destinationNode.querySelector('button[type="button"].removeDestinationBtn')

    // make clicking on the label toggle the checkbox by setting unique id
    destinationNode.querySelector('.form-check-label').setAttribute('for', `send-max-amount-${this.destinationIndex}`)
    sendMaxCheckbox.setAttribute('id', `send-max-amount-${this.destinationIndex}`)

    // disable checkbox if some other checkbox is currently checked
    if (this.maxSendDestinationIndex >= 0) {
      sendMaxCheckbox.setAttribute('readonly', 'readonly')
      sendMaxCheckbox.parentElement.classList.add('disabled')
    }

    destinationNode.setAttribute('data-index', this.destinationIndex)
    addressInput.setAttribute('data-index', this.destinationIndex)
    addressErrorDiv.setAttribute('data-index', this.destinationIndex)
    amountInput.setAttribute('data-index', this.destinationIndex)
    amountUsdInput.setAttribute('data-index', this.destinationIndex)
    sendMaxCheckbox.setAttribute('data-index', this.destinationIndex)
    removeDestinationButton.setAttribute('data-index', this.destinationIndex)

    this.destinationsTarget.appendChild(destinationTemplate)

    this.destinationIndex++
    this.destinationCount++

    if (this.destinationCount === 1) {
      hide(removeDestinationButton)
    } else {
      this.removeDestinationBtnTargets.forEach(btn => {
        show(btn)
      })
    }

    this.updateSendButtonState()
  }

  newDestinationAccount () {
    if (!this.destinationFieldsValid()) {
      return
    }

    const destinationTemplate = document.importNode(this.destinationAccountTemplateTarget.content, true)

    const destinationNode = destinationTemplate.firstElementChild
    const accountInput = destinationNode.querySelector('select[name="destination-account-number"]')
    const amountInput = destinationNode.querySelector('input.amount')
    const amountUsdInput = destinationNode.querySelector('input.amount-usd')
    const sendMaxCheckbox = destinationNode.querySelector('input[type="checkbox"]')
    const removeDestinationButton = destinationNode.querySelector('button[type="button"].removeDestinationBtn')

    // make clicking on the label toggle the checkbox by setting unique id
    destinationNode.querySelector('.form-check-label').setAttribute('for', `send-max-amount-${this.destinationIndex}`)
    sendMaxCheckbox.setAttribute('id', `send-max-amount-${this.destinationIndex}`)

    // disable checkbox if some other checkbox is currently checked
    if (this.maxSendDestinationIndex >= 0) {
      sendMaxCheckbox.setAttribute('readonly', 'readonly')
      sendMaxCheckbox.parentElement.classList.add('disabled')
    }

    destinationNode.setAttribute('data-index', this.destinationIndex)
    accountInput.setAttribute('data-index', this.destinationIndex)
    amountInput.setAttribute('data-index', this.destinationIndex)
    amountUsdInput.setAttribute('data-index', this.destinationIndex)
    sendMaxCheckbox.setAttribute('data-index', this.destinationIndex)
    removeDestinationButton.setAttribute('data-index', this.destinationIndex)

    this.destinationAccountsTarget.appendChild(destinationTemplate)

    this.destinationIndex++
    this.destinationCount++

    if (this.destinationCount === 1) {
      hide(removeDestinationButton)
    } else {
      this.removeDestinationBtnTargets.forEach(btn => {
        show(btn)
      })
    }

    this.updateSendButtonState()
  }

  destinationAddressEdited (event) {
    const editedAddress = event.currentTarget
    const _this = this

    this.updateSendButtonState()

    axios.post('/validate-address?address=' + editedAddress.value)
      .then((response) => {
        let result = response.data
        if (!result.valid) {
          _this.setDestinationFieldError(editedAddress, result.error ? result.error : 'Invalid address')
          return
        }
        _this.clearDestinationFieldError(editedAddress)
      })
      .catch(() => {
        _this.setDestinationFieldError(editedAddress, 'Cannot validate address. You can continue if you are sure')
      })
  }

  destinationAmountEdited (event) {
    this.updateSendButtonState()

    const amountTarget = event.currentTarget
    const amount = parseFloat(amountTarget.value)
    if (isNaN(amount) || amount <= 0) {
      this.setDestinationFieldError(amountTarget, 'Amount must be a non-zero positive number.', false)
      return
    }
    const accountBalance = this.getAccountBalance()
    const totalSendAmount = this.getTotalSendAmount()
    if (totalSendAmount > accountBalance) {
      this.setDestinationFieldError(amountTarget, `Amount exceeds balance. Please enter ${accountBalance - (totalSendAmount - amount)} or less.`, false)
      return
    }

    this.clearDestinationFieldError(amountTarget)
    // update max send amount field if some other amount field has been updated
    const editedAmountFieldIndex = event.target.getAttribute('data-index')
    if (this.maxSendDestinationIndex !== editedAmountFieldIndex) {
      this.updateMaxAmountFieldIfSet()
    }

    this.setUsdField(amountTarget)

    this.calculateCustomInputsPercentage()
    if (this.useRandomChangeOutputsTarget.checked) {
      this.generateChangeOutputs()
    }
  }

  destinationAmountUsdEdited (event) {
    const amountTarget = event.currentTarget
    const amount = parseFloat(amountTarget.value)
    if (isNaN(amount) || amount <= 0) {
      this.setDestinationFieldError(amountTarget, 'Amount must be a non-zero positive number.', false)
      return
    }

    let dcrAmount = parseFloat(amountTarget.value) / this.exchangeRate

    let dcrAmountTarget
    this.amountTargets.forEach(target => {
      if (target.getAttribute('data-index') === amountTarget.getAttribute('data-index')) {
        dcrAmountTarget = target
      }
    })

    const accountBalance = this.getAccountBalance()
    const totalSendAmount = (this.getTotalSendAmount() - parseFloat(dcrAmountTarget.value))

    if (totalSendAmount + dcrAmount > accountBalance) {
      let amountLeft = this.exchangeRate * (accountBalance - totalSendAmount)
      this.setDestinationFieldError(amountTarget, `Amount exceeds balance. Please enter ${amountLeft.toFixed(4)} or less.`, false)
      this.updateSendButtonState()
      return
    }

    this.clearDestinationFieldError(amountTarget)
    // update max send amount field if some other amount field has been updated
    const editedAmountFieldIndex = event.target.getAttribute('data-index')
    if (this.maxSendDestinationIndex !== editedAmountFieldIndex) {
      this.updateMaxAmountFieldIfSet()
    }

    this.setDcrField(amountTarget)

    this.calculateCustomInputsPercentage()
    if (this.useRandomChangeOutputsTarget.checked) {
      this.generateChangeOutputs()
    }

    this.updateSendButtonState()
  }

  updateMaxAmountFieldIfSet () {
    if (this.maxSendDestinationIndex >= 0) {
      const activeSendMaxCheckbox = document.getElementById(`send-max-amount-${this.maxSendDestinationIndex}`)
      this.setMaxAmountForDestination(activeSendMaxCheckbox)
    }
  }

  setMaxAmountForDestination (sendMaxCheckbox) {
    if (sendMaxCheckbox.hasAttribute('readonly')) {
      sendMaxCheckbox.checked = false
      return
    }

    if (!sendMaxCheckbox.checked && this.useCustomTarget.checked) {
      show(this.changeOutputsTarget)
    }

    const index = parseInt(sendMaxCheckbox.getAttribute('data-index'))
    const destinationNode = document.querySelector(`div.destination[data-index="${index}"]`)
    const amountField = destinationNode.querySelector('input.amount')
    const amountUsdField = destinationNode.querySelector('input.amount-usd')

    this.maxSendDestinationIndex = index
    const currentAmount = amountField.value
    amountField.setAttribute('readonly', 'readonly')
    amountUsdField.setAttribute('readonly', 'readonly')
    this.maxSendAmountCheckTargets.forEach(checkbox => {
      checkbox.setAttribute('readonly', 'readonly')
      checkbox.parentElement.classList.add('disabled')
    })

    const uncheckCurrentMaxCheckbox = () => {
      sendMaxCheckbox.checked = false
      this.maxSendDestinationIndex = -1
      amountField.value = currentAmount
      amountField.removeAttribute('readonly')
      amountUsdField.removeAttribute('readonly')
      this.maxSendAmountCheckTargets.forEach(checkbox => {
        checkbox.removeAttribute('readonly')
        checkbox.parentElement.classList.remove('disabled')
      })
    }

    if (!sendMaxCheckbox.checked) {
      uncheckCurrentMaxCheckbox()
      return
    }

    // temporarily set the destination amount field to 1 to make destination validation pass
    // value will be reset afterwards if there are other destination validation errors or if getting max amount fails
    amountField.value = 1
    amountUsdField.value = 1
    amountField.classList.remove('is-invalid')
    amountUsdField.classList.remove('is-invalid')
    if (!this.destinationFieldsValid()) {
      uncheckCurrentMaxCheckbox()
      return
    }

    amountField.value = ''
    amountUsdField.value = ''
    let _this = this
    this.getMaxSendAmount(amount => {
      amountField.value = amount
      this.setUsdField(amountField)
      _this.clearDestinationFieldError(amountField)
      sendMaxCheckbox.removeAttribute('readonly')
      sendMaxCheckbox.parentElement.classList.remove('disabled')
      _this.hideChangeOutputPanel()
      _this.updateSendButtonState()
      _this.calculateCustomInputsPercentage()
    }, (errMsg) => {
      uncheckCurrentMaxCheckbox()
      _this.setDestinationFieldError(amountField, errMsg, false)
    })

    this.setUsdField(amountField)
  }

  getMaxSendAmount (successCallback, errorCallback) {
    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`
    if (this.spendUnconfirmedTarget.checked) {
      queryParams += '&spend-unconfirmed=true'
    }

    this.setBusy(true)
    let _this = this
    axios.get('/max-send-amount?' + queryParams)
      .then((response) => {
        let result = response.data
        if (!result.error) {
          successCallback(result.amount)
        } else if (errorCallback) {
          errorCallback(result.error)
        } else {
          _this.setErrorMessage(result.error)
        }
      })
      .catch(() => {
        if (errorCallback) {
          errorCallback('A server error occurred')
        } else {
          _this.setErrorMessage('A server error occurred')
        }
      })
      .then(() => {
        _this.setBusy(false)
      })
  }

  setUsdField (dcrFieldTarget) {
    const _this = this
    if (this.exchangeRate > 0) {
      let usdAmount = parseFloat(dcrFieldTarget.value) * this.exchangeRate
      if (isNaN(usdAmount)) {
        usdAmount = 0
      }
      this.amountUsdTargets.forEach(target => {
        if (target.getAttribute('data-index') === dcrFieldTarget.getAttribute('data-index')) {
          target.value = usdAmount.toFixed(2)
          _this.clearDestinationFieldError(target)
        }
      })
    }
  }

  setDcrField (usdFieldTarget) {
    const _this = this
    if (this.exchangeRate > 0) {
      let dcrAmount = parseFloat(usdFieldTarget.value) / this.exchangeRate
      this.amountTargets.forEach(target => {
        if (target.getAttribute('data-index') === usdFieldTarget.getAttribute('data-index')) {
          target.value = dcrAmount.toFixed(2)
          _this.clearDestinationFieldError(target)
        }
      })
    }
  }

  destinationFieldsValid (dontModifyErrorOutput) {
    if (!dontModifyErrorOutput) {
      this.clearMessages()
    }
    let fieldsAreValid = true

    for (const addressTarget of this.addressTargets) {
      if (addressTarget.value === '') {
        if (!dontModifyErrorOutput) {
          this.setDestinationFieldError(addressTarget, 'Destination address should not be empty', false)
        }
        fieldsAreValid = false
      } else if (addressTarget.classList.contains('is-invalid')) {
        fieldsAreValid = false
      } else {
        if (!dontModifyErrorOutput) {
          this.clearDestinationFieldError(addressTarget)
        }
      }
    }

    for (const amountTarget of this.amountTargets) {
      const amount = parseFloat(amountTarget.value)
      if (isNaN(amount) || amount <= 0) {
        if (!dontModifyErrorOutput) {
          this.setDestinationFieldError(amountTarget, 'Amount must be a non-zero positive number', false)
        }
        fieldsAreValid = false
      } else if (amountTarget.classList.contains('is-invalid')) {
        fieldsAreValid = false
      }
    }

    for (const amountTarget of this.amountUsdTargets) {
      const amount = parseFloat(amountTarget.value)
      if (isNaN(amount) || amount <= 0) {
        if (!dontModifyErrorOutput) {
          this.setDestinationFieldError(amountTarget, 'Amount must be a non-zero positive number', false)
        }
        fieldsAreValid = false
      } else if (amountTarget.classList.contains('is-invalid')) {
        fieldsAreValid = false
      }
    }

    return fieldsAreValid
  }

  setDestinationFieldError (element, errorMessage, append) {
    const errorElement = element.parentElement.lastElementChild
    if (append && errorElement.innerText !== '') {
      errorElement.innerText += `, ${errorMessage.toLowerCase()}`
    } else {
      errorElement.innerText = errorMessage
    }
    element.classList.add('is-invalid')
    show(errorElement)
    element.select()
    element.focus()

    this.updateSendButtonState()

    this.alignDestinationField(element)
  }

  clearDestinationFieldError (element) {
    const errorElement = element.parentElement.lastElementChild
    errorElement.innerText = ''
    hide(errorElement)
    element.classList.remove('is-invalid')

    this.updateSendButtonState()

    this.alignDestinationField(element)
  }

  alignDestinationField (element) {
    const index = element.getAttribute('data-index')
    let addressTarget, amountTarget, amountUsdTarget, sendMaxTarget, accountTarget
    this.addressTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        addressTarget = el
      }
    })
    this.amountTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        amountTarget = el
      }
    })
    this.amountUsdTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        amountUsdTarget = el
      }
    })
    this.maxSendAmountCheckTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        sendMaxTarget = el
      }
    })
    this.destinationAccountTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        accountTarget = el
      }
    })

    const amountErr = amountTarget.parentElement.lastElementChild.innerHTML
    const amountUsdErr = amountUsdTarget.parentElement.lastElementChild.innerHTML
    let addressErr = ''
    if (this.sendingToAddress) {
      addressErr = addressTarget.parentElement.lastElementChild.innerHTML
    }

    sendMaxTarget.parentElement.parentElement.style.marginBottom = '0px'
    amountTarget.parentElement.style.marginBottom = '0px'
    amountUsdTarget.parentElement.style.marginBottom = '0px'
    if (this.sendingToAddress) {
      addressTarget.parentElement.style.marginBottom = '0px'
    } else {
      accountTarget.parentElement.style.marginBottom = '0px'
    }

    if (amountErr === '' && amountUsdErr === '' && addressErr === '') {
      return
    }

    // this is the number of char that is shown per line in the error div associated with the element
    const addressErrCharPerLine = 54
    const amountErrCharPerLine = 27

    let numberOfAddressErrLines = addressErr.length / addressErrCharPerLine
    numberOfAddressErrLines = Math.round(numberOfAddressErrLines) < numberOfAddressErrLines ? Math.round(numberOfAddressErrLines) + 1 : Math.round(numberOfAddressErrLines)
    let numberOfAmountErrLines = amountErr.length / amountErrCharPerLine
    numberOfAmountErrLines = Math.round(numberOfAmountErrLines) < numberOfAmountErrLines ? Math.round(numberOfAmountErrLines) + 1 : Math.round(numberOfAmountErrLines)

    let numberOfAmountUsdErrLines = amountUsdErr.length / amountErrCharPerLine
    numberOfAmountUsdErrLines = Math.round(numberOfAmountUsdErrLines) < numberOfAmountUsdErrLines ? Math.round(numberOfAmountUsdErrLines) + 1 : Math.round(numberOfAmountUsdErrLines)

    let maxLines = numberOfAmountErrLines
    if (numberOfAddressErrLines > maxLines) {
      maxLines = numberOfAddressErrLines
    }

    if (numberOfAddressErrLines > maxLines) {
      maxLines = numberOfAddressErrLines
    }

    const pixelPerLine = 20
    sendMaxTarget.parentElement.parentElement.style.marginBottom = `${pixelPerLine * maxLines}px`
    if (this.sendingToAddress) {
      addressTarget.parentElement.style.marginBottom = `${pixelPerLine * (maxLines - numberOfAddressErrLines)}px`
    } else {
      accountTarget.parentElement.style.marginBottom = `${pixelPerLine * (maxLines - numberOfAddressErrLines)}px`
    }
    amountTarget.parentElement.style.marginBottom = `${pixelPerLine * (maxLines - numberOfAmountErrLines)}px`
    amountUsdTarget.parentElement.style.marginBottom = `${pixelPerLine * (maxLines - numberOfAmountUsdErrLines)}px`
  }

  removeDestination (event) {
    if (this.destinationCount === 1) {
      return
    }

    const targetElement = event.currentTarget
    const index = parseInt(targetElement.getAttribute('data-index'))

    if (this.sendingToAddress) {
      this.destinationsTarget.removeChild(this.destinationsTarget.querySelector(`div.destination[data-index="${index}"]`))
    } else {
      this.destinationAccountsTarget.removeChild(this.destinationAccountsTarget.querySelector(`div.destination[data-index="${index}"]`))
    }

    this.destinationCount--

    if (this.destinationCount === 1) {
      hide(this.removeDestinationBtnTarget)
    }

    if (this.maxSendDestinationIndex === index) {
      this.maxSendAmountCheckTargets.forEach(checkbox => {
        checkbox.removeAttribute('readonly')
        checkbox.parentElement.classList.remove('disabled')
      })
    }

    this.updateMaxAmountFieldIfSet()
    this.updateSendButtonState()
  }

  resetCustomInputsAndChangeOutputs () {
    show(this.fetchingUtxosTarget)

    $('#custom-inputs').slideUp()
    this.customInputsTableTarget.innerHTML = ''

    this.hideChangeOutputPanel()
  }

  openCustomInputsAndChangeOutputsPanel () {
    this.resetCustomInputsAndChangeOutputs()
    if (!this.customInputPnlOpnen) {
      $('#custom-inputs').slideDown()
    }

    const _this = this
    const fetchUtxoSuccess = unspentOutputs => {
      const utxos = unspentOutputs.map(utxo => {
        let receiveDateTime = new Date(utxo.receive_time * 1000).toString().split(' ').slice(0, 5).join(' ')
        let dcrAmount = utxo.amount / 100000000
        return `<tr>
                  <td width='5%'>
                    <input data-target='send.utxoCheckbox' data-action='click->send#utxoSelectedOrDeselected' type='checkbox' class='custom-input' 
                    name='utxo' value='${utxo.key}' data-amount='${dcrAmount}' data-address='${utxo.address}' />
                  </td>
                  <td width='40%'>${utxo.address}</td>
                  <td width='15%'>${dcrAmount} DCR</td>
                  <td width='25%'>${receiveDateTime}</td>
                  <td width='15%'>${utxo.confirmations} confirmation(s)</td>
                </tr>`
      })

      _this.customInputsTableTarget.innerHTML = utxos.join('')
      hide(this.fetchingUtxosTarget)
      if (!(this.maxSendDestinationIndex >= 0)) {
        show(_this.changeOutputsTarget)
      }
    }

    const accountNumber = this.sourceAccountTarget.value
    this.getUnspentOutputs(accountNumber, fetchUtxoSuccess)
  }

  getUnspentOutputs (accountNumber, successCallback) {
    this.nextButtonTarget.innerHTML = 'Loading...'
    this.setBusy(true)
    this.updateSendButtonState()

    let url = `/unspent-outputs/${accountNumber}`
    if (this.spendUnconfirmedTarget.checked) {
      url += '?spend-unconfirmed=true'
    }

    let _this = this
    axios.get(url).then(function (response) {
      let result = response.data
      if (result.success) {
        successCallback(result.message)
      } else {
        _this.setErrorMessage(result.message)
      }
    }).catch(function () {
      _this.setErrorMessage('A server error occurred')
    }).then(function () {
      _this.nextButtonTarget.innerHTML = 'Send'
      _this.setBusy(false)
    })
  }

  // triggered when destination amount fields are edited or when utxo is selected
  calculateCustomInputsPercentage () {
    if (!this.useCustomTarget.checked) {
      return
    }

    const sendAmount = this.getTotalSendAmount()
    const selectedInputSum = this.getSelectedInputsSum()

    let percentage = 0
    if (selectedInputSum >= sendAmount) {
      percentage = 100
    } else {
      percentage = (selectedInputSum / sendAmount) * 100
    }

    this.utxoSelectionProgressBarTarget.style.width = `${percentage}%`
  }

  getAccountBalance () {
    let target = this.sourceAccountTarget
    if (this.sourceAccountTarget.options) {
      target = this.sourceAccountTarget.options[this.sourceAccountTarget.selectedIndex]
    }
    let amount = parseFloat(target.getAttribute('data-spendable-balance'))
    let unconfirmed = parseFloat(target.getAttribute('data-unconfirmed-balance'))

    return this.spendUnconfirmedTarget.checked ? amount + unconfirmed : amount
  }

  getTotalAccountBalance () {
    let target = this.sourceAccountTarget
    if (this.sourceAccountTarget.options) {
      target = this.sourceAccountTarget.options[this.sourceAccountTarget.selectedIndex]
    }
    return parseFloat(target.getAttribute('data-total-balance'))
  }

  getTotalSendAmount () {
    let amount = 0
    this.amountTargets.forEach(amountTarget => {
      amount += parseFloat(amountTarget.value)
    })
    return amount > 0 ? amount : 0
  }

  getSelectedInputsSum () {
    let sum = 0
    this.customInputsTableTarget.querySelectorAll('input.custom-input:checked').forEach(selectedInputElement => {
      sum += parseFloat(selectedInputElement.dataset.amount)
    })
    return sum
  }

  hideChangeOutputPanel () {
    this.resetChangeOutput()
    hide(this.changeOutputsTarget)
    hide(this.generatedChangeOutputsTarget)
  }

  resetChangeOutput () {
    this.useRandomChangeOutputsTarget.checked = false
    this.numberOfChangeOutputsTarget.value = ''
    this.generatedChangeOutputsTarget.innerHTML = ''
  }

  generateChangeOutputs () {
    if (this.generatingChangeOutputs || !this.usingCustomChangeOutput()) {
      return
    }

    this.clearMessages()

    const numberOfChangeOutput = parseFloat(this.numberOfChangeOutputsTarget.value)
    if (isNaN(numberOfChangeOutput) || numberOfChangeOutput < 1) {
      this.showError('Number of change outputs must be 1 or more')
      return
    }

    if (!this.validateSendForm()) {
      return
    }

    let _this = this
    this.getRandomChangeOutputs(numberOfChangeOutput, function (changeOutputs) {
      if (!_this.usingCustomChangeOutput()) {
        return
      }

      // first calculate total change amount to use below in calculating percentages
      _this.totalChangeAmount = 0
      changeOutputs.forEach((changeOutput) => {
        _this.totalChangeAmount += changeOutput.Amount
      })

      changeOutputs.forEach((changeOutput, i) => {
        let template = document.importNode(_this.changeOutputTemplateTarget.content, true)
        const addressElement = template.querySelector('input[name="change-output-address"]')
        const percentageElement = template.querySelector('input[name="change-output-amount-percentage"]')
        const amountElement = template.querySelector('input[name="change-output-amount"]')

        let percentage = 0
        if (_this.useRandomChangeOutputsTarget.checked) {
          amountElement.value = changeOutput.Amount
          percentageElement.setAttribute('disabled', 'disabled')
          percentage = (changeOutput.Amount / _this.totalChangeAmount) * 100
        }
        percentageElement.value = percentage
        addressElement.value = changeOutput.Address

        addressElement.setAttribute('data-index', i)
        percentageElement.setAttribute('data-index', i)
        amountElement.setAttribute('data-index', i)

        _this.generatedChangeOutputsTarget.appendChild(template)
      })

      show(_this.generatedChangeOutputsTarget)
    })
  }

  getRandomChangeOutputs (numberOfOutputs, successCallback, errorCallback) {
    this.generatingChangeOutputs = true
    this.generatedChangeOutputsTarget.innerHTML = ''
    hide(this.generatedChangeOutputsTarget)
    this.generateOutputsButtonTarget.setAttribute('disabled', 'disabled')
    this.generateOutputsButtonTarget.innerHTML = 'Loading...'
    this.numberOfChangeOutputsTarget.setAttribute('disabled', 'disabled')

    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`

    queryParams += `&nChangeOutput=${numberOfOutputs}`

    this.setBusy(true)
    let _this = this
    axios.get('/random-change-outputs?' + queryParams)
      .then((response) => {
        let result = response.data
        if (!result.error) {
          successCallback(result.message)
        } else if (errorCallback) {
          errorCallback(result.error)
        } else {
          _this.setErrorMessage(result.error)
        }
      })
      .catch(() => {
        _this.setErrorMessage('A server error occurred')
      })
      .then(() => {
        _this.generateOutputsButtonTarget.removeAttribute('disabled')
        _this.generateOutputsButtonTarget.innerHTML = 'Generate Change Outputs'
        _this.numberOfChangeOutputsTarget.removeAttribute('disabled')
        _this.generatingChangeOutputs = false
        this.setBusy(false)
      })
  }

  usingCustomChangeOutput () {
    return this.useCustomChangeOutput && !isHidden(this.changeOutputsTarget)
  }

  changeOutputAmountPercentageChanged (event) {
    const targetElement = event.currentTarget
    const index = parseInt(targetElement.getAttribute('data-index'))
    let currentPercentage = parseFloat(targetElement.value)

    let totalPercentageAllotted = 0
    this.changeOutputPercentageTargets.forEach(percentageTarget => {
      totalPercentageAllotted += parseFloat(percentageTarget.value)
    })

    if (totalPercentageAllotted > 100) {
      const previouslyAllotted = totalPercentageAllotted - currentPercentage
      const availablePercentage = 100 - previouslyAllotted
      targetElement.value = availablePercentage
      currentPercentage = availablePercentage
    }

    const totalChangeAmount = this.totalChangeAmount
    this.changeOutputAmountTargets.forEach(function (amountTarget) {
      if (parseInt(amountTarget.getAttribute('data-index')) === index) {
        amountTarget.value = totalChangeAmount * currentPercentage / 100
      }
    })
  }

  getWalletPassphraseAndSubmit () {
    const _this = this
    this.clearMessages()
    if (!this.validateSendForm() || !this.validateChangeOutputAmount()) {
      return
    }

    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`
    axios.get('/tx-fee-and-size?' + queryParams).then(response => {
      const result = response.data
      if (result.error) {
        this.passwordErrorTarget.innerHTML = `<div class="error">${result.error}</div>`
        return
      }
      _this.feeTarget.textContent = `${result.fee} DCR`
      _this.estimateSizeTarget.textContent = `${result.size} bytes`
      _this.balanceAfterTarget.textContent = `${_this.getTotalAccountBalance() - (_this.getTotalSendAmount() + result.fee)} DCR`
    })

    this.transactionDetailsTarget.innerHTML = this.summaryHTML()

    $('#passphrase-modal').modal()
  }

  summaryHTML () {
    const _this = this
    let summaryHTML
    if (this.useCustomTarget.checked) {
      summaryHTML = '<p>You are about to spend these inputs:</p>'
      let inputs = ''
      this.utxoCheckboxTargets.forEach(utxoCheckbox => {
        if (!utxoCheckbox.checked) {
          return
        }
        let utxo = `${utxoCheckbox.value.substring(0, 15)}...${utxoCheckbox.value.substring(utxoCheckbox.value.length - 6)}`
        const amount = parseFloat(utxoCheckbox.getAttribute('data-amount'))
        let usdAmountStr = ''
        if (_this.exchangeRate > 0) {
          usdAmountStr = `($${(amount * _this.exchangeRate).toFixed(2)}) `
        }
        inputs += `<li>${amount} DCR ${usdAmountStr}from ${utxo}`
      })
      summaryHTML += `<ul>${inputs}</ul> <p>Sending them to:</p>`
    } else {
      summaryHTML = '<p>You about to send</p>'
    }
    let destinations = ''
    if (this.sendingToAddress) {
      this.addressTargets.forEach(addressTarget => {
        const index = addressTarget.getAttribute('data-index')
        let currentAmountTarget
        _this.amountTargets.forEach(function (target) {
          if (target.getAttribute('data-index') === index) {
            currentAmountTarget = target
          }
        })
        if (!currentAmountTarget) {
          return
        }
        let amount = parseFloat(currentAmountTarget.value)
        let usdAmountStr = ''
        if (_this.exchangeRate > 0) {
          usdAmountStr = `($${(amount * _this.exchangeRate).toFixed(2)}) `
        }

        destinations += `<li>${amount} DCR ${usdAmountStr}to ${addressTarget.value}</li>`
      })
    } else {
      this.destinationAccountTargets.forEach(accountTarget => {
        const index = accountTarget.getAttribute('data-index')
        let currentAmountTarget
        _this.amountTargets.forEach(function (target) {
          if (target.getAttribute('data-index') === index) {
            currentAmountTarget = target
          }
        })
        if (!currentAmountTarget) {
          return
        }
        let amount = parseFloat(currentAmountTarget.value)
        let usdAmountStr = ''
        if (_this.exchangeRate > 0) {
          usdAmountStr = `($${(amount * _this.exchangeRate).toFixed(2)}) `
        }
        const accountName = accountTarget.options[accountTarget.selectedIndex].getAttribute('data-account-name')
        destinations += `<li>${parseFloat(currentAmountTarget.value)} DCR ${usdAmountStr}to <b>${accountName}</b></li>`
      })
    }

    let changeOutputs = ''
    this.changeOutputAddressTargets.forEach(changeOutputAddressTarget => {
      const index = changeOutputAddressTarget.getAttribute('data-index')
      let currentAmountTarget
      _this.changeOutputAmountTargets.forEach(function (target) {
        if (target.getAttribute('data-index') === index) {
          currentAmountTarget = target
        }
      })
      if (!currentAmountTarget) {
        return
      }
      let amount = parseFloat(currentAmountTarget.value)
      let usdAmountStr = ''
      if (_this.exchangeRate > 0) {
        usdAmountStr = `($${(amount * _this.exchangeRate).toFixed(2)}) `
      }
      changeOutputs += `<li>${(amount.toFixed(2))} DCR ${usdAmountStr}to ${changeOutputAddressTarget.value} (change)</li>`
    })

    summaryHTML += `<ul>${destinations}</ul>`
    if (changeOutputs !== '') {
      summaryHTML += `<ul>${changeOutputs}</ul>`
    }
    return summaryHTML
  }

  validateSendForm (dontModifyErrorOutput) {
    if (!dontModifyErrorOutput) {
      this.errorsTarget.innerHTML = ''
      hide(this.errorsTarget)
    }
    let valid = this.destinationFieldsValid(dontModifyErrorOutput)

    if (this.sourceAccountTarget.value === '') {
      if (!dontModifyErrorOutput) {
        this.showError('The source account is required')
      }
      valid = false
    }

    if (this.useCustomTarget.checked && this.getSelectedInputsSum() < this.getTotalSendAmount()) {
      if (!dontModifyErrorOutput) {
        this.showError('The sum of selected inputs is less than send amount')
      }
      valid = false
    }

    return valid
  }

  validateChangeOutputAmount () {
    this.clearMessages()

    // if the user is using random change output, the server will get accurate values
    if (this.useRandomChangeOutputsTarget.checked) {
      return true
    }

    let totalPercentageAllotted = 0
    this.changeOutputPercentageTargets.forEach(percentageTarget => {
      const thisPercent = parseFloat(percentageTarget.value)
      if (isNaN(thisPercent) || thisPercent <= 0) {
        this.showError('Change amount percentage must be greater than 0')
        return false
      }

      totalPercentageAllotted += thisPercent
    })

    if (this.usingCustomChangeOutput() && totalPercentageAllotted !== 100) {
      this.showError(`Total change amount percentage must be equal to 100. Current total is ${totalPercentageAllotted}`)
      return false
    }

    return true
  }

  submitForm () {
    if (!this.validatePassphrase()) {
      return
    }
    this.setBusy(true)

    $('#passphrase-modal').modal('hide')

    this.nextButtonTarget.innerHTML = 'Sending...'
    this.nextButtonTarget.setAttribute('disabled', 'disabled')

    let postData = $('#send-form').serialize()
    postData += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`

    // clear password input
    this.walletPassphraseTarget.value = ''

    let _this = this
    axios.post('/send', postData).then((response) => {
      let result = response.data
      if (result.error !== undefined) {
        _this.setErrorMessage(result.error)
      } else {
        _this.resetSendForm()
        let txHash = `The transaction was published successfully. Hash: <strong>${result.txHash}</strong>`
        _this.setSuccessMessage(txHash)
      }
    }).catch((e) => {
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      _this.nextButtonTarget.innerHTML = 'Send'
      _this.setBusy(false)
    })
  }

  validatePassphrase () {
    if (this.walletPassphraseTarget.value === '') {
      this.passwordErrorTarget.innerHTML = '<div class="error">Your wallet passphrase is required</div>'
      return false
    }

    return true
  }

  clearFields () {
    this.resetSendForm()
    this.clearMessages()
  }

  resetSendForm () {
    this.resetCustomInputsAndChangeOutputs()

    this.destinationIndex = 0
    this.destinationCount = 0

    if (this.sendingToAddress) {
      this.destinationsTarget.innerHTML = ''
      this.newDestination()
    } else {
      this.destinationAccountsTarget.innerHTML = ''
      this.newDestinationAccount()
    }

    this.spendUnconfirmedTarget.checked = false
    $(this.useCustomTarget).bootstrapToggle('off')
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    hide(this.successMessageTarget)
    show(this.errorMessageTarget)
  }

  setSuccessMessage (message) {
    this.successMessageTarget.innerHTML = message
    hide(this.errorMessageTarget)
    show(this.successMessageTarget)
  }

  clearMessages () {
    hide(this.errorMessageTarget)
    hide(this.successMessageTarget)
    this.errorsTarget.innerHTML = ''
    hide(this.errorsTarget)
  }

  showError (error) {
    this.errorsTarget.innerHTML += `<div class="error">${error}</div>`
    show(this.errorsTarget)
  }
}
