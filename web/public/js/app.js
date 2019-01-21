/**==================================================================*
 *                    VALIDATE FUNCTIONS                             *
 *===================================================================*/
 function validateFormFields() {
    // clear errors first 
    $(".error").remove()
    var isClean = true;

    $('[data-validate="true"]').each(function(){
        var self = $(this);
        var elementType = self.prop("type");
        
        var value = "";
        if (elementType && elementType.toLowerCase() === 'radio') {
            value = self.filter(":checked").val();
        } else {
            value = self.val();
        }

        if (self.data("required") === true && value === "") {
            isClean = false;
            // convert input name to uppercase words to display errors
            var elementName = self.prop("name").replace(/(\-\w)/g, function(text){return " " + text[1].toUpperCase();});
            elementName = elementName.charAt(0).toUpperCase() + elementName.slice(1);
            $(".errors").append("<div class='error'>" + elementName + " is required</div>");
        }
    });

    return isClean;
}

/**==================================================================*
 *                    PASSPHRASE FUNCTIONS                           *
 *===================================================================*/
function getWalletPassphraseAndSubmit(submitFunction) {
    var passphraseModal = $("#passphrase-modal");

    $("#passphrase-submit").off("click").on("click", function(){
        if (validatePassphrase()) {
            passphraseModal.modal("hide");
            submitFunction();
        }
    });
    passphraseModal.modal();
}

function validatePassphrase() {
    $(".passphrase-error").remove();
    var passphraseEl = $("#walletPassphrase");

    if (passphraseEl.val() === "") {
        passphraseEl.after("<div class='passphrase-error error'>Your wallet passphrase is required</div>");
        return false
    }

    return true
}


/**==================================================================*
 *                      GENERAL                                      *
 *===================================================================*/
function setErrorMessage(message) {
    $(".alert-success").hide();
    $(".alert-danger").html(message).show();
}

function setSuccessMessage(message) {
    $(".alert-danger").hide();
    $(".alert-success").html(message).show();
}

function clearMessages() {
    $(".alert-success").hide();
    $(".alert-danger").hide();
}

/**====================================================================*
 *                            AJAX ABSTRACTION                         *
 *=====================================================================*/
function makeGetRequest(requestInfo, onSuccessFunc, onErrorFunc, onCompleteFunc) {
    if (typeof requestInfo.data === "undefined") {
        requestInfo.data = {}
    }

    requestInfo.method = "GET"
    makeRequest(requestInfo, onSuccessFunc, onErrorFunc, onCompleteFunc)
}

function makePostRequest(requestInfo, onSuccessFunc, onErrorFunc, onCompleteFunc) {
    if (typeof requestInfo.data === "undefined") {
        requestInfo.data = {}
    }

    requestInfo.method = "POST"
    makeRequest(requestInfo, onSuccessFunc, onErrorFunc, onCompleteFunc)
}

function makeRequest(requestInfo, onSuccessFunc, onErrorFunc, onCompleteFunc) {
    $.ajax({
        url: requestInfo.url,
        method: requestInfo.method,
        data: requestInfo.data,
        success: onSuccessFunc,
        error: onErrorFunc,
        complete: onCompleteFunc
    });
}