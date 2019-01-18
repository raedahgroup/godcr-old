// Renders latest peer count and best block to title
function connectionInfo(string){
  document.getElementById("bchn-status-h").style.display = "block";
  document.getElementById("bchn-status").innerHTML = string;
}
