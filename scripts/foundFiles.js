$('#fileRequestForm').submit(function(e){
    e.preventDefault();
    $.ajax({
        url:"/DownloadFoundFile",
        type:'post',
        data:$('#fileRequestForm').serialize(),
        success:function(){
            //whatever you wanna do after the form is successfully submitted
        }
    });
});