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
 * Controller for the search box and list of results.
 */
function SearchController($scope, $http) {
  /** The user's query string. */
  $scope.query = '';

  /** Array of all results for the query. */
  $scope.results = [];

  /**
   * Action in response to submitting a search query.
   *
   * Sends a request to the backend and stores the results in results.
   */
  $scope.search = function() {
    $http.get('/_/search', {params: {q: $scope.query}})
        .success(function(data) {
          $scope.results = data;
        });
  };

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
