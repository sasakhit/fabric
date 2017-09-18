myApp.controller('slController',
  ['$http', '$rootScope', '$scope', '$route', '$log', '$mdDialog', '$filter', 'DataService', 'socket',
  function($http, $rootScope, $scope, $route, $log, $mdDialog, $filter, DataService, socket) {

    var user;
    var now = new Date();
    var today = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 0, 0, 0);

    $scope.username = 'Jim';
    $scope.orgname = 'org1';
    $scope.channelName = 'mychannel';
    $scope.chaincodeName = 'chaincode_sl';

    $scope.brInd = 'B';
    $scope.borrower = '0876111111';
    $scope.lender = '0876222222';
    $scope.tradeDate = today;
    $scope.settleDate = today;
    $scope.secCode = '8756';
    $scope.qty = 5000;
    $scope.ccy = 'JPY';
    $scope.amt = 4000000;

    socket.on("block_change", function(data){
      $scope.blockheight = data.height.low;
    });

    $scope.queryChainInfo = function() {
      DataService.queryChainInfo($scope.channelName)
        .then(function(data) {
          alert('Block Height: ' + data.height.low);
          //$scope.comment = data;
        })
        .catch(function(error) {
          alert(error)
          $scope.comment = error;
        });
    }

    $scope.getOutstandings = function() {
      $scope.comment = 'Getting outstandings ...';
      getOutstandings();
    }

    $scope.getTransactions = function() {
      $scope.comment = 'Getting transactions ...';
      getTransactions();
    }

    $scope.tradeSl = function(ev) {
      function dialogController($scope, $mdDialog, brInd, borrower, lender, tradeDate, settleDate, secCode, qty, ccy, amt, channelName, chaincodeName) {
        $scope.brInd = brInd;
        $scope.borrower = borrower;
        $scope.lender = lender;
        $scope.tradeDate = tradeDate;
        $scope.settleDate = settleDate;
        $scope.secCode = secCode;
        $scope.qty = qty;
        $scope.ccy = ccy;
        $scope.amt = amt;

        $scope.ok = function(brInd, borrower, lender, tradeDate, settleDate, secCode, qty, ccy, amt) {
          if (!brInd || !borrower || !lender || !tradeDate || !settleDate || !secCode || !qty || !ccy || !amt) {
            alert('All the fields should be filled in');
          }
          else {
            $scope.comment = 'Trading SL ...';
            bookSlTrade(brInd, borrower, lender, tradeDate, settleDate, secCode, qty, ccy, amt)
              .then(function() {
                $scope.comment = 'Offsetting Outstandings ...';
                return offsetOutstandings();
              })
              .then(function() {
                $scope.comment = 'Revaluating ...';
                return revaluateMtm();
              })
              .then(function() {
                $scope.comment = 'Done!';
                getOutstandings();
                getTransactions();
                alert("Trade is recorded in Blockchain");
                $mdDialog.hide();
              })
              .catch(function(error) {
                $scope.comment = error;
                alert(error);
              });
          }
        }

        $scope.cancel = function() {
          $mdDialog.hide();
        }
      }

      $mdDialog.show({
        controller: dialogController,
        targetEvent: ev,
        ariaLabel:  'SL Trade Entry',
        clickOutsideToClose: true,
        templateUrl: 'views/templates/slTradeEntry.html',
        onComplete: afterShowAnimation,
        size: 'large',
        bindToController: true,
        autoWrap: false,
        parent: angular.element(document.body),
        preserveScope: true,
        locals: {
          brInd: $scope.brInd,
          borrower: $scope.borrower,
          lender: $scope.lender,
          tradeDate: $scope.tradeDate,
          settleDate: $scope.settleDate,
          secCode: $scope.secCode,
          qty: $scope.qty,
          ccy: $scope.ccy,
          amt: $scope.amt,
          channelName: $scope.channelName,
          chaincodeName: $scope.chaincodeName
        }
      });

      function afterShowAnimation(scope, element, options) {
        // post-show code here: DOM element focus, etc.
      }
    }

    $scope.marginCall = function() {
      $scope.comment = 'Calculating Margin Calls ...';
      revaluateMtm()
        .then(function() {
          return calcMarginCall();
        })
        .then(function() {
          return offsetOutstandings();
        })
        .then(function() {
          $scope.comment = 'Done';
          alert('Margin Calls are recorded in Blockchain');
          getOutstandings();
          getTransactions();
        })
        .catch(function(error) {
          $scope.comment = error;
        });
    }

    function bookSlTrade(brInd, borrower, lender, tradeDate, settleDate, secCode, qty, ccy, amt) {
      $scope.comment = 'Booking SL Trade ...';
      var tradeDateMs = tradeDate.getTime().toString();
      var settleDateMs = settleDate.getTime().toString();
      var promise = DataService.invokeSlTrade(brInd, borrower, lender, tradeDateMs, settleDateMs, secCode, qty, ccy, amt, $scope.channelName, $scope.chaincodeName, 'tradeSl')
        .then(function(data) {
          $scope.comment = data;
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

    function calcMarginCall() {
      $scope.comment = 'Calculating Margin Calls ...';
      var promise = DataService.invokeSlCommon($scope.channelName, $scope.chaincodeName, 'calcMarginCall')
        .then(function(data) {
          $scope.comment = data;
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

    function offsetOutstandings() {
      $scope.comment = 'Offsetting Outstandings ...';
      var promise = DataService.invokeSlCommon($scope.channelName, $scope.chaincodeName, 'offsetOutstandings')
        .then(function(data) {
          $scope.comment = data;
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

    function revaluateMtm() {
      $scope.comment = 'Revaluating ...';
      var promise = DataService.invokeSlCommon($scope.channelName, $scope.chaincodeName, 'revaluateMtm')
        .then(function(data) {
          $scope.comment = data;
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

    function getOutstandings() {
      var promise = DataService.querySl($scope.channelName, $scope.chaincodeName, 'getOutstandings')
        .then(function(data) {
          $scope.outstandings = data;
          $scope.comment = '';
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

    function getTransactions() {
      var promise = DataService.querySl($scope.channelName, $scope.chaincodeName, 'getTransactions')
        .then(function(data) {
          data.tradeDate = new Date(data.tradeDate);
          data.settleDate = new Date(data.settleDate);
          $scope.transactions = data;
          $scope.comment = '';
        })
        .catch(function(error) {
          $scope.comment = error;
        });

      return promise;
    }

  }
]);
