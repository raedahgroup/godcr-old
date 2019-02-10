function submitPurchaseTicketsForm() {
    $("#success-msg-container").empty();
    $(".alert-danger").empty();

    var form = $("#purchase-tickets-form");
    var submit_btn = $("#purchase-btn");
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
            var successMsg = ["<p>You have purchased " + response.message.length + " ticket(s)</p>"];
            var ticketHashes = response.message.map(ticketHash => "<p><strong>" + ticketHash + "</strong></p>");
            successMsg.push(ticketHashes);
            setSuccessMessage(successMsg.join(""));
        }
    };

    var errorFunc = function(){
        setErrorMessage("A server error occurred");
    };

    var completeFunc = function() {
        submit_btn.removeAttr("disabled").html("Purchase");
    };

    submit_btn.attr("disabled", "disabled").html("Sending...");
    makePostRequest(requestInfo, successFunc, errorFunc, completeFunc)
}

$(function(){
    $("#next-btn").on("click", function(e){
        e.preventDefault();
        clearMessages();
        if (validateFormFields()) {
            getWalletPassphraseAndSubmit(submitPurchaseTicketsForm);
        }
    })
});