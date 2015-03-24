var app = angular.module("app", ["ngRoute", "ngCookies", "ngMaterial", "smart-table"]);

app.config(["$routeProvider", function($routeProvider) {
    $routeProvider
        .when("/login", {
            templateUrl: "views/login.html",
            controller: "loginController",
        }).when("/form", {
            templateUrl: "views/form.html",
            controller: "formController",
        }).when("/admin/login", {
            templateUrl: "views/login.html",
            controller: "loginController",
        }).when("/admin/list", {
            templateUrl: "views/list.html",
            controller: "listController",
        }).when("/done", {
            templateUrl: "views/done.html",
            controller: "doneController",
        }).otherwise({ redirectTo: "/login" });
}]);

app.config(["$mdThemingProvider", function($mdThemingProvider) {
    $mdThemingProvider.theme("default")
        .primaryPalette("blue");
}]);

app.factory("session", ["$cookies", function($cookies) {
    return { 
        setID: function(id) {
            $cookies.sessionID = id;
        },
        getID: function() {
            return $cookies.sessionID || "";
        },
        deleteID : function() {
            delete $cookies.sessionID;
        }
    };
}]);

app.factory("alert", function() {
    return {
        hidden: true,
        message: "",
    };
});

app.controller("loginController", ["$scope", "$http", "$location", "session", "alert", function($scope, $http, $location, session, alert) {
    $scope.submit = function(login) {
        $scope.alert.hidden = true;

        var url = $scope.admin ? "api/1.0/admin/auth" : "api/1.0/auth";

        $http({
            method: "POST",
            url: url,
            data: login,
            headers: {
                "Accept": "application/json",
            },
        }).success(function(data, status) {
            if (status != 200) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Login error: ", status, data);
                return;
            }
            if (data.SessionID == null || data.SessionID == "") {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Login error: ", status, data);
                return;
            }
            if (data.Completed && !$scope.admin) {
                $location.path("/done");
                return;
            }
            session.setID(data.SessionID);
            if ($scope.admin) {
                $location.path("/admin/list");
                return;
            } else {
                $location.path("/form");
            }

        }).error(function(data, status) {
            $scope.alert.hidden = false;
            if (status == 401) {
                $scope.alert.message = "Bad username or password";
            } else {
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
            }
            console.log("Login error: ", status, data);
        });
    };

    $scope.login = {};

    $scope.admin = ($location.path().indexOf("admin") > -1);

    $scope.alert = alert;

    if (session.getID() != "") {
        if ($scope.admin) {
            $location.path("/admin/list");
            return;
        } else {
            $location.path("/form");
        }
    }
}]);

app.controller("formController", ["$scope", "$http", "$location", "$window", "session", "alert", function($scope, $http, $location, $window, session, alert) {

    $scope.logout = function(expired) {
        if (expired) {
            $scope.alert.hidden = false;
            $scope.alert.message = "Your session expired. Please log in again";
        }
        session.deleteID();
        $location.path("/login");
    };

    $scope.open_handbook = function() {
        $scope.data.clicked = true;
        $window.open("images/handbook.pdf", "BISD Handbook", "toolbar=no");
    };

    $scope.submit = function(data) {
        $scope.alert.hidden = true;

        $http({
            method: "POST",
            url: "api/1.0/submit",
            data: {Campus: data.Campus, Agree: data.Agree},
            headers: {
                "Accept": "application/json",
                "X-Session-Key": $scope.sessionID,
            },
        }).success(function(data, status) {
            if (status != 200) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Submit error: ", status, data);
                return;
            }
            if (data.Status == null || data.Status != true) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Submit error: ", status, data);
                return;
            }
            $location.path("/done");

        }).error(function(data, status) {
            if (status == 401) {
                $scope.logout(true);
            } else {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
            }
            console.log("Submit error: ", status, data);
        });
    };

    // setup data
    $scope.sessionID = session.getID();

    $scope.alert = alert;
    $scope.alert.hidden = true;

    $scope.data = {
        clicked: false,
    };

    // check for login
    if ($scope.sessionID == "") {
        $scope.logout();
    }
}]);

app.controller("listController", ["$scope", "$http", "$location", "session", "alert", function($scope, $http, $location, session, alert) {

    $scope.logout = function(expired) {
        if (expired) {
            $scope.alert.hidden = false;
            $scope.alert.message = "Your session expired. Please log in again";
        }
        session.deleteID();
        $location.path("/admin/login");
    };

    $scope.fetch = function() {
        $scope.alert.hidden = true;

        $http({
            method: "GET",
            url: "api/1.0/admin/list",
            headers: {
                "Accept": "application/json",
                "X-Session-Key": $scope.sessionID,
            },
        }).success(function(data, status) {
            if (status != 200) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Fetch error: ", status, data);
                return;
            }
            if (data.List == null) {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                console.log("Fetch error: ", status, data);
                return;
            }
            $scope.ajaxList = data.List;
            angular.forEach($scope.ajaxList, function(val,key) {
                $scope.ajaxList[key].Time = new Date(val.Time); 
                delete $scope.ajaxList[key].Username;
                delete $scope.ajaxList[key].Headers;
            });
            $scope.displayList = [].concat($scope.ajaxList);

        }).error(function(data, status) {
            if (status == 401) {
                $scope.logout(true);
            } else {
                $scope.alert.hidden = false;
                $scope.alert.message = "Something bad happened: " + angular.toJson(data);
                $scope.logout();
            }
            console.log("Fetch error: ", status, data);
        });
    };

    // setup data
    $scope.sessionID = session.getID();

    $scope.alert = alert;
    $scope.alert.hidden = true;

    $scope.displayList = [];
    $scope.ajaxList = [];

    $scope.filter = {
        search: "",
    };

    // check for login
    if ($scope.sessionID == "") {
        $scope.logout();
        return;
    }

    $scope.fetch();
}]);

app.controller("doneController", ["$scope", "$location", "$window", "session", function($scope, $location, $window, session) {
    $scope.open_handbook = function() {
        $window.open("images/handbook.pdf", "BISD Handbook", "toolbar=no");
    };

    $scope.logout = function() {
        session.deleteID();
        $location.path("/login");
    };
}]);
