myApp.factory('DataService',
    ['$http', '$mdDialog', '$q', '$route',
    function($http, $mdDialog, $q, $route) {

      return {
        login : login,
        queryChainInfo : queryChainInfo,
        queryBlock : queryBlock,
        queryTransaction : queryTransaction,
        queryChannel : queryChannel,
        queryChaincode : queryChaincode,
        createChannel : createChannel,
        joinChannel : joinChannel,
        queryChainInfo : queryChainInfo,
        installChaincode : installChaincode,
        instantiateChaincode : instantiateChaincode,
        deployHelloworld : deployHelloworld,
        queryHelloworld : queryHelloworld,
        invokeHelloworld : invokeHelloworld,
        deployCounter : deployCounter,
        queryCounter : queryCounter,
        invokeCounter : invokeCounter,
        deployFx : deployFx,
        queryFx : queryFx,
        invokeFx : invokeFx,
        deploySl : deploySl,
        querySl : querySl,
        invokeSlTrade : invokeSlTrade,
        invokeSlCommon : invokeSlCommon
      };

      function login(username, orgname) {
        return $http.post('/users', {username: username, orgname: orgname}).then(function(response) {
          return response.data;
        });
      }

      function queryChainInfo(channelName) {
        return $http.get('/channels/' + channelName, {params: {peer: 'peer1'}}).then(function(response) {
          return response.data;
        });
      }

      function queryBlock(channelName, blockNumber) {
        return $http.get('/channels/' + channelName + '/blocks/' + blockNumber, {params: {peer: 'peer1'}}).then(function(response) {
          return response.data;
        });
      }

      function queryTransaction(channelName, transactionId) {
        return $http.get('/channels/' + channelName + '/transactions/' + transactionId, {params: {peer: 'peer1'}}).then(function(response) {
          return response.data;
        });
      }

      function queryChannel() {
        return $http.get('/channels', {params: {peer: 'peer1'}}).then(function(response) {
          return response.data;
        });
      }

      function queryChaincode(type) {
        return $http.get('/chaincodes', {params: {peer: 'peer1', type: type}}).then(function(response) {
          return response.data;
        });
      }

      function createChannel(channelName, channelConfigPath) {
        return $http.post('/channels', {channelName: channelName, channelConfigPath: channelConfigPath}).then(function(response) {
          return response.data;
        });
      }

      function joinChannel(channelName) {
        return $http.post('/channels/' + channelName + '/peers').then(function(response) {
          return response.data;
        });
      }

      function installChaincode(chaincodeName, chaincodePath, chaincodeVersion) {
        return $http.post('/chaincodes', {chaincodeName: chaincodeName, chaincodePath: chaincodePath, chaincodeVersion: chaincodeVersion}).then(function(response) {
          return response.data;
        });
      }

      function instantiateChaincode(channelName, chaincodeName, chaincodeVersion, functionName, args) {
        return $http.post('/channels/' + channelName + '/chaincodes', {chaincodeName: chaincodeName, chaincodeVersion: chaincodeVersion, functionName: functionName, args: args}).then(function(response) {
          return response.data;
        });
      }

      function deployHelloworld() {
        return $http.get('/helloworld/deploy', {timeout: 10000}).then(function(response) {
          return response.data;
        });
      }

      function queryHelloworld(owner, channelName, chaincodeName) {
        return $http.get('/channels/' + channelName + '/chaincodes/' + chaincodeName, {params: {peer: 'peer1', fcn: 'query', args: '["' + owner + '"]'}}).then(function(response) {
          return response.data;
        });
      }

      function invokeHelloworld(fromOwner, toOwner, moveQuantity, channelName, chaincodeName) {
        //return $http.post('/channels/' + channelName + '/chaincodes/' + chaincodeName, {fcn: 'move', args: '["' + fromOwner + '","' + toOwner + '","' + moveQuantity + '"]'}).then(function(response) {
        return $http.post('/channels/' + channelName + '/chaincodes/' + chaincodeName, {fcn: 'move', args: [fromOwner,toOwner,moveQuantity]}).then(function(response) {
          return response.data;
        });
      }

      function deployCounter() {
        return $http.get('/counter/deploy', {timeout: 10000}).then(function(response) {
          return response.data;
        });
      }

      function queryCounter(chaincodeID) {
        return $http.get('/counter/query', {params: {chaincodeID: chaincodeID}}).then(function(response) {
          return response.data;
        }).catch(function(error) {
          throw error.data;
        });
      }

      function invokeCounter(counterId, chaincodeID) {
        return $http.post('/counter/invoke', {counterId: counterId, chaincodeID: chaincodeID}).then(function(response) {
          return response.data;
        });
      }

      function deployFx() {
        return $http.get('/fx/deploy', {timeout: 10000}).then(function(response) {
          return response.data;
        });
      }

      function queryFx(chaincodeID) {
        return $http.get('/fx/query', {params: {chaincodeID: chaincodeID}}).then(function(response) {
          return response.data;
        }).catch(function(error) {
          throw error.data;
        });
      }

      function invokeFx(fromAccount, fromCcy, fromAmt, toAccount, toCcy, toAmt, chaincodeID) {
        return $http.post('/fx/invoke', {fromAccount: fromAccount, fromCcy: fromCcy, fromAmt: fromAmt, toAccount: toAccount, toCcy: toCcy, toAmt: toAmt, chaincodeID: chaincodeID}).then(function(response) {
          return response.data;
        });
      }

      function deploySl() {
        return $http.get('/sl/deploy', {timeout: 10000}).then(function(response) {
          return response.data;
        });
      }

      function querySl(channelName, chaincodeName, functionName) {
        return $http.get('/channels/' + channelName + '/chaincodes/' + chaincodeName, {params: {peer: 'peer1', fcn: functionName, args: '["n/a"]'}}).then(function(response) {
          return response.data;
        });
      }

      function invokeSlTrade(brInd, borrower, lender, tradeDate, settleDate, secCode, qty, ccy, amt, channelName, chaincodeName, functionName) {
        var trade = {brInd: brInd, borrower: borrower, lender: lender, tradeDate: tradeDate, settleDate: settleDate, secCode: secCode, qty: qty, ccy: ccy, amt: amt};
        var args = [JSON.stringify(trade)];
        return $http.post('/channels/' + channelName + '/chaincodes/' + chaincodeName, {fcn: functionName, args: args}).then(function(response) {
          return response.data;
        });
      }

      function invokeSlCommon(channelName, chaincodeName, functionName) {
        return $http.post('/channels/' + channelName + '/chaincodes/' + chaincodeName, {fcn: functionName, args: []}).then(function(response) {
          return response.data;
        });
      }

    }
]);
