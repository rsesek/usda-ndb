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

angular.module('foodle', [])
    .config(function($routeProvider) {
      $routeProvider
          .when('/search', {templateUrl: '/partials/search.html'})
          .when('/food/:NDBID', {templateUrl: '/partials/detail.html'});
    })
    .service('FoodGroups', function($http) {
      var service = {
        _data: {},
        nameFromId: function(id) {
          if (id in this._data)
            return this._data[id].Description;
          return 'Unknown';
        }
      };
      $http.get('/_/foodGroups').success(function(data) {
        for (var i = 0; i < data.length; ++i) {
          service._data[data[i].GroupCode] = data[i];
        }
      });
      return service;
    })
    .filter('foodGroupName', function(FoodGroups) {
      return function(input) {
        return FoodGroups.nameFromId(input);
      };
    })
    .service('NutrientDefinitions', function($http) {
      var service = {
        _data: {},
        nameFromId: function(id) {
          if (id in this._data)
            return this._data[id].Description;
          return 'Unknown';
        },
        unitsForId: function(id) {
          if (id in this._data)
            return this._data[id].Units;
          return '';
        },
        sortFoodNutrients: function(nutrients) {
          nutrients.sort(function(a, b) {
            var nA = service._data[a.NutrientID];
            var nB = service._data[b.NutrientID];
            if (!nA || !nB)
              return -1;
            return nA.SortOrder - nB.SortOrder;
          });
        },
      };
      $http.get('/_/nutrients').success(function(data) {
        for (var i = 0; i < data.length; ++i) {
          service._data[data[i].NutrientID] = data[i];
        }
      });
      return service;
    })
    .filter('nutrientName', function(NutrientDefinitions) {
      return function(input) {
        return NutrientDefinitions.nameFromId(input);
      };
    })
    .filter('nutrientUnits', function(NutrientDefinitions) {
      return function(input) {
        return NutrientDefinitions.unitsForId(input);
      };
    });
