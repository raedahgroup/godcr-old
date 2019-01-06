/**==================================================================*
 *                    PASSPHRASE FUNCTIONS                           *
 *===================================================================*/

function validatePassphrase() {
    $(".passphrase-error").remove();
    var error = "";

    var passphraseEl = $("#wallet-passphrase");

    if (passphraseEl.val() === "") {
        error = "Your wallet passphrase is required";
    }

    if (error !== "") {
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
function validateAmountField() {
    var amountEl = $("#amount");
    if (amountEl.val() === "") {
        amountEl.after("<div class='error'>Please enter an amount first</div>");
        return false;
    }

    return true;
}
function validateSendForm() {
    // clear errors first
    $(".error").remove();
    var errors = {};

    var sourceAccountEl = $("#source-account");
    var amountEl = $("#amount");
    var destinationAddressEl = $("#destination-address");
    var isClean = true;

    if (sourceAccountEl.find(":selected").text() === "") {
        errors["#source-account"] = "The source account is required";
    }

    if (amountEl.val() === "") {
        errors["#amount"] = "The amount is required"
    }

    if (destinationAddressEl.val() === "") {
        errors["#destination-address"] = "The destination address is required"
    }

    var isCustomTransaction = $("#use-custom").prop("checked")?true:false;
    if (isCustomTransaction && (getSelectedInputsSum() < amountEl.val())) { 
        errors["#custom-tx-row .alert-info"] = "The sum of selected inputs is less than send amount";
    }

    if (!$.isEmptyObject(errors)) {
        isClean = false;
        for (var i in errors) {
            $(".errors").append("<div class='error'>" + errors[i] + "</div>");
        }
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

function resetCustomizePanel() {
    $("#custom-tx-row tbody").empty();
    $("#custom-tx-row .status").show();
    $("#custom-tx-row .alert-danger").remove();
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

function openCustomizePanel() {
    $("#custom-tx-row").removeClass("d-none");
    resetCustomizePanel();
   
    var account_number = $("#source-account").find(":selected").val();
    var callback = function(txs) {
        // populate outputs 
        var utxoHtml = txs.map(tx => {
            var receiveDateTime = new Date(tx.receive_time * 1000)
            return  "<tr>" + 
                        "<td width='5%'><input type='checkbox' class='custom-input' name='tx' value="+ tx.key+" data-amount='" + tx.formatted_amount + "' /></td>" +
                        "<td width='50%'>" + tx.key + "</td>" + 
                        "<td width='20%'>" + tx.formatted_amount + "</td>" + 
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
    getUnspentOutputs(account_number, callback);
}


$(function(){
    $("#use-custom").on("change", function(){
        if (this.checked) {
            if (validateAmountField()) {
                openCustomizePanel();
            } else {
                $(this).prop("checked", false);
            }
        } else {
            resetCustomizePanel();
            $("#custom-tx-row").addClass("d-none");
        }
    });

    // clear validation errors on type 
    $("input[type=text], input[type=number]").each(function(){
        $(this).on("keyup", function(){
            if ($(this).val() !== "") {
                $(this).next(".error").remove();
            }
        });
    });

    $("#submit-btn").on("click", function(e){
        e.preventDefault();
        if (validateSendForm()) {
            getWalletPassphraseAndSubmit(submitSendForm)
        }
    })

    // set current nav item 
    $("#nav-send").addClass("active");
});


