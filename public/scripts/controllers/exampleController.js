myApp.controller('exampleController',
  ['$http', '$rootScope', '$scope', '$route', '$log', '$mdDialog', 'DataService',
  function($http, $rootScope, $scope, $route, $log, $mdDialog, DataService) {

    $scope.username = 'Jim';
    $scope.orgname = 'org1';
    $scope.channelName = 'mychannel';
    $scope.chaincodeName = 'example_cc';

    $scope.query = function(owner) {
      $scope.comment = 'Querying ...';
      DataService.queryHelloworld(owner, $scope.channelName, $scope.chaincodeName)
        .then(function(data) {
          $scope.owner = owner;
          $scope.quantity = data;
          $scope.comment = '';
        })
        .catch(function(error) {
          $scope.comment = error.data;
        });
    }

    $scope.invoke = function(fromOwner, toOwner, moveQuantity) {
      $scope.comment = 'Invoking ...';
      DataService.invokeHelloworld(fromOwner, toOwner, moveQuantity, $scope.channelName, $scope.chaincodeName)
        .then(function(data) {
          $scope.comment = 'Invoking ...';
        })
        .catch(function(error) {
          $scope.comment = error.data;
        });
    }

  }
]);
