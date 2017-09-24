module.exports = function(io) {

  // S04. connectionイベント・データを受信する
  io.sockets.on('connection', function(socket) {
    //io.of('/chat').on('connection', function(socket) {
      var name;
      // S05. client_to_serverイベント・データを受信する
      socket.on('client_to_server', function(data) {
          // S06. server_to_clientイベント・データを送信する
          io.sockets.emit('server_to_client', {value : data.value});
      });
      // S07. client_to_server_broadcastイベント・データを受信し、送信元以外に送信する
      socket.on('client_to_server_broadcast', function(data) {
          socket.broadcast.emit('server_to_client', {value : data.value});
      });
      // S08. client_to_server_personalイベント・データを受信し、送信元のみに送信する
      socket.on('client_to_server_personal', function(data) {
          var id = socket.id;
          name = data.value;
          var personalMessage = "あなたは、" + name + "さんとして入室しました。"
          io.to(id).emit('server_to_client', {value : personalMessage});
      });
      // S09. disconnectイベントを受信し、退出メッセージを送信する
      socket.on('disconnect', function() {
          if (name == 'undefined') {
              console.log("未入室のまま、どこかへ去っていきました。");
          } else {
              var endMessage = name + "さんが退出しました。"
              io.sockets.emit('server_to_client', {value : endMessage});
          }
      });
  });
};
