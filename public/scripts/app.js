var myApp = angular.module('myApp',
  ['ngRoute', 'ngAnimate', 'ui.bootstrap', 'ngTagsInput' , 'ngMaterial', 'ngCookies', 'chart.js']);

myApp.config(['$routeProvider', '$mdThemingProvider', function($routeProvider, $mdThemingProvider) {

  $routeProvider
    .when('/setup', {
      templateUrl: '/views/templates/setup.html',
      controller: 'setupController'
    })
    .when('/example', {
      templateUrl: '/views/templates/example.html',
      controller: 'exampleController'
    })
    .when('/counter', {
      templateUrl: '/views/templates/counter.html',
      controller: 'counterController'
    })
    .when('/fx', {
      templateUrl: '/views/templates/fx.html',
      controller: 'fxController'
    })
    .when('/sl', {
      templateUrl: '/views/templates/sl.html',
      controller: 'slController'
    })
    .when('/chat', {
      templateUrl: '/views/templates/chat.html',
      controller: 'chatController'
    })
    .when('/vote', {
      templateUrl: '/views/templates/vote.html',
      controller: 'voteController'
    })
    .when('/notImplemented', {
      templateUrl: '/views/templates/notImplemented.html'
    })
    .otherwise({
      redirectTo: '/setup'
    });

    $mdThemingProvider.theme('default')
      .primaryPalette('indigo')
      .accentPalette('pink')
      .warnPalette('green')
      .backgroundPalette('grey');
}]);

myApp.controller('navController', ['$scope', '$location', function($scope, $location) {
  $scope.currentNavItem = $location.path().replace('/','');
}]);
