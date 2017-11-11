myApp.controller('voteController',
  ['$scope', 'socket',
  function($scope, socket) {

    $scope.title = "Japan Town Hall Topics"
    /*
    $scope.candidates = [ { id: 1, name: "Topic 1", count: 0 },
                          { id: 2, name: "Topic 2", count: 0 },
                          { id: 3, name: "Topic 3", count: 0 } ]
    */
    $scope.label = "Please select one option, and click VOTE button";
    $scope.button = "Vote";
    $scope.isVoted = false;
    $scope.colors = ['#45B7CD', '#FF6384', '#803690', '#00ADF9', '#DCDCDC', '#46BFBD', '#FDB45C', '#949FB1', '#4D5360'];
    $scope.options = {
    scales: {
        xAxes: [{
            stacked: true
        }],
        yAxes: [{
          barPercentage: 0.9,
          categoryPercentage: 0.8
        }]
    }
};

    socket.on("vote_server_to_client", function(candidates){
      //alert(candidates);
      $scope.candidates = candidates;
      $scope.labels = _.map(candidates, function(candidate) { return candidate.id });
      $scope.values = _.map(candidates, function(candidate) { return candidate.count });
    });

    /*
    socket.on("vote_server_to_client", function(data){
      $scope.candidates.forEach(function(c) {if (c.id == data.id) c.count++ ;})
    });
    */

    $scope.submit = function(selectedItem) {
      if ($scope.isVoted) {
        // Do nothing
      } else {
        //alert(selectedItem);
        socket.emit("vote_client_to_server", selectedItem);
        $scope.isVoted = true;
        $scope.label = "Thank you for voting!";
        $scope.button = "Already Voted";
      }
    }
  }
]);
