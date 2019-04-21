import { Controller } from 'stimulus'
import axios from 'axios'
import { hide, show, listenForBalanceUpdate } from '../utils'

export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'form',
      'sourceAccount', 'sourceAccountSpan',
      'spendUnconfirmed',
      'destinations', 'destinationTemplate', 'address', 'addressError', 'amount', 'maxSendAmountCheck', 'removeDestinationBtn',
      'useCustom', 'fetchingUtxos', 'utxoSelectionProgressBar', 'customInputsTable', 'utxoCheckbox',
      'changeOutputs', 'numberOfChangeOutputs', 'useRandomChangeOutputs', 'generateOutputsButton', 'generatedChangeOutputs',
      'changeOutputTemplate', 'changeOutputPercentage', 'changeOutputAddress', 'changeOutputAmount',
      'errors',
      'nextButton',
      // from wallet passphrase modal (utils.html)
      'walletPassphrase', 'passwordError', 'transactionDetails'
    ]
  }

  connect () {
    listenForBalanceUpdate(this)
  }

  initialize () {
    this.destinationCount = 0
    this.destinationIndex = 0
    this.newDestination()
    let _this = this

    // bootstrap4-toggle is not triggering stimulusjs change action directly
    this.useCustomTarget.onchange = function () {
      _this.toggleUseCustom()
    }
  }

  toggleSpendUnconfirmed () {
    if (this.useCustomTarget.checked) {
      this.openCustomInputsAndChangeOutputsPanel()
    }
  }

  toggleUseCustom () {
    if (!this.useCustomTarget.checked) {
      this.resetCustomInputsAndChangeOutputs()
      return
    }

    this.openCustomInputsAndChangeOutputsPanel()
  }

  maxSendAmountCheckboxToggle (event) {
    this.setMaxAmountForDestination(event.currentTarget)
  }

  utxoSelectedOrDeselected () {
    this.calculateCustomInputsPercentage()
    this.updateMaxAmountFieldIfSet()
  }

  toggleCustomChangeOutputsVisibility () {
    this.clearMessages()
    this.useCustomChangeOutput = !this.useCustomChangeOutput
    this.resetChangeOutput()
  }

  newDestination () {
    if (!this.destinationFieldsValid()) {
      return
    }

    const destinationTemplate = document.importNode(this.destinationTemplateTarget.content, true)

    const destinationNode = destinationTemplate.firstElementChild
    const addressInput = destinationNode.querySelector('input[name="destination-address"]')
    const addressErrorDiv = destinationNode.querySelector('div.address-error')
    const amountInput = destinationNode.querySelector('input[name="destination-amount"]')
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
  }

  destinationAddressEdited (event) {
    const _this = this
    const index = event.currentTarget.getAttribute('data-index')
    let addressErrorTarget, amountTarget, sendMaxTarget
    this.addressErrorTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        addressErrorTarget = el
      }
    })
    this.amountTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        amountTarget = el
      }
    })
    this.maxSendAmountCheckTargets.forEach(el => {
      if (el.getAttribute('data-index') === index) {
        sendMaxTarget = el
      }
    })

    axios.post('/validate-address?address=' + event.currentTarget.value)
      .then((response) => {
        let result = response.data
        if (!result.valid) {
          addressErrorTarget.textContent = result.error ? result.error : 'Invalid address'
          amountTarget.parentElement.style.marginBottom = '20px'
          sendMaxTarget.parentElement.style.marginBottom = '20px'
          return
        }
        addressErrorTarget.textContent = ''
        amountTarget.parentElement.style.marginBottom = ''
        sendMaxTarget.parentElement.style.marginBottom = ''
      })
      .catch(() => {
        addressErrorTarget.textContent = 'Cannot validate address. You can continue if you are sure'
        amountTarget.parentElement.style.marginBottom = ''
        sendMaxTarget.parentElement.style.marginBottom = ''
      })
  }

  destinationAmountEdited (event) {
    // update max send amount field if some other amount field has been updated
    const editedAmountFieldIndex = event.target.getAttribute('data-index')
    if (this.maxSendDestinationIndex !== editedAmountFieldIndex) {
      this.updateMaxAmountFieldIfSet()
    }

    this.calculateCustomInputsPercentage()
    if (this.useRandomChangeOutputsTarget.checked) {
      this.generateChangeOutputs()
    }
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

    if (!sendMaxCheckbox.checked) {
      show(this.changeOutputsTarget)
    }

    const index = parseInt(sendMaxCheckbox.getAttribute('data-index'))
    const destinationNode = document.querySelector(`div.destination[data-index="${index}"]`)
    const amountField = destinationNode.querySelector('input[name="destination-amount"]')

    this.maxSendDestinationIndex = index
    const currentAmount = amountField.value
    amountField.setAttribute('readonly', 'readonly')
    this.maxSendAmountCheckTargets.forEach(checkbox => {
      checkbox.setAttribute('readonly', 'readonly')
      checkbox.parentElement.classList.add('disabled')
    })

    const uncheckCurrentMaxCheckbox = () => {
      sendMaxCheckbox.checked = false
      this.maxSendDestinationIndex = -1
      amountField.value = currentAmount
      amountField.removeAttribute('readonly')
      this.maxSendAmountCheckTargets.forEach(checkbox => {
        checkbox.removeAttribute('readonly')
        checkbox.parentElement.classList.remove('disabled')
      })
    }

    if (!sendMaxCheckbox.checked) {
      uncheckCurrentMaxCheckbox()
      amountField.value = ''
      return
    }

    // temporarily set the destination amount field to 1 to make destination validation pass
    // value will be reset afterwards if there are other destination validation errors or if getting max amount fails
    amountField.value = 1
    if (!this.destinationFieldsValid()) {
      uncheckCurrentMaxCheckbox()
      return
    }

    amountField.value = ''
    let _this = this
    this.getMaxSendAmount(amount => {
      amountField.value = amount
      sendMaxCheckbox.removeAttribute('readonly')
      sendMaxCheckbox.parentElement.classList.remove('disabled')
      _this.hideChangeOutputPanel()
      _this.calculateCustomInputsPercentage()
    }, (errMsg) => {
      uncheckCurrentMaxCheckbox()
      _this.setDestinationFieldError(amountField, errMsg, false)
    })
  }

  getMaxSendAmount (successCallback, errorCallback) {
    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`
    if (this.spendUnconfirmedTarget.checked) {
      queryParams += '&spend-unconfirmed=true'
    }

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
  }

  destinationFieldsValid () {
    this.clearMessages()
    let fieldsAreValid = true

    for (const addressTarget of this.addressTargets) {
      if (addressTarget.value === '') {
        this.setDestinationFieldError(addressTarget, 'Destination address should not be empty', false)
        fieldsAreValid = false
      } else {
        this.clearDestinationFieldError(addressTarget)
      }
    }

    for (const amountTarget of this.amountTargets) {
      const amount = parseFloat(amountTarget.value)
      if (isNaN(amount) || amount <= 0) {
        this.setDestinationFieldError(amountTarget, 'Amount must be a non-zero positive number', true)
        fieldsAreValid = false
      } else {
        amountTarget.classList.remove('is-invalid')
      }
    }

    return fieldsAreValid
  }

  setDestinationFieldError (element, errorMessage, append) {
    const errorElement = element.parentElement.parentElement.lastElementChild
    if (append && errorElement.innerText !== '') {
      errorElement.innerText += `, ${errorMessage.toLowerCase()}`
    } else {
      errorElement.innerText = errorMessage
    }

    // element.classList.add('is-invalid')
    show(errorElement)
  }

  clearDestinationFieldError (element) {
    const errorElement = element.parentElement.parentElement.lastElementChild
    errorElement.innerText = ''
    hide(errorElement)
    element.classList.remove('is-invalid')
  }

  removeDestination (event) {
    if (this.destinationCount === 1) {
      return
    }

    const targetElement = event.currentTarget
    const index = parseInt(targetElement.getAttribute('data-index'))

    this.destinationsTarget.removeChild(this.destinationsTarget.querySelector(`div.destination[data-index="${index}"]`))
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
  }

  resetCustomInputsAndChangeOutputs () {
    show(this.fetchingUtxosTarget)

    $('#custom-inputs').slideUp()
    this.customInputsTableTarget.innerHTML = ''

    this.hideChangeOutputPanel()
  }

  openCustomInputsAndChangeOutputsPanel () {
    this.resetCustomInputsAndChangeOutputs()
    $('#custom-inputs').slideDown()

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
    this.nextButtonTarget.setAttribute('disabled', 'disabled')

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
      _this.nextButtonTarget.innerHTML = 'Next'
      _this.nextButtonTarget.removeAttribute('disabled')
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

  getTotalSendAmount () {
    let amount = 0
    this.amountTargets.forEach(amountTarget => {
      amount += parseFloat(amountTarget.value)
    })
    return amount
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
    if (this.generatingChangeOutputs || !this.useCustomChangeOutput) {
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
      if (!_this.useCustomChangeOutput) {
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
      })
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
    let summaryHTML
    if (this.useCustomTarget.checked) {
      summaryHTML = '<p>You are about to spend the input(s)</p>'
      let inputs = ''
      this.utxoCheckboxTargets.forEach(utxoCheckbox => {
        if (!utxoCheckbox.checked) {
          return
        }
        inputs += `<li>${parseFloat(utxoCheckbox.getAttribute('data-amount')).toFixed(8)} from ${utxoCheckbox.getAttribute('data-address')}`
      })
      summaryHTML += `<ul>${inputs}</ul> <p>and send</p>`
    } else {
      summaryHTML = '<p>You about to send</p>'
    }
    let destinations = ''
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
      destinations += `<li>${parseFloat(currentAmountTarget.value).toFixed(8)} DCR to ${addressTarget.value}</li>`
    })

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
      destinations += `<li>${parseFloat(currentAmountTarget.value).toFixed(8)} DCR to ${changeOutputAddressTarget.value} (change)</li>`
    })

    summaryHTML += `<ul>${destinations}</ul>`
    this.transactionDetailsTarget.innerHTML = summaryHTML
    $('#passphrase-modal').modal()
  }

  validateSendForm () {
    this.errorsTarget.innerHTML = ''
    hide(this.errorsTarget)
    let valid = this.destinationFieldsValid()

    if (this.sourceAccountTarget.value === '') {
      this.showError('The source account is required')
      valid = false
    }

    if (this.useCustomTarget.checked && this.getSelectedInputsSum() < this.getTotalSendAmount()) {
      this.showError('The sum of selected inputs is less than send amount')
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

    if (this.useCustomChangeOutput && totalPercentageAllotted !== 100) {
      this.showError(`Total change amount percentage must be equal to 100. Current total is ${totalPercentageAllotted}`)
      return false
    }

    return true
  }

  submitForm () {
    if (!this.validatePassphrase()) {
      return
    }

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
      _this.nextButtonTarget.innerHTML = 'Next'
      _this.nextButtonTarget.removeAttribute('disabled')
    })
  }

  validatePassphrase () {
    if (this.walletPassphraseTarget.value === '') {
      this.passwordErrorTarget.innerHTML = '<div class="error">Your wallet passphrase is required</div>'
      return false
    }

    return true
  }

  resetSendForm () {
    this.resetCustomInputsAndChangeOutputs()

    this.destinationsTarget.innerHTML = ''
    this.destinationIndex = 0
    this.destinationCount = 0
    this.newDestination()

    this.addressTargets.forEach(ele => {
      ele.value = ''
    })
    this.amountTargets.forEach(ele => {
      ele.value = ''
    })
    this.spendUnconfirmedTarget.checked = false
    $(this.useCustomTarget).bootstrapToggle('off')

    this.clearMessages()
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
