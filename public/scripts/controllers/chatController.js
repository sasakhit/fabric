myApp.controller('chatController',
  ['$scope', 'socket',
  function($scope, socket) {

    $scope.label = "名前：";
    $scope.button = "入室";
    $scope.chatlogs = [];
    var isEnter = false;
    var name = "";

    socket.on("server_to_client", function(data){
      $scope.chatlogs.push(data.value);
    });

    $scope.submit = function(message) {
      if (isEnter) {
        message = "[" + name + "]: " + message;
        socket.emit("client_to_server", {value : message});
      } else {
        name = message;
        var entryMessage = name + "さんが入室しました。";
        socket.emit("client_to_server_broadcast", {value : entryMessage});
        socket.emit("client_to_server_personal", {value : name});
        isEnter = true;
        $scope.label = "メッセージ：";
        $scope.button = "送信";
      }
    }
  }
]);
