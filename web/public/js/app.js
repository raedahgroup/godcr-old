/**==================================================================*
 *                    VALIDATE FUNCTIONS                             *
 *===================================================================*/
 function validateFormFields() {
    // clear errors first 
    $(".error").remove()
    var isClean = true;

    $('[data-validate="true"]').each(function(){
        var value = "";
        var elementType = $(this).prop("type");
        var elementName = $(this).prop("name").replace(/(\-\w)/g, function(text){return " " + text[1].toUpperCase();});
        elementName = elementName.charAt(0).toUpperCase() + elementName.slice(1);
        

        if (elementType && elementType.toLowerCase() === 'radio') {
            value = $(this).filter(":checked").val();
        } else {
            value = $(this).val();
        }

        if ($(this).data("required") === true && value === "") {
            isClean = false;
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
    $("#walletPassphrase").val("");

    $("#passphrase-submit").on("click", function(){
        if (validatePassphrase()) {
            passphraseModal.modal("hide");
            submitFunction();
        }
    });

    passphraseModal.modal();
}

function validatePassphrase() {
    $(".passphrase-error").remove();
    var error = "";

    var passphraseEl = $("#walletPassphrase");

    if (passphraseEl.val() == "") {
        error = "Your wallet passphrase is required";
    }

    if (error != "") {
        passphraseEl.after("<div class='passphrase-error error'>" + error + "</div>");
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
        success: function(response){
            if (typeof onSuccessFunc === "function") {
                onSuccessFunc(response);
            }
        },
        error: function(){
            if (typeof onErrorFunc === "function") {
                onErrorFunc();
            }
        },
        complete: function(){
            if (typeof onCompleteFunc === "function") {
                onCompleteFunc();
            }
        }
    });
}