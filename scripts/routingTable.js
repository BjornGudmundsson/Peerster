$('#messageForm').submit(function(e){
    e.preventDefault();
    $.ajax({
        url:"/PostPrivateMessage",
        type:'post',
        data:$('#messageForm').serialize(),
        success:function(data){
            //whatever you wanna do after the form is successfully submitted
            console.log(data);
        }
    });
});