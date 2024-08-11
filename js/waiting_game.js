let url = "http://" + window.location.host

function update(){
    let request = new XMLHttpRequest();
    request.open("POST", url+"/get_waiting");
    request.responseType = "text";
    request.send();
    request.onload = function() {
        let id = request.response;
        if(id !== "wait"){
            window.location.href = url+"/game?id=" + id;
        }
    }
}

setInterval(update, 250)