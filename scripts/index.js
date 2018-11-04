const messageButton = document.getElementById("messageButton");
const container = document.getElementById("container");
const form = document.forms.namedItem("fileInfo");
form.addEventListener('submit', function(ev) {

    var oOutput = document.querySelector("div"),
    oData = new FormData(form);
    console.log("Bjorn");

    var oReq = new XMLHttpRequest();
    oReq.open("POST", "/AddFile", true);
    oReq.onload = function(oEvent) {
        if (oReq.status == 200) {
            oOutput.innerHTML = "Uploaded!";
        } 
        else {
            oOutput.innerHTML = "Error " + oReq.status + " occurred when trying to upload your file.<br \/>";
        }
    };

    oReq.send(oData);
    ev.preventDefault();
}, false);

messageButton.addEventListener("click", async (e) => {
    e.preventDefault()
    $.ajax({
        url : '/GetMessages',
        method : 'get',
        success : (data) => {
            container.innerHTML = "";
            container.innerHTML = data
        }
    });
});

$('#messageForm').submit(function(e){
    e.preventDefault();
    $.ajax({
        url:'/AddMessage',
        type:'post',
        data:$('#messageForm').serialize(),
        success:function(){
            //whatever you wanna do after the form is successfully submitted
        }
    });
});