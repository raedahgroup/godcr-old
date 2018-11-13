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

function getWalletPassphraseAndSubmit(submitFunc) {
    var passphraseModal = $("#passphrase-modal");
    
    $("#passphrase-submit").on("click", function(){
        if (validatePassphrase()) {
            passphraseModal.modal('hide');
            submitFunc($("#walletPassphrase").val());
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

function submitSendForm(passphrase) {
    var form = $("#send-form");
    var submit_btn = $("#send-form #submit-btn");
    submit_btn.attr("disabled", "disabled").html("Sending...");

    $.ajax({
        url: form.attr("action"),
        method: "POST",
        data: form.serialize(),
        success: function(response) {
            if (typeof response.success != "undefined") {
                var m = "The transaction was published successfully. Hash: <strong>" + response.success + "</strong>";
                setSuccessMessage(m)
            } else {
                setErrorMessage(response.error) 
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
    $("#submit-btn").on("click", function(){
        if (validateSendForm()) {
            getWalletPassphraseAndSubmit(submitSendForm);
        }
    });
});
