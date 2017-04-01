var Conn = {
	connect: function(username) {
		this.username = username;
		this.conn = new WebSocket("ws://localhost:8080/websocket/"+username);

		this.conn.onopen = function() {
			console.log("connected!");
		}

		this.conn.onmessage = function(event) {
			var data = JSON.parse(event.data);

			if (data.message == undefined || data.type == undefined) {
				console.log("something was undefined");
				return
			}

			if(data.type == "welcome") {
				$('.username-request').hide();
			}

			switch (data.type) {
				case "chat-message":
					$('.messages ul').append('<li><small>'+data.from+'</small><p>'+data.message+'</p> <small>sent: '+data.created_at+'</small></li>');
					break;
				case "refresh":
					break;
				case "control-message":
					console.log("got control-message:");
					console.log(data.message);
					break;
				default:
					console.log("message doesn't fit the types we have:");
					console.log(event.data);
			}
		};
	},

	send: function(message) {
		this.conn.send(JSON.stringify({username: this.username, message: message}));
	},

	recieve: function() {

	}
};

$(document).ready(function() {
	var conn = Conn;

	$('#start-chat').click(function(e) {
		e.preventDefault();

		var username = $('#username').val();

		conn.connect(username);

		$('#new-message-form').show();
	});

	$('#new-message-form').submit(function(e) {
		e.preventDefault();

		var message = $('#message').val();
		console.log(message);

		conn.send(message);

		$('#message').val('');
	})
});
