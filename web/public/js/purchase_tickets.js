function submitPurchaseTicketsForm() {
    $("#success-msg-container").empty();
    $(".alert-danger").empty();

    var form = $("#purchase-tickets-form");
    var submit_btn = $("#next-btn");
    var postData = form.serialize();

    // reset passphrase-field 
    $("#wallet-passphrase").val("");

    if ($("#source-account").prop("disabled")) {
        postData += "&source-account=" + $("#source-account").val();
    }

    var requestInfo = {
        data: postData,
        url: "/purchase-tickets"
    };

    var successFunc = function(response){
        if (!response.success) {
            setErrorMessage(response.message)
        } else {
            $("#success-msg-container").append("<p>You have purchased " + response.message.length + " ticket(s)</p>");
            for (var i in response.message) {
                $("#success-msg-container").append("<p><strong>" + response.message[i] + "</strong></p>");
            }
        }
    };

    var errorFunc = function(){
        setErrorMessage("A server error occurred");
    };

    var completeFunc = function() {
        submit_btn.removeAttr("disabled").html("Next");
    };

    submit_btn.attr("disabled", "disabled").html("Sending...");
    makePostRequest(requestInfo, successFunc, errorFunc, completeFunc)
}

$(function(){
    $("#next-btn").on("click", function(e){
        e.preventDefault();
        if (validateFormFields()) {
            getWalletPassphraseAndSubmit(submitPurchaseTicketsForm);
        }
    })
});