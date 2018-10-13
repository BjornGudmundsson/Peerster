const messageButton = document.getElementById("messageButton");
const container = document.getElementById("container");

messageButton.addEventListener("click", async (e) => {
    e.preventDefault()
    $.ajax({
        url : '/GetMessages',
        method : 'get',
        success : (data) => {
            console.log(data);
            container.innerHTML = "";
            container.innerHTML = data
        }
    });
});

$("#messageForm").submit((e)=>{
    $ajax({
        type : "POST",
        url : "/AddMessage",
        data : {
            text : $(this).text.value
        }
    });
    return false;
});