var mciServices = mciServices || {};

mciServices.rest = angular.module('mciServices.rest', []);

mciServices.rest.factory('mciBaseRestService', ['$http', function($http) {
    // private vars
    var baseUrl = '';
    var resources = {};

    // the service that will be returned
    var service = {};

    var httpCall = function(method, resource, idents, config, callbacks) {
        if (!$http[method] || typeof $http[method] !== 'function') {
            alert('invalid http method: ' + method);
            return;
        }

        config.method = method;
        config.url = [baseUrl, resource].concat(idents).join('/');

        $http(config)
            .success(callbacks.success || function() {})
            .error(callbacks.error || function() {});
    };

    ['get', 'post', 'put', 'delete'].forEach(function(method) {
        service[method+'Resource'] = function(resource, idents, config, callbacks) {
            httpCall(method, resource, idents, config, callbacks);
        };
    });

    return service;
}]);

mciServices.rest.factory('mciTaskDrawerRestService', 
    [
      'mciBaseRestService', 
      function(baseSvc) {
        var resource = 'task_history_json';
        var service = {};
        var defaultRadius = 10;

        service.fetchHistory = function(taskId, historyType, radius, callbacks) {
          var config = { params: { radius: radius || defaultRadius } };
          baseSvc.getResource(resource, [taskId, historyType], config, callbacks);
        };

        return service;
      }
    ]
);

mciServices.rest.factory('mciTasksRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'tasks';

    var service = {};

    service.getResource = function() {
        return resource;
    };

    service.takeActionOnTask = function(taskId, action, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        baseSvc.putResource(resource, [taskId], config, callbacks);
    };

    service.getTask = function(taskId, callbacks) {
        baseSvc.getResource(resource, [taskId], {}, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciBuildsRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'builds';

    var service = {};

    service.takeActionOnBuild = function(buildId, action, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        baseSvc.putResource(resource, [buildId], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciHostRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'host';

    var service = {};

    service.updateStatus = function(hostId, action, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        baseSvc.putResource(resource, [hostId], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciHostsRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'hosts';

    var service = {};

    service.updateStatus = function(hostIds, action, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        config.data['host_ids'] = hostIds;
        baseSvc.putResource(resource, [], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciVersionsRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'version';

    var service = {};

    service.takeActionOnVersion = function(versionId, action, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        baseSvc.putResource(resource, [versionId], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciBuildVariantHistoryRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'build_variant';

    var service = {};

    service.getBuildVariantHistory = function(project, buildVariant, params, callbacks) {
        var _project = encodeURIComponent(project);
        var _buildVariant = encodeURIComponent(buildVariant);

        var config = { params: params };
        baseSvc.getResource(resource, [_project, _buildVariant], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciTaskHistoryRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'task_history';

    var service = {};

    service.getTaskHistory = function(project, taskName, params, callbacks) {
        var _project = encodeURIComponent(project);
        var _taskName = encodeURIComponent(taskName);

        var config = { params: params };
        baseSvc.getResource(resource, [_project, _taskName], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciLoginRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'login';

    var service = {};

    service.authenticate = function(username, password, data, callbacks) {
        var config = { data: data };
        config.data['username'] = username;
        config.data['password'] = password;
        baseSvc.postResource(resource, [], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciSpawnRestService', ['mciBaseRestService', function(baseSvc) {
    var resource = 'spawn';

    var service = {};

    service.getSpawnedHosts = function(action, params, callbacks) {
        baseSvc.getResource(resource, action, params, callbacks);
    }

    service.getSpawnableDistros = function(action, params, callbacks) {
        var config = { params: params };
        baseSvc.getResource(resource, action, config, callbacks);
    };

    service.getUserKeys = function(action, params, callbacks) {
        var config = { params: params };
        baseSvc.getResource(resource, action, config, callbacks);
    };

    service.spawnHost = function(spawnInfo, data, callbacks) {
        var config = { data: data };
        config.data['distro'] = spawnInfo.distroName;
        config.data['save_key'] = spawnInfo.saveKey;
        config.data['key_name'] = spawnInfo.spawnKey.name;
        config.data['public_key'] = spawnInfo.spawnKey.key;
        config.data['userdata'] = spawnInfo.userData;
        baseSvc.putResource(resource, [], config, callbacks);
    };

    service.terminateHost = function(action, hostId, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        config.data['host_id'] = hostId;
        baseSvc.postResource(resource, [], config, callbacks);
    };

    service.updateRDPPassword = function(action, hostId, rdpPassword, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        config.data['host_id'] = hostId;
        config.data['rdp_pwd'] = rdpPassword;
        baseSvc.postResource(resource, [], config, callbacks);
    };

    service.extendHostExpiration = function(action, hostId, addHours, data, callbacks) {
        var config = { data: data };
        config.data['action'] = action;
        config.data['host_id'] = hostId;
        config.data['add_hours'] = addHours;
        baseSvc.postResource(resource, [], config, callbacks);
    };

    return service;
}]);

mciServices.rest.factory('mciTaskStatisticsRestService', ['mciBaseRestService', function(baseSvc) {
  var resource = 'task_stats';
  var service = {};

  service.getTimeStatistics = function getTimeStatistics(field1, field2, groupByField, days, callbacks) {
    baseSvc.getResource(resource, [field1, field2, groupByField, days], {}, callbacks);
  };

  return service;
}]);
