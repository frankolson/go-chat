$(function() {
  let websocket = new WebSocket('ws://' + window.location.host + '/websocket');
  let room = $('#chat-room');

  websocket.addEventListener('message', function(event) {
    let data = JSON.parse(event.data);
    let message = `<p><strong>${data.username}</strong>: ${data.text}</p>`;
    
    room.append(message);
    room.scrollTop = room.scrollHeight;
  });

  $('#chat-form').submit(function(event) {
    event.preventDefault();
    let username = $('#input-username').val();
    let message = $('#input-message').val();
    
    websocket.send(JSON.stringify({
      username: username,
      text: message,
    }));
    
    $('#input-message').val('');
  });
});