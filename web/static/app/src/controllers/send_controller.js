import { Controller } from 'stimulus'
import axios from 'axios'
// customTxRow
export default class extends Controller {
  static get targets () {
    return [
      'errorMessage', 'successMessage',
      'form',
      'sourceAccount',
      'spendUnconfirmed',
      'destinations', 'destinationTemplate', 'address', 'amount', 'removeDestinationButton',
      'useCustom', 'fetchingUtxos', 'utxoSelectionProgressBar', 'customInputsTable',
      'changeOutputs', 'numberOfChangeOutputs', 'useRandomChangeOutputs', 'generateOutputsButton', 'generatedChangeOutputs',
      'changeDestinationTemplate', 'changeOutputPercentage', 'changeOutputAmount',
      'errors',
      'nextButton',
      // from wallet passphrase modal (utils.html)
      'walletPassphrase', 'passwordError'
    ]
  }

  initialize () {
    this.newDestination()
  }

  newDestination () {
    if (!this.destinationFieldsValid()) {
      return
    }

    const destinationTemplate = document.importNode(this.destinationTemplateTarget.content, true)
    this.destinationsTarget.appendChild(destinationTemplate)

    if (this.destinationCount() > 1) {
      this.show(this.removeDestinationButtonTarget)
    }
  }

  destinationFieldsValid () {
    this.clearMessages()
    let fieldsAreValid = true

    for (const addressTarget of this.addressTargets) {
      if (addressTarget.value === '') {
        this.showError('Destination address should not be empty')
        fieldsAreValid = false
        break
      }
    }

    for (const amountTarget of this.amountTargets) {
      const amount = parseFloat(amountTarget.value)
      if (amount < 0) {
        this.showError('Amount must be a non-zero positive number')
        fieldsAreValid = false
        break
      }
    }

    return fieldsAreValid
  }

  destinationCount () {
    return this.destinationsTarget.querySelectorAll('div.destination').length
  }

  removeDestination () {
    if (this.destinationCount() > 1) {
      this.destinationsTarget.removeChild(this.destinationsTarget.querySelector('div.destination:last-child'))
    }
    if (this.destinationCount() <= 1) {
      this.hide(this.removeDestinationButtonTarget)
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

    if (!this.destinationFieldsValid()) {
      this.useCustomTarget.checked = false
      return
    }

    this.openCustomInputsAndChangeOutputsPanel()
  }

  resetCustomInputsAndChangeOutputs () {
    this.show(this.fetchingUtxosTarget)

    $('#custom-inputs').slideUp()
    this.customInputsTableTarget.innerHTML = ''

    this.hide(this.changeOutputsTarget)
    this.useRandomChangeOutputsTarget.checked = false
    this.numberOfChangeOutputsTarget.value = ''
    this.generatedChangeOutputsTarget.innerHTML = ''
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
                    <input data-action='click->send#calculateCustomInputsPercentage' type='checkbox' class='custom-input' name='utxo' value='${utxo.key}' data-amount='${dcrAmount}' />
                  </td>
                  <td width='40%'>${utxo.address}</td>
                  <td width='15%'>${dcrAmount} DCR</td>
                  <td width='25%'>${receiveDateTime}</td>
                  <td width='15%'>${utxo.confirmations} confirmation(s)</td>
                </tr>`
      })

      _this.customInputsTableTarget.innerHTML = utxos.join('')
      _this.hide(this.fetchingUtxosTarget)
      _this.show(_this.changeOutputsTarget)
    }

    const accountNumber = this.sourceAccountTarget.value
    this.getUnspentOutputs(accountNumber, fetchUtxoSuccess)
  }

  getUnspentOutputs (accountNumber, successCallback) {
    this.nextButtonTarget.innerHTML = 'Loading...'
    this.nextButtonTarget.setAttribute('disabled', 'disabled')

    let url = `/unspent-outputs/${accountNumber}`
    if (this.spendUnconfirmedTarget.checked) {
      url += '?getUnconfirmed=true'
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

  toggleCustomChangeOutputsVisibility () {
    this.clearMessages()
    this.useCustomChangeOutput = !this.useCustomChangeOutput
    this.generatedChangeOutputsTarget.innerHTML = ''
    this.numberOfChangeOutputsTarget.value = ''
  }

  generateChangeOutputs () {
    if (this.generatingChangeOutputs || !this.useCustomChangeOutput) {
      return
    }

    const numberOfChangeOutput = parseFloat(this.numberOfChangeOutputsTarget.value)
    if (numberOfChangeOutput < 1) {
      this.showError('Number of change outputs must be 1 or more')
      return
    }

    if (!this.validateSendForm()) {
      return
    }

    let _this = this
    this.getRandomChangeOutputs(numberOfChangeOutput, function (changeOutputDestinations) {
      if (!_this.useCustomChangeOutput) {
        return
      }

      _this.totalChangeAmount = 0

      changeOutputDestinations.forEach((destination, i) => {
        _this.totalChangeAmount += destination.Amount

        let template = document.importNode(_this.changeDestinationTemplateTarget.content, true)
        const addressElement = template.querySelector('input[name="change-output-address"]')
        const percentageElement = template.querySelector('input[name="change-output-amount-percentage"]')
        const amountElement = template.querySelector('input[name="change-output-amount"]')

        let percentage = 0
        if (_this.useRandomChangeOutputsTarget.checked) {
          amountElement.value = destination.Amount
          percentageElement.setAttribute('disabled', 'disabled')
          percentage = (destination.Amount / _this.totalChangeAmount) * 100
        }
        percentageElement.value = percentage
        addressElement.value = destination.Address

        addressElement.setAttribute('data-index', i)
        percentageElement.setAttribute('data-index', i)
        amountElement.setAttribute('data-index', i)

        _this.generatedChangeOutputsTarget.appendChild(template)
      })

      _this.show(_this.generatedChangeOutputsTarget)
    })
  }

  getRandomChangeOutputs (numberOfOutputs, successCallback) {
    this.generatingChangeOutputs = true
    this.generatedChangeOutputsTarget.innerHTML = ''
    this.generateOutputsButtonTarget.setAttribute('disabled', 'disabled')
    this.generateOutputsButtonTarget.innerHTML = 'Loading...'
    this.numberOfChangeOutputsTarget.setAttribute('disabled', 'disabled')

    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`

    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      queryParams += `&source-account=${this.sourceAccountTarget.value}`
    }

    queryParams += `&nChangeOutput=${numberOfOutputs}`

    let _this = this
    axios.get('/random-change-outputs?' + queryParams)
      .then((response) => {
        let result = response.data
        if (result.error !== undefined) {
          _this.setErrorMessage(result.error)
        } else {
          successCallback(result.message)
        }
      })
      .catch((error) => {
        console.log(error)
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
    this.clearMessages()
    if (!this.validateSendForm() || !this.validateChangeOutputAmount()) {
      return
    }
    $('#passphrase-modal').modal()
  }

  validateSendForm () {
    this.errorsTarget.innerHTML = ''
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

    let totalPercentageAllotted = 0
    this.changeOutputPercentageTargets.forEach(percentageTarget => {
      const thisPercent = parseFloat(percentageTarget.value)
      if (thisPercent <= 0) {
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

    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      postData += `&source-account=${this.sourceAccountTarget.value}`
    }

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
    }).catch(() => {
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
    let destinationCount = this.destinationCount()
    while (destinationCount > 1) {
      this.removeDestination()
      destinationCount--
    }
    this.addressTargets.forEach(ele => {
      ele.value = ''
    })
    this.amountTargets.forEach(ele => {
      ele.value = ''
    })
    this.spendUnconfirmedTarget.checked = false
    this.useCustomTarget.checked = false

    this.clearMessages()
  }

  setErrorMessage (message) {
    this.errorMessageTarget.innerHTML = message
    this.hide(this.successMessageTarget)
    this.show(this.errorMessageTarget)
  }

  setSuccessMessage (message) {
    this.successMessageTarget.innerHTML = message
    this.hide(this.errorMessageTarget)
    this.show(this.successMessageTarget)
  }

  clearMessages () {
    this.hide(this.errorMessageTarget)
    this.hide(this.successMessageTarget)
    this.errorsTarget.innerHTML = ''
  }

  hide (el) {
    el.classList.add('d-none')
  }

  show (el) {
    el.classList.remove('d-none')
  }

  showError (error) {
    this.errorsTarget.innerHTML += `<div class="error">${error}</div>`
  }
}
