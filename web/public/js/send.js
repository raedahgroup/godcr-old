/**==================================================================*
 *                    PASSPHRASE FUNCTIONS                           *
 *===================================================================*/

function validatePassphrase() {
    $(".passphrase-error").remove();
    var error = "";

    var passphraseEl = $("#wallet-passphrase");

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

    var sourceAccountEl = $("#source-account");
    var amountEl = $("#amount");
    var destinationAddressEl = $("#destination-address");
    var isClean = true;

    if (sourceAccountEl.find(":selected").text() == "") {
        errors["#source-account"] = "The source account is required";
    }

    if (amountEl.val() == "") {
        errors["#amount"] = "The amount is required"
    }

    if (destinationAddressEl.val() == "") {
        errors["#destination-address"] = "The destination address is required"
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
                var hash = "<strong>" + response.txHash + "</strong>";
                setSuccessMessage(hash)
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

function openCustomizePanel() {
    $("#customize-checkbox").prop("checked", false);
    if ((!$("#customize-panel").hasClass("show") || $("#customize-panel").css("display") == "none") && validateSendForm()) {
        $("form .collapse").slideUp();
        $("#customize-panel").slideDown();

        $("#customize-checkbox").prop("checked", true);
        $("#customize-panel .status").show();

        var account_number = $("#source-account").find(":selected").val();
        var callback = function(txs) {
            // populate outputs 
            var utxoHtml = txs.map(tx => {
                var receiveDateTime = new Date(tx.receive_time * 1000)
                return  "<tr>" + 
                            "<td width='5%'><input type='checkbox' name='tx' value="+ tx.key+" /></td>" +
                            "<td width='60%'>" + tx.key + "</td>" + 
                            "<td width='15%'>" + tx.amount_string + "</td>" + 
                            "<td width='20%'>" + receiveDateTime.toString() + "</td>" +
                        "</tr>"
            });
            $("#customize-panel tbody").html(utxoHtml.join('\n'));
            $("#customize-panel .status").hide();
        }
        getUnspentOutputs(account_number, callback);
    }
}


$(function(){
    $("#form-panel-card button").on("click", function(){
        if (!$("#form-panel").hasClass("show") || $("#form-panel").css("display") == "none") {
            $("form .collapse").slideUp().removeClass("show");
            $("#form-panel").slideDown();
        }
    });

    $("#customize-panel-card button").on("click", function(){
        openCustomizePanel();
    });

    $("#customize-checkbox").on("change", function(){
        if (this.checked) {
            openCustomizePanel();
        }
    })


    $("#submit-btn").on("click", function(e){
        e.preventDefault();
        if (validateSendForm()) {
            getWalletPassphraseAndSubmit(submitSendForm)
        }
    })

    // set current nav item 
    $("#nav-send").addClass("active");
});


