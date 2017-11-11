module.exports = function(io) {

  var candidates = [ { id: 1, name: "Topic 1", count: 0 },
                     { id: 2, name: "Topic 2", count: 0 },
                     { id: 3, name: "Topic 3", count: 0 } ]

  io.sockets.on('connection', function(socket) {
    var id = socket.id;
    io.to(id).emit('vote_server_to_client', candidates);

    socket.on('vote_client_to_server', function(id) {
      candidates.forEach(function(c) {if (c.id == id) c.count++ ;})
      io.sockets.emit('vote_server_to_client', candidates);
    });
  });
};
