import { Controller } from 'stimulus'
import axios from 'axios'

export default class extends Controller {
  static get targets () {
    return ['sourceAccount', 'address', 'amount', 'destinations', 'destinationTemplate', 'changeAddress', 'changeOutputPercentage',
      'changeAmount', 'errors', 'customInput', 'customTxRow', 'customInputContent', 'submitButton', 'nextButton',
      'removeDestinationButton', 'form', 'walletPassphrase', 'passwordError', 'useCustom',
      'spendUnconfirmed', 'errorMessage', 'successMessage',
      'progressBar', 'changeOutputsCard', 'changeOutputPnl', 'numberOfChangeOutputs',
      'useRandomChangeOutputs', 'generateOutputsButton', 'changeOutputContent', 'changeDestinationTemplate', 'changeOutputAddress',
      'changeOutputAmount']
  }

  initialize () {
    this.newDestination()
  }

  validateSendForm () {
    this.errorsTarget.innerHTML = ''
    let errors = []
    let isClean = this.validateDestinationsField()

    if (this.sourceAccountTarget.value === '') {
      errors.push('The source account is required')
    }

    if (this.useCustomTarget.value && this.getSelectedInputsSum() < this.getTotalSendAmount()) {
      errors.push('The sum of selected inputs is less than send amount')
    }

    if (errors.length) {
      for (let i in errors) {
        this.showError(errors[i])
      }
      isClean = false
    }

    return isClean
  }

  validateAndRefreshPercentage () {
    if (this.validateDestinationsField()) {
      this.calculateSelectedInputPercentage()
    }
  }

  validateDestinationsField () {
    this.clearMessages()
    let isClean = true
    this.addressTargets.forEach((el, i) => {
      if (el.value === '') {
        this.showError('The destination address and amount are required')
        isClean = false
      }
    })
    this.amountTargets.forEach((el, i) => {
      let amount = parseFloat(el.value)
      if (!(amount > 0)) {
        this.showError('Amount must be a non-zero positive number')
        isClean = false
      }
    })
    return isClean
  }

  validateChangeOutputAmount () {
    this.clearMessages()
    let isClean = true
    this.changeOutputAmountTargets.forEach(el => {
      let amount = parseFloat(el.value)
      if (!(amount > 0)) {
        this.showError('Change amount must be a non-zero positive number')
        isClean = false
      }
    })
    return isClean
  }

  getSelectedInputsSum () {
    let sum = 0
    this.customTxRowTarget.querySelectorAll('input.custom-input:checked').forEach((el, i) => {
      sum += parseFloat(el.dataset.amount)
    })
    return sum
  }

  getTotalSendAmount () {
    let amount = 0
    this.amountTargets.forEach((el, i) => {
      amount += parseFloat(el.value)
    })
    return amount
  }

  calculateSelectedInputPercentage () {
    const sendAmount = this.getTotalSendAmount()
    const selectedInputSum = this.getSelectedInputsSum()
    let percentage = 0

    if (selectedInputSum >= sendAmount) {
      percentage = 100
    } else {
      percentage = (selectedInputSum / sendAmount) * 100
    }
    this.progressBarTarget.style.width = `${percentage}%`
  }

  openCustomizePanel () {
    let _this = this
    this.resetCustomizePanel()
    $('#custom-tx-row').slideDown()

    let accountNumber = _this.sourceAccountTarget.value
    let callback = function (txs) {
      // populate outputs
      let utxoHtml = txs.map(tx => {
        let receiveDateTime = new Date(tx.receive_time * 1000).toString().split(' ').slice(0, 5).join(' ')
        let dcrAmount = tx.amount / 100000000
        return `<tr>
                  <td width='5%'>
                    <input data-action='click->send#calculateSelectedInputPercentage' type='checkbox' class='custom-input' name='utxo' value='${tx.key}' data-amount='${dcrAmount}' />
                  </td>
                  <td width='40%'>${tx.address}</td>
                  <td width='15%'>${dcrAmount} DCR</td>
                  <td width='25%'>${receiveDateTime}</td>
                  <td width='15%'>${tx.confirmations} confirmation(s)</td>
                </tr>`
      }).join('')
      _this.customInputContentTarget.innerHTML = utxoHtml
      $('#custom-tx-row .status').hide()
      _this.show(_this.changeOutputsCardTarget)
    }
    this.getUnspentOutputs(accountNumber, callback)
  }

  getUnspentOutputs (accountNumber, successCallback) {
    this.nextButtonTarget.innerHTML = 'Loading...'
    this.nextButtonTarget.setAttribute('disabled', 'disabled')
    let _this = this

    let url = `/unspent-outputs/${accountNumber}`
    if (_this.spendUnconfirmedTarget.checked) {
      url += '?getUnconfirmed=true'
    }

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

  toggleChangeOutputPanel (event) {
    this.useCustomChangeOutput = !this.useCustomChangeOutput
    this.clearMessages()
    this.changeOutputContentTarget.innerHTML = ''
    this.numberOfChangeOutputsTarget.value = ''
  }

  toggleUseRandomChangeOutputs () {
    let numberOfChangeOutput = parseFloat(this.numberOfChangeOutputsTarget.value)
    if (!(numberOfChangeOutput > 0)) {
      return
    }
    this.generateChangeOutputs()
  }

  changeOutputNumberChanged () {
    let numberOfChangeOutput = parseFloat(this.numberOfChangeOutputsTarget.value)
    if (!(numberOfChangeOutput > 0)) {
      return
    }
    this.generateChangeOutputs()
  }

  changeOutputAmountChanged (event) {
    let _this = this
    let targetElement = event.currentTarget

    let index = parseInt(targetElement.getAttribute('data-index'))
    let amount = parseFloat(targetElement.value)

    let totalAmountEntered = 0
    _this.changeOutputAmountTargets.forEach(ele => {
      totalAmountEntered += parseFloat(ele.value)
    })
    let availableChange = _this.totalChangeAmount - (totalAmountEntered - amount)
    if (totalAmountEntered > _this.totalChangeAmount) {
      amount = availableChange
      targetElement.value = availableChange
    }

    _this.changeOutputPercentageTargets.forEach(function (ele) {
      if (parseInt(ele.getAttribute('data-index')) === index) {
        let percentage = amount / _this.totalChangeAmount * 100
        ele.value = percentage
      }
    })
  }

  changeOutputAmountPercentageChanged (event) {
    let _this = this
    let targetElement = event.currentTarget

    let index = parseInt(targetElement.getAttribute('data-index'))
    let percentage = parseFloat(targetElement.value)

    let totalAmountEntered = 0
    _this.changeOutputPercentageTargets.forEach(ele => {
      totalAmountEntered += parseFloat(ele.value)
    })
    if (totalAmountEntered > 100) {
      let availablePercentage = 100 - (totalAmountEntered - percentage)
      targetElement.value = availablePercentage
      percentage = availablePercentage
    }

    _this.changeOutputAmountTargets.forEach(function (ele) {
      if (parseInt(ele.getAttribute('data-index')) === index) {
        let amount = percentage / 100 * _this.totalChangeAmount
        ele.value = amount
      }
    })
  }

  generateChangeOutputs () {
    if (!this.validateSendForm()) {
      return
    }

    if (this.generatingChangeOutputs || !this.useCustomChangeOutput) {
      return
    }
    this.generatingChangeOutputs = true

    this.changeOutputContentTarget.innerHTML = ''

    let numberOfChangeOutput = parseFloat(this.numberOfChangeOutputsTarget.value)
    if (!(numberOfChangeOutput > 0)) {
      this.showError('Number of change outputs must be a non-zero positive number')
      return
    }

    this.generateOutputsButtonTarget.setAttribute('disabled', 'disabled')
    this.generateOutputsButtonTarget.innerHTML = 'Loading...'
    this.numberOfChangeOutputsTarget.setAttribute('disabled', 'disabled')

    let _this = this
    _this.getRandomChangeOutputs(numberOfChangeOutput, function (changeOutputdestinations) {
      if (!this.useCustomChangeOutput) {
        return
      }
      _this.totalChangeAmount = 0
      changeOutputdestinations.forEach(destination => {
        _this.totalChangeAmount += destination.Amount
      })
      changeOutputdestinations.forEach((destination, i) => {
        let template = document.importNode(_this.changeDestinationTemplateTarget.content, true)

        template.querySelector('input[name="change-output-address"]').value = destination.Address
        let percentage = 0
        if (_this.useRandomChangeOutputsTarget.checked) {
          template.querySelector('input[name="change-output-amount"]').value = destination.Amount
          percentage = destination.Amount / _this.totalChangeAmount * 100
        }
        template.querySelector('input[name="change-output-amount-percentage"]').value = percentage

        template.querySelector('input[name="change-output-address"]').setAttribute('data-index', i)
        template.querySelector('input[name="change-output-amount"]').setAttribute('data-index', i)
        template.querySelector('input[name="change-output-amount-percentage"]').setAttribute('data-index', i)

        _this.changeOutputContentTarget.appendChild(template)
      })

      _this.show(_this.changeOutputContentTarget)

      _this.changeOutputAmountTargets.forEach(function (ele) {
        ele.setAttribute('readonly', 'true')
      })
      if (_this.useRandomChangeOutputsTarget.checked) {
        _this.changeOutputPercentageTargets.forEach(function (ele) {
          ele.setAttribute('disabled', 'disabled')
        })
      }

      _this.generateOutputsButtonTarget.removeAttribute('disabled')
      _this.generateOutputsButtonTarget.innerHTML = 'Generate Change Outputs'
      _this.numberOfChangeOutputsTarget.removeAttribute('disabled')
      _this.generatingChangeOutputs = false
    }, function () {
      _this.generateOutputsButtonTarget.removeAttribute('disabled')
      _this.generateOutputsButtonTarget.innerHTML = 'Generate Change Outputs'
      _this.numberOfChangeOutputsTarget.removeAttribute('disabled')
      _this.generatingChangeOutputs = false
    })
  }

  getRandomChangeOutputs (numberOfOutputs, successCallback, completeCallback) {
    let queryParams = $('#send-form').serialize()
    queryParams += `&totalSelectedInputAmountDcr=${this.getSelectedInputsSum()}`

    // add source-account value to post data if source-account element is disabled
    if (this.sourceAccountTarget.disabled) {
      queryParams += `&source-account=${this.sourceAccountTarget.value}`
    }

    queryParams += `&nChangeOutput=${numberOfOutputs}`

    let _this = this

    axios.get('/random-change-outputs?' + queryParams).then((response) => {
      let result = response.data
      if (result.error !== undefined) {
        _this.setErrorMessage(result.error)
      } else {
        successCallback(result.message)
      }
    }).catch((error) => {
      console.log(error)
      _this.setErrorMessage('A server error occurred')
    }).then(() => {
      if (completeCallback) {
        completeCallback()
      }
    })
  }

  normalizeChangeOutputAmount () {
    let totalAmountEntered = 0
    this.changeOutputAmountTargets.forEach(ele => {
      totalAmountEntered += parseFloat(ele.value)
    })

    let availableChange = this.totalChangeAmount - totalAmountEntered

    if (availableChange > 0) {
      let changeOutputCount = this.numberOfChangeOutputsTarget.value
      this.changeOutputAmountTargets.forEach(ele => {
        let index = parseInt(ele.getAttribute('data-index'))
        if (index === changeOutputCount - 1) {
          ele.value = parseFloat(ele.value) + availableChange
        }
      })
    }
  }

  submitForm () {
    if (!this.validatePassphrase()) {
      return
    }

    this.normalizeChangeOutputAmount()

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

  resetSendForm () {
    this.resetCustomizePanel()
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

    $('#custom-tx-row').slideUp()
    this.hide(this.changeOutputsCardTarget)

    this.clearMessages()
  }

  resetCustomizePanel () {
    this.customTxRowTarget.querySelectorAll('tbody').forEach(el => {
      el.innerHTML = ''
    })
    this.customTxRowTarget.querySelectorAll('.status').forEach(el => {
      this.show(el)
    })
    this.customTxRowTarget.querySelectorAll('.alert-danger').forEach(el => {
      el.parentNode.removeChild(el)
    })

    this.hide(this.changeOutputsCardTarget)
    this.useRandomChangeOutputsTarget.checked = false
    this.numberOfChangeOutputsTarget.value = ''
    this.changeOutputContentTarget.innerHTML = ''
  }

  newDestination () {
    if (!this.validateDestinationsField()) {
      return
    }
    let destinationTemplate = document.importNode(this.destinationTemplateTarget.content, true)
    this.destinationsTarget.appendChild(destinationTemplate)
    if (this.destinationCount() > 1) {
      this.show(this.removeDestinationButtonTarget)
    }
  }

  removeDestination () {
    let count = this.destinationsTarget.querySelectorAll('div.destination').length
    if (!(count > 1)) {
      return
    }
    this.destinationsTarget.removeChild(this.destinationsTarget.querySelector('div.destination:last-child'))
    if (!(this.destinationCount() > 1)) {
      this.hide(this.removeDestinationButtonTarget)
    }
  }

  destinationCount () {
    return this.destinationsTarget.querySelectorAll('div.destination').length
  }

  validatePassphrase () {
    if (this.walletPassphraseTarget.value === '') {
      this.passwordErrorTarget.innerHTML = '<div class="error">Your wallet passphrase is required</div>'
      return false
    }

    return true
  }

  getWalletPassphraseAndSubmit () {
    this.clearMessages()
    if (!this.validateSendForm() || !this.validateChangeOutputAmount()) {
      return
    }
    $('#passphrase-modal').modal()
  }

  toggleUseCustom () {
    if (!this.useCustomTarget.checked) {
      $('#custom-tx-row').slideUp()
      this.resetCustomizePanel()
      this.hide(this.changeOutputsCardTarget)
      return
    }
    if (!this.validateDestinationsField()) {
      this.useCustomTarget.checked = false
      return
    }
    this.clearMessages()
    this.openCustomizePanel()
  }

  toggleSpendUnconfirmed () {
    if (this.useCustomTarget.checked) {
      this.resetCustomizePanel()
      this.openCustomizePanel()
    }
  }

  setErrorMessage (message) {
    this.hide(this.successMessageTarget)
    this.show(this.errorMessageTarget)
    this.errorMessageTarget.innerHTML = message
  }

  setSuccessMessage (message) {
    this.hide(this.errorMessageTarget)
    this.show(this.successMessageTarget)
    this.successMessageTarget.innerHTML = message
  }

  clearMessages () {
    this.hide(this.errorMessageTarget)
    this.hide(this.successMessageTarget)
    this.errorsTarget.innerHTML = ''
    this.successMessageTarget.innerHTML = ''
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
