let url = "http://" + window.location.host
$(".submit_btn").click(function(){
    let name = document.getElementById("name").value;
    let login = document.getElementById("login").value;
    let password = document.getElementById("password").value;
    let request = new XMLHttpRequest();
    request.open("POST", url+"/registration?name=" + name + "&login=" + login + "&password=" + password);
    request.responseType = "text";
    request.send();
    request.onload = function(){
        let response = request.response
        if(response !== "ok"){
            alert(response);
        }
        else{
            window.location.href = url+"/";
        }
    }
});