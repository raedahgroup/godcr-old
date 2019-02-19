/**==================================================================*
 *                  SEND PAGE FUNCTIONS                              *
 *===================================================================*/
function validateDestinationFields() {
    for (const field of $('#destinations input')) {
        if ($(field).val() === "") {
            $(".errors").html("<div class='error'>The destination address and amount are required</div>");
            return false;
        }
        if ($(field).hasClass("amount") && !(parseFloat($(field).val()) > 0)) {
            $(".errors").html("<div class='error'>Amount must be a non-zero positive number</div>");
            return false;
        }
    }
    return true;
}

function validateSendForm() {
    // clear errors first
    $(".errors").empty();
    var errors = [];
    var isClean = true;

    if ($("#source-account").find(":selected").text() === "") {
        errors.push("The source account is required");
    }

    isClean = validateDestinationFields();

    if ($("#use-custom").prop("checked") && (getSelectedInputsSum() < getTotalSendAmount()) ) {
        errors.push("The sum of selected inputs is less than send amount");
    }

    if (errors.length > 0) {
        for (var i in errors) {
          $(".errors").append("<div class='error'>" + errors[i] + "</div>");
        }
        isClean = false;
    }

    return isClean;
}

function getSelectedInputsSum() {
    var sum = 0;
    $(".custom-input:checked").each(function(){
       sum += $(this).data("amount");
    });

    return sum
}

function getTotalSendAmount() {
    let amount = 0;
    for (const field of $('#destinations .amount')) {
        amount += parseFloat($(field).val());
    }
    return amount;
}

function calculateSelectedInputPercentage() {
    var sendAmount = getTotalSendAmount();
    var selectedInputSum = getSelectedInputsSum();
    var percentage = 0;

    if (selectedInputSum >= sendAmount) {
        percentage = 100;
    } else {
        percentage = (selectedInputSum / sendAmount) * 100;
    }

    $("#custom-tx-row .progress-bar").css("width", percentage+"%");
}


function resetCustomizePanel() {
    $("#custom-tx-row tbody").empty();
    $("#custom-tx-row .status").show();
    $("#custom-tx-row .alert-danger").remove();

    $("#changeOutputsCard").hide();
    $("#autoChangeOutputsDiv").hide();
    $("#autoChangeOutputsDiv").html('')
    $('#numberOfChangeOutput').val('')
}


function openCustomizePanel(get_unconfirmed) {
    resetCustomizePanel();
    $("#custom-tx-row").slideDown();
   
    var account_number = $("#source-account").find(":selected").val();
    var callback = function(txs) {
        // populate outputs 
        var utxoHtml = txs.map(tx => {
            var receiveDateTime = new Date(tx.receive_time * 1000);
            var dcrAmount = tx.amount / 100000000;
            return  "<tr>" + 
                        "<td width='5%'><input type='checkbox' class='custom-input' name='utxo' value="+ tx.key+" data-amount='" + dcrAmount + "' /></td>" +
                        "<td width='40%'>" + tx.address + "</td>" +
                        "<td width='15%'>" + dcrAmount + " DCR</td>" +
                        "<td width='25%'>" + receiveDateTime.toString().split(' ').slice(0,5).join(' ') + "</td>" +
                        "<td width='15%'>" + tx.confirmations + " confirmation(s)</td>" +
                    "</tr>"
        });
        $("#custom-tx-row tbody").html(utxoHtml.join('\n'));
        $("#custom-tx-row .status").hide();

        // register check listener 
        $(".custom-input").on("click", function(){
           calculateSelectedInputPercentage();
        });

        $("#destinations").on("keyup", ".amount", function(){
            validateDestinationFields();
            calculateSelectedInputPercentage();
        });

        $("#changeOutputsCard").show();
    }
    getUnspentOutputs(account_number, get_unconfirmed, callback);
}

function getUnspentOutputs(account_number, get_unconfirmed, success_callback) {
    var next_btn = $(".next-btn");
    next_btn.attr("disabled", "disabled").html("Loading...");

    var data = {}
    if (get_unconfirmed) {
        data.getUnconfirmed = true
    }

    $.ajax({
        url: "/unspent-outputs/" + account_number,
        method: "GET",
        data: data,
        success: function(response) {
            if (response.success) {
                success_callback(response.message)
            } else {
                setErrorMessage(response.message);
            }
        },
        error: function(error) {
            setErrorMessage("A server error occurred")
        },
        complete: function() {
            next_btn.removeAttr("disabled").html("Next");
        }
    })
}

var changeOutputDestinations = undefined;

function generateRandomChangeOutputs() {
    if (!validateSendForm()) {
        return
    }

    $("#autoChangeOutputsDiv").html('');

    var numberOfChangeOutputTxt = $("#numberOfChangeOutput");
    var numberOfChangeOutput = parseFloat(numberOfChangeOutputTxt.val());
    if(!(numberOfChangeOutput > 0)){
        setErrorMessage("Number of change outputs must be a non-zero positive number");
        return;
    }
    

    var generate_outputs_btn = $("generate-outputs-btn");
    generate_outputs_btn.attr("disabled", "disabled").html("Loading...");
    numberOfChangeOutputTxt.attr("disabled", "disabled")

    getRandomChangeOutputs(numberOfChangeOutput, function(changeOutputdestinations) {
        var outputsHtml = changeOutputdestinations.map(destination => {
            return  `<div class="row">
                <div class="col-md-6 col-sm-12">
                    <div class="form-group">
                        <label>Change Address</label>
                        <input type="text" class="form-control" readonly name="change-output-address" value="${destination.Address}" />
                    </div>
                </div>

                <div class="col-md-6 col-sm-12">
                    <div class="form-group">
                        <label>Amount (DCR)</label>
                        <input type="number" class="form-control change-amount" readonly name="change-output-amount" value="${destination.Amount}" />
                    </div>
                </div>
            </div>`; 
        });
        
        $("#autoChangeOutputsDiv").show();
        $("#autoChangeOutputsDiv").html(outputsHtml.join('\n'));
    }, function () {
        generate_outputs_btn.removeAttr("disabled").html("Send");
        numberOfChangeOutputTxt.removeAttr("disabled")
    });
}

function addCustomChangeDestination() {
    var addChangeOutputBtn = $("#add-change-destination-btn")
    var removeChangeOutputDestination = $("#remove-change-destination-btn")

    addChangeOutputBtn.attr("disabled", "disabled").html("Loading")
    removeChangeOutputDestination.attr("disabled", "disabled")

    getRandomChangeOutputs(1, function (changeOurputDestinations) {
        var outputsHtml = changeOurputDestinations.map(destination => {
            return  `<div class="row">
                <div class="col-md-6 col-sm-12">
                    <div class="form-group">
                        <label>Change Address</label>
                        <input type="text" class="form-control" readonly name="change-output-address" value="${destination.Address}" />
                    </div>
                </div>

                <div class="col-md-6 col-sm-12">
                    <div class="form-group">
                        <label>Amount (DCR)</label>
                        <input type="number" class="form-control change-amount" name="change-output-amount" value="${destination.Amount}" />
                    </div>
                </div>
            </div>`; 
        });

        $("#custom-change-destinations").append(outputsHtml)
    }, function () {
        addChangeOutputBtn.removeAttr("disabled").html("Add another address");
        removeChangeOutputDestination.removeAttr("disabled")
        if($("#custom-change-destinations .row").length > 1) {
            $("#remove-change-destination-btn").show();
        }
    })
}

function getRandomChangeOutputs(number_of_outputs, success_callback, complete_callback) {
    var postData = $("#send-form").serialize();
    postData += "&totalSelectedInputAmountDcr=" + getSelectedInputsSum();

    // add source-account value to post data if source-account element is disabled
    if ($("#source-account").prop("disabled")) {
        postData += "&source-account=" + $("#source-account").val();
    }

    postData += "&nChangeOutput=" + number_of_outputs;

    $.ajax({
        url: "/random-change-outputs",
        method: "POST",
        data: postData,
        success: function(response) {
            if (response.error) {
                setErrorMessage(response.error)
            } else {
                changeOutputDestinations = response.message;
                success_callback(response.message);
            }
        },
        error: function(error) {
            changeOutputDestinations = undefined;
            setErrorMessage("A server error occurred")
        },
        complete: function() {
            if (complete_callback) {
                complete_callback()
            }
        }
    })
}

function submitSendForm() {
    var form = $("#send-form");
    var submit_btn = $("#send-form #submit-btn");
    submit_btn.attr("disabled", "disabled").html("Sending...");

    var postData = form.serialize();
    postData += "&totalSelectedInputAmountDcr=" + getSelectedInputsSum();

    // add source-account value to post data if source-account element is disabled
    if ($("#source-account").prop("disabled")) {
        postData += "&source-account=" + $("#source-account").val();
    }

    $.ajax({
        url: form.attr("action"),
        method: "POST",
        data: postData,
        success: function(response) {
            if (response.error) {
                setErrorMessage(response.error)
            } else {
                var txHash = "The transaction was published successfully. Hash: <strong>" + response.txHash + "</strong>";
                setSuccessMessage(txHash)
            }
        },
        error: function(error) {
            setErrorMessage("A server error occurred")
        },
        complete: function() {
            submit_btn.removeAttr("disabled").html("Send");
        }
    })
}

/**==================================================================*
 *                    MULTI-ADDRESS FUNCTIONS                        *
 *===================================================================*/
function newDestination() {
    let html = `<div class="row">
                    <div class="col-md-6 col-sm-12">
                        <div class="form-group">
                            <label>Destination Address</label>
                            <input type="text" class="form-control" name="destination-address" />
                        </div>
                    </div>
                    <div class="col-md-6 col-sm-12">
                        <div class="form-group">
                            <label for="amount-">Amount (DCR)</label>
                            <input type="number" class="form-control amount" name="destination-amount" />
                        </div>
                    </div>
                </div>
    `
    $("#destinations").append(html)
}

/**==================================================================*
 *                    PASSPHRASE FUNCTIONS                           *
 *===================================================================*/

function validatePassphrase() {
    if ($("#walletPassphrase").val() === "") {
        $("#passphrase-modal .errors").html("<div class='error'>Your wallet passphrase is required</div>");
        return false
    }

    return true
}

function getWalletPassphraseAndSubmit() {
    var passphraseModal = $("#passphrase-modal");
    var submitPassphrase = $("#passphrase-submit");
    
    submitPassphrase.off("click");
    submitPassphrase.on("click", function(){
        if (validatePassphrase()) {
            passphraseModal.modal('hide');
            submitSendForm();
        }
    });
    passphraseModal.modal();
}


/**==================================================================*
 *                      GENERAL                                      *
 *===================================================================*/

function setErrorMessage(message) {
    $(".alert-success").hide();
    $(".alert-danger").html(message).show();
}

function setSuccessMessage(message) {
    var m = "The transaction was published successfully. Hash: <strong>" + message + "</strong>";
    $(".alert-danger").hide();
    $(".alert-success").html(m).show();
}

function clearMessages() {
    $(".alert-success").hide();
    $(".alert-danger").hide();
}

$(function(){
    $("#use-custom").on("change", function(){
        if (this.checked) {
            if (validateDestinationFields()) {
                $(".errors").empty();
                var get_unconfirmed = false;
                if ($("#spend-unconfirmed").is(":checked")) {
                    get_unconfirmed = true;
                }
                openCustomizePanel(get_unconfirmed);
            } else {
                $(this).prop("checked", false);
            }
        } else {
            resetCustomizePanel();
            $("#custom-tx-row").slideUp();
        }
    });
    
    $("#spend-unconfirmed").on("change", function(){
        var use_custom = $("#use-custom").is(":checked");

        if (use_custom) {
            var get_unconfirmed_utxos = false;
            if (this.checked) {
                get_unconfirmed_utxos = true;
            }

            resetCustomizePanel();
            openCustomizePanel(get_unconfirmed_utxos);
        }
    });

    $("#add-destination-btn").on("click", function () {
        newDestination();
        $("#remove-destination-btn").show();
    })

    $("#remove-destination-btn").hide();

    $("#remove-destination-btn").on("click", function () {
        $("#destinations .row:last-child").remove();
        if($("#destinations .row").length < 2) {
            $("#remove-destination-btn").hide();
        }
    })


    $("#changeOutputsCard").hide();
    $("#autoChangeOutputsDiv").hide();

    $("#automaticChangeOutputPnl").on("hide.bs.collapse", function () {
        $("#autoChangeOutputsDiv").html("")
    })

    $("#generate-outputs-btn").on("click", function() {
        generateRandomChangeOutputs()
    })

    $("#customChangeOutputPnl").on("shown.bs.collapse", function () {
        addCustomChangeDestination()
    })

    $("#customChangeOutputPnl").on("hide.bs.collapse", function () {
        $("#custom-change-destinations").html("")
    })

    $("#add-change-destination-btn").on("click", function () {
        addCustomChangeDestination()
    })

    $("#remove-change-destination-btn").hide()
    $("#remove-change-destination-btn").on("click", function () {
        $("#custom-change-destinations .row:last-child").remove()
        if($("#custom-change-destinations .row").length < 2) {
            $("#remove-change-destination-btn").hide();
        }
    })

    $("#submit-btn").on("click", function(e){
        e.preventDefault();
        if (validateSendForm()) {
            getWalletPassphraseAndSubmit();
        }
    });
});


