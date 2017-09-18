myApp.factory('Utils',
      ['$rootScope', '$mdDialog', '$mdToast', function($rootScope, $mdDialog, $mdToast) {

    return {
      toast : toast,
      dialog : dialog
    };

    function toast(data) {
      $mdToast.show(
        $mdToast.simple()
        .textContent(data)
        .position("top right")
        .hideDelay(3000)
      );
    }

    function dialog(title, data) {
      $mdDialog.show(
        $mdDialog.alert()
        .parent(angular.element(document.querySelector('#popupContainer')))
        .clickOutsideToClose(true)
        .title(title)
        .textContent(data)
        .ariaLabel('Dialog')
        .ok('OK')
      );
    }
  }
]);
