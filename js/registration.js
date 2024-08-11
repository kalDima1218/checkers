let url = "http://" + window.location.host

function hashText(text) {
    const msgUint8 = new TextEncoder().encode(text);
    return crypto.subtle.digest("SHA-512", msgUint8)
        .then(hashBuffer => {
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            const hashHex = hashArray
                .map(b => b.toString(16).padStart(2, "0"))
                .join("");
            return hashHex;
        });
}

$(".submit_btn").click(function(){
    let username = document.getElementById("username").value;
    let login = document.getElementById("login").value;
    let password = document.getElementById("password").value;
    hashText(password).then(passwordHash => {
        let request = new XMLHttpRequest();
        request.open("POST", url + "/registration?username=" + username + "&login=" + login + "&password=" + passwordHash);
        request.responseType = "text";
        request.send();
        request.onload = function () {
            let response = request.response
            if (response === "ok"){
                window.location.href = url + "/";
            } else if (response === "not free login") {
                alert("Логин занят");
            } else if (response === "too long") {
                alert("Логин или ник слишком длинный")
            } else {
                alert("Что-то пошло не так")
            }
        }
    });
});