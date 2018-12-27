/**==================================================================*
 *                    PASSPHRASE FUNCTIONS                           *
 *===================================================================*/

function validatePassphrase() {
    $(".passphrase-error").remove();
    var error = "";

    var passphraseEl = $("#walletPassphrase");

    if (passphraseEl.val() == "") {
        error = "Your wallet passphrase is required";
    }

    if (error != "") {
        passphraseEl.after("<div class='passphrase-error'>" + error + "</div>");
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
 *                  SEND PAGE FUNCTIONS                              *
 *===================================================================*/
function validateSendForm() {
    // clear errors first
    $(".error").remove();
    var errors = {};

    var sourceAccountEl = $("#sourceAccount");
    var amountEl = $("#amount");
    var destinationAddressEl = $("#destinationAddress");
    var isClean = true;

    if (sourceAccountEl.find(":selected").text() == "") {
        errors["#sourceAccount"] = "The source account is required";
    }

    if (amountEl.val() == "") {
        errors["#amount"] = "The amount is required"
    }

    if (destinationAddressEl.val() == "") {
        errors["#destinationAddress"] = "The destination address is required"
    }

    if (!$.isEmptyObject(errors)) {
        isClean = false;
        for (var i in errors) {
            $(i).after("<div class='error'>" + errors[i] + "</div>");
        }
    }

    return isClean;
}

function submitSendForm() {
    var form = $("#send-form");
    var submit_btn = $("#send-form #submit-btn");
    submit_btn.attr("disabled", "disabled").html("Sending...");

    $.ajax({
        url: form.attr("action"),
        method: "POST",
        data: form.serialize(),
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

function getUnspentOutputs(account_number, success_callback) {
    var next_btn = $(".next-btn");
    next_btn.attr("disabled", "disabled").html("Loading...");

    $.ajax({
        url: "/unspent-outputs/" + account_number,
        method: "GET",
        data: {},
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

function setProgressWidth(width) {
    $(".stepper .progress-bar").css("width", width + "%");
}

function goToTab(tab_item) {
    tab_item.trigger("click").addClass("stepper-active");
}

function closeTab(tab_tem) {
    tab_item.trigger("click");
}

$(function(){
    var stepper_index = 0;
    var steppers = $(".stepper .nav-link");

    $(".next-btn").on("click", function(){
        if (validateSendForm() && stepper_index <= steppers.length) {
            var account_number = $("#sourceAccount").find(":selected").val();
            var callback = function(utxos) {
                // populate outputs 
                var utxoHtml = utxos.map(utxo => {
                    var receiveDateTime = new Date(utxo.receive_time * 1000)
                    return  "<tr>" + 
                                "<td width='5%'><input type='checkbox' name='tx' value="+ utxo.key+" /></td>" +
                                "<td width='60%'>" + utxo.key + "</td>" + 
                                "<td width='15%'>" + utxo.amount / 100000000 + " DCR</td>" + 
                                "<td width='20%'>" + receiveDateTime.toString() + "</td>" +
                            "</tr>"
                });
                $("#wallet-outputs tbody").html(utxoHtml.join('\n'));

                stepper_index += 1;
                goToTab(steppers.eq(stepper_index));  
                setProgressWidth((stepper_index+1 / steppers.length) * 100); 
            }
            getUnspentOutputs(account_number, callback);
        }
    });

    $(".previous-btn").on("click", function(){
        if (stepper_index > 0) {
            steppers.eq(stepper_index).removeClass("stepper-active");
            stepper_index -= 1;
            goToTab(steppers.eq(stepper_index));
            setProgressWidth((stepper_index+1 / steppers.length) * 100);     
        }
    })

    $("#submit-btn").on("click", function(e){
        e.preventDefault();
        getWalletPassphraseAndSubmit()
    })
});


