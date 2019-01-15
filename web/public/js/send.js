/**==================================================================*
 *                  SEND PAGE FUNCTIONS                              *
 *===================================================================*/
function validateAmountField() { 
    if ($("#amount").val() === "") {
        $(".errors").html("<div class='error'>Please enter an amount first</div>");
        return false;
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

    isClean = validateAmountField();

    if ($("#destination-address").val() === "") {
        errors.push("The destination address is required");
    }

    if ($("#use-custom").prop("checked") && (getSelectedInputsSum() < $("#amount").val()) ) {
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

function calculateSelectedInputPercentage() {
    var sendAmount = $("#amount").val();
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
                        "<td width='50%'>" + tx.key + "</td>" + 
                        "<td width='20%'>" + dcrAmount + " DCR</td>" + 
                        "<td width='25%'>" + receiveDateTime.toString().split(' ').slice(0,5).join(' '); + "</td>" +
                    "</tr>"
        });
        $("#custom-tx-row tbody").html(utxoHtml.join('\n'));
        $("#custom-tx-row .status").hide();

        // register check listener 
        $(".custom-input").on("click", function(){
           calculateSelectedInputPercentage();
        });

        $("#amount").on("keyup", function(){
            validateAmountField();
            calculateSelectedInputPercentage();
        });
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
                success_callback(response.message);
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
    
    $("#passphrase-submit").on("click", function(){
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
            if (validateAmountField()) {
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
            } else {
                get_unconfirmed_utxos = false;
            }

            resetCustomizePanel();
            openCustomizePanel(get_unconfirmed_utxos);
        }
    });

    $("#submit-btn").on("click", function(e){
        e.preventDefault();
        if (validateSendForm()) {
            getWalletPassphraseAndSubmit();
        }
    });
});


