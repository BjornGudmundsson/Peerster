const messageButton = document.getElementById("messageButton");
const container = document.getElementById("container");

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