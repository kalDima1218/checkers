let url = "http://" + window.location.host
$(".logout_btn").click(function(){
    window.location.href = url+"/logout";
});

$(".new_game_btn").click(function(){
    window.location.href = url+"/start_game";
});

$(".enter_game_btn").click(function(){
    let id = document.getElementById("id").value;
    window.location.href = url+"/game?id=" + id;
});

$(".new_bot_game_btn").click(function(){
    window.location.href = url+"/start_bot_game";
});

function update(){
    let request = new XMLHttpRequest();
    request.open("POST", url+"/get_waiting");
    request.responseType = "text";
    request.send();
    request.onload = function() {
        let id = request.response;
        if(id !== "wrong"){
            window.location.href = url+"/game?id=" + id;
        }
    }
}

setInterval(update, 300)