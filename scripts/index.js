const messageButton = document.getElementById("messageButton");
const container = document.getElementById("container");
const form = document.forms.namedItem("fileInfo");
form.addEventListener('submit', function(ev) {

    var oOutput = document.querySelector("div"),
    oData = new FormData(form);
    console.log("Handling upload file");

    var oReq = new XMLHttpRequest();
    oReq.open("POST", "/AddFile", true);
    oReq.onload = function(oEvent) {
        if (oReq.status == 200) {
            oOutput.innerHTML = "Uploaded!";
            console.log("Uploaded!");
        } else {
            var message = "Error " + oReq.status + " occurred when trying to upload your file.<br \/>";
            oOutput.innerHTML = message;
            console.log(message);
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

$('#downloadForm').submit(function(e){
    e.preventDefault();
    $.ajax({
        url:'/RequestFile',
        type:'post',
        data:$('#fileRequestForm').serialize(),
        success:function(){
            //whatever you wanna do after the form is successfully submitted
        }
    });
});

$('#DownloadMetaFileForm').submit(function(e){
    e.preventDefault();
    $.ajax({
        url:'/DownloadMetaFile',
        type:'post',
        data:$('#DownloadMetaFileForm').serialize(),
        success:function(){
            //whatever you wanna do after the form is successfully submitted
        }
    });
});
