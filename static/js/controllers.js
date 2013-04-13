//
// USDA-NDB Viewer
// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

/**
 * Controller for the home page.
 */
function HomeController($scope, $location) {
  /** The user's query string. */
  $scope.query = $location.search().q;

  /**
   * Action for the search button that redirects to the search result list.
   */
  $scope.search = function() {
    $location.search('q', $scope.query);
    $location.path('/search');
  };
}

/**
 * Controller for the search box and list of results.
 */
function SearchController($scope, $http) {
  /** Array of all results for the query. */
  $scope.results = [];

  // $scope.query gets inherited from the parent scope.
  $http.get('/_/search', {params: {q: $scope.query}})
      .success(function(data) {
        $scope.results = data;
      });

  /**
   * Filters the result list to just 10 items or fewer.
   */
  $scope.filteredResults = function() {
    var results = [];
    for (var i = 0; i < 10 && i < $scope.results.length; ++i) {
      results.push($scope.results[i]);
    }
    return results;
  };
}

/**
 * Controller for the food detail page.
 */
function DetailController($scope, $routeParams, $http) {
  /** The food object. */
  $scope.food = {};

  $http.get('/_/food/' + $routeParams.NDBID)
      .success(function(data) {
        $scope.food = data;
        $scope.unit = $scope.food.Weights[0];
        $scope.onUnitsChanged();
      })
      .error(function(data) {
        $scope.error = data;
      });

  /**
   * Called when the units selector changes.
   */
  $scope.onUnitsChanged = function() {
    $scope.unitAmount = $scope.unit.Amount;
  };
}
