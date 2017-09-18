myApp.controller('fxController',
  ['$http', '$rootScope', '$scope', '$route', '$log', '$mdDialog', 'DataService',
  function($http, $rootScope, $scope, $route, $log, $mdDialog, DataService) {

    var user;

    $scope.chaincodeID = '01451afb7bc9a1848f0e94cc2d33f2cce6be13d6ac1c4fdbf963c8c62d54e17d';
    $scope.username = 'user_type1_0';
    $scope.password = '53a2377961';
    $scope.fromAccount = '0895123456';
    $scope.fromCcy = 'JPY';
    $scope.fromAmt = 5000;
    $scope.toAccount = '0895999999';
    $scope.toCcy = 'USD';
    $scope.toAmt = 50;

    $scope.login = function() {
      $scope.comment = 'Logging in ...';
      DataService.login($scope.username, $scope.password).then(function(data) {
        $scope.comment = data;
      });
    }

    $scope.deploy = function() {
      $scope.comment = 'Deploying chaincode ...';
      DataService.deployFx().then(function(data) {
        $scope.chaincodeID = data;
      });
    }

    $scope.query = function() {
      $scope.comment = 'Querying ...';
      DataService.queryFx($scope.chaincodeID)
        .then(function(data) {
          $scope.accounts = data;
          $scope.comment = '';
        })
        .catch(function(error) {
          $scope.comment = error;
        });
    }

    $scope.invoke = function(fromAccount,fromCcy,fromAmt,toAccount,toCcy,toAmt) {
      $scope.comment = 'Invoking ...';
      DataService.invokeFx(fromAccount,fromCcy,fromAmt,toAccount,toCcy,toAmt, $scope.chaincodeID)
        .then(function(data) {
          $scope.comment = data;
        })
        .catch(function(error) {
          $scope.comment = error;
        });
    }

  }
]);
