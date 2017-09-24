myApp.controller('setupController',
  ['$http', '$rootScope', '$scope', '$route', '$log', 'DataService', 'Utils',
  function($http, $rootScope, $scope, $route, $log, DataService, Utils) {

    $scope.username = 'Jim';
    $scope.password = '53a2377961';
    $scope.orgname = 'org1';
    $scope.channelName = 'mychannel';
    $scope.channelConfigPath = '../artifacts/channel/mychannel.tx';
    $scope.chaincodeName_ex = 'example_cc';
    $scope.chaincodePath_ex = 'example_cc';
    $scope.chaincodeVersion_ex = 'v0';
    $scope.functionName_ex = 'init';
    $scope.args_ex = ["a","100","b","200"];
    $scope.chaincodeName_sl = 'chaincode_sl';
    $scope.chaincodePath_sl = 'chaincode_sl';
    $scope.chaincodeVersion_sl = 'v1';
    $scope.functionName_sl = 'n/a';
    $scope.args_sl = ["n/a"];

    $scope.login = function() {
      $scope.comment = 'Logging in ...';
      DataService.login($scope.username, $scope.orgname).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.createChannel = function() {
      $scope.comment = 'Creating channel ...';
      DataService.createChannel($scope.channelName, $scope.channelConfigPath).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.joinChannel = function() {
      $scope.comment = 'Joining channel ...';
      DataService.joinChannel($scope.channelName).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.queryChainInfo = function() {
      DataService.queryChainInfo($scope.channelName)
        .then(function(data) {
          Utils.dialog('Chain Information', data);
        })
        .catch(function(error) {
          alert(error);
          $scope.comment = error;
        });
    }

    $scope.queryBlock = function(blockNumber) {
      DataService.queryBlock($scope.channelName, blockNumber)
        .then(function(data) {
          Utils.dialog('Block Information', data);
          $scope.transactionId = data.data.data[0].payload.header.channel_header.tx_id;
        })
        .catch(function(error) {
          alert(error);
          $scope.comment = error;
        });
    }

    $scope.queryTransaction = function(transactionId) {
      DataService.queryTransaction($scope.channelName, transactionId)
        .then(function(data) {
          //alert(JSON.stringify(data));
          Utils.dialog('Transaction Information', data);
        })
        .catch(function(error) {
          alert(error);
          $scope.comment = error;
        });
    }

    $scope.queryChannel = function() {
      DataService.queryChannel()
        .then(function(data) {
          //alert(JSON.stringify(data));
          Utils.dialog('Channel Information', data);
        })
        .catch(function(error) {
          alert(error);
          $scope.comment = error;
        });
    }

    $scope.queryChaincode = function(type) {
      DataService.queryChaincode(type)
        .then(function(data) {
          //alert(JSON.stringify(data));
          Utils.dialog('Chaincode Information - ' + type, data);
        })
        .catch(function(error) {
          alert(error);
          $scope.comment = error;
        });
    }

    $scope.installChaincode_ex = function() {
      $scope.comment = 'Installing chaincode ...';
      DataService.installChaincode($scope.chaincodeName_ex, $scope.chaincodePath_ex, $scope.chaincodeVersion_ex).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.instantiateChaincode_ex = function() {
      $scope.comment = 'Instantiating chaincode ...';
      DataService.instantiateChaincode($scope.channelName, $scope.chaincodeName_ex, $scope.chaincodeVersion_ex, $scope.functionName_ex, $scope.args_ex).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.installChaincode_sl = function() {
      $scope.comment = 'Installing chaincode ...';
      DataService.installChaincode($scope.chaincodeName_sl, $scope.chaincodePath_sl, $scope.chaincodeVersion_sl).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

    $scope.instantiateChaincode_sl = function() {
      $scope.comment = 'Instantiating chaincode ...';
      DataService.instantiateChaincode($scope.channelName, $scope.chaincodeName_sl, $scope.chaincodeVersion_sl, $scope.functionName_sl, $scope.args_sl).then(function(data) {
        $scope.comment = data;
        Utils.toast(data);
      });
    }

  }
]);
