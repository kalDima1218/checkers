let url = "http://" + window.location.host;
let selected = [-1, -1];
let current_turn = -1;
let last_turn = -1;
let board;
let side = 0;
let params = window
    .location
    .search
    .replace('?','')
    .split('&')
    .reduce(
		function(p,e){
			let a = e.split('=');
			p[decodeURIComponent(a[0])] = decodeURIComponent(a[1]);
			return p;
			},
		{}
		);

function rerender(){
	let id = params["id"];
	let request1 = new XMLHttpRequest();
	request1.open("GET", url+"/get_board_hist?id="+params["id"]+"&turn="+current_turn);
	console.log(current_turn)
	request1.responseType = "text";
	request1.send();
	request1.onload = function() {
		board = JSON.parse(request1.response);
		let new_board = ""
		if(side === 1) {
			for (let i = 0; i < 8; i += 1) {
				for (let j = 0; j < 8; j += 1) {
					new_board += '<div class="field_square">';
					if (board[i][j] === 0) {
						new_board += '<p class="void-piece" id="cell-' + i.toString() + '-' + j.toString() + '"></p>';
					} else if (board[i][j] === 1) {
						new_board += '<p class="red-piece cell" id="cell-' + i.toString() + '-' + j.toString() + '"></p>';
					} else if (board[i][j] === 2) {
						new_board += '<p class="black-piece cell" id="cell-' + i.toString() + '-' + j.toString() + '"></p>';
					} else if (board[i][j] === 3) {
						new_board += '<p class="red-piece king cell" id="cell-' + i.toString() + '-' + j.toString() + '"></p>';
					} else if (board[i][j] === 4) {
						new_board += '<p class="black-piece king cell" id="cell-' + i.toString() + '-' + j.toString() + '"></p>';
					}
					new_board += '</div>\n';
				}
				new_board += '<div style="display: none"></div>\n';
			}
		}
		else{
			for (let i = 0; i < 8; i += 1) {
				for (let j = 0; j < 8; j += 1) {
					new_board += '<div class="field_square">';
					if (board[7-i][7-j] === 0) {
						new_board += '<p class="void-piece" id="cell-' + (7-i).toString() + '-' + (7-j).toString() + '"></p>';
					} else if (board[7-i][7-j] === 1) {
						new_board += '<p class="red-piece cell" id="cell-' + (7-i).toString() + '-' + (7-j).toString() + '"></p>';
					} else if (board[7-i][7-j] === 2) {
						new_board += '<p class="black-piece cell" id="cell-' + (7-i).toString() + '-' + (7-j).toString() + '"></p>';
					} else if (board[7-i][7-j] === 3) {
						new_board += '<p class="red-piece king cell" id="cell-' + (7-i).toString() + '-' + (7-j).toString() + '"></p>';
					} else if (board[7-i][7-j] === 4) {
						new_board += '<p class="black-piece king cell" id="cell-' + (7-i).toString() + '-' + (7-j).toString() + '"></p>';
					}
					new_board += '</div>\n';
				}
				new_board += '<div style="display: none"></div>\n';
			}
		}
		document.getElementById("board").innerHTML = new_board
		$(".cell").click(function() {
			for(let i of document.getElementsByClassName("selected")){
				i.classList.remove("selected");
			}
			$(this).addClass("selected")
			selected_id = $(this).attr("id").split("-")
			selected = [parseInt(selected_id[1]), parseInt(selected_id[2])];
		});
		$(".void-piece").click(function() {
			if(selected[0] === -1){
				return;
			}
			xy = $(this).attr("id").split("-")
			let request = new XMLHttpRequest();
			request.open("GET", url+"/make_move?id=" + id + "&from_x=" + selected[0] + "&from_y=" + selected[1] + "&to_x=" + xy[1] + "&to_y=" + xy[2]);
			request.responseType = "text";
			selected = [-1, -1]
			request.send();
			request.onload = function(){}
		});
		if(selected[0] !== -1){
			document.getElementById("cell-" + selected[0].toString() + "-" + selected[1].toString()).classList.add("selected");
		}
	};
	for(let i of document.getElementsByClassName("selected")){
		i.classList.remove("selected");
	}
	let request2 = new XMLHttpRequest();
	request2.open("GET", url+"/who_win?id=" + id);
	request2.responseType = "text";
	request2.send();
	request2.onload = function() {
		if(request2.response == "0"){
			document.getElementById("win-red").classList.remove("hidden");
		}
		else if(request2.response == "1"){
			document.getElementById("win-black").classList.remove("hidden");
		}
	}
}

$(".end_turn_btn").click(function(){
	let request = new XMLHttpRequest();
	request.open("GET", url+"/end_move?id="+params["id"]);
	request.responseType = "text";
	request.send();
	request.onload = function(){}
});

$(".prev_turn_btn").click(function(){
	if(current_turn === 0){
		return;
	}
	current_turn -= 1;
	rerender();
});

$(".next_turn_btn").click(function(){
	if(current_turn === last_turn){
		return;
	}
	current_turn += 1;
	rerender();
});

function setPlayers() {
	let request = new XMLHttpRequest();
	request.open("GET", url + "/get_players?id=" + params["id"]);
	request.responseType = "text";
	request.send();
	request.onload = function () {
		let players = JSON.parse(request.response);
		document.getElementById("player0").innerHTML = "Красный: " + players[0];
		document.getElementById("player1").innerHTML = "Черный: " + players[1];
	}
	let request1 = new XMLHttpRequest();
	request1.open("GET", url + "/get_side?id=" + params["id"]);
	request1.responseType = "text";
	request1.send();
	request1.onload = function () {
		if(request1.response === "0"){
			side = 0;
		}
		else{
			side = 1;
		}
	}
}

function update() {
	let id = params["id"];
	let request1 = new XMLHttpRequest();
	request1.open("GET", url + "/get_last_move_number?id=" + params["id"]);
	request1.responseType = "text";
	request1.send();
	request1.onload = function () {
		let new_last_turn = JSON.parse(request1.response);
		if(new_last_turn !== last_turn){
			last_turn = new_last_turn;
			current_turn = last_turn;
			rerender();
		}
	}
	let request2 = new XMLHttpRequest();
	request2.open("GET", url+"/whose_move?id=" + id);
	request2.responseType = "text";
	request2.send();
	request2.onload = function() {
		if(request2.response === "0"){
			document.getElementById("turn").innerHTML = "Ходит: красный";
		}
		else{
			document.getElementById("turn").innerHTML = "Ходит: черный";
		}
	}
}


setPlayers();

setInterval(update, 300);