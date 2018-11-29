let Server = function () {
};

Server.prototype.sendMessage = function (name, payload, callback) {
    document.addEventListener('astilectron-ready', function () {
        // This will send a message to Go
        astilectron.sendMessage({name: name, payload: payload}, function (message) {
            callback(message)
        });
    });
};

(function () {
    let system = {
        currentGroupName: ''
    };
    let server = new Server();
    let app = new Vue({
        el: '#app',
        data: {
            systemHosts: [],//[{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true}],
            systemHostsLoading: false,
            fullscreenLoading: false,
            hostGroups: [],
            addIp: '',
            addHost: '',
            system: system
        },
        methods: {
            groupSwitch: function (value) {
                if (value) {
                    this.$message({
                        message: '启用分组成功',
                        type: 'success'
                    });
                } else {
                    this.$message({
                        message: '关闭分组成功',
                        type: 'info'
                    });
                }
            },
            initSystemHosts: function () {
                let _this = this;
                // this.$notify({
                //     title: '提示',
                //     message: '左侧分组可通过点击复选框快速启用及关闭哦~',
                //     position: 'top-left'
                // });
                Promise.all([this.getList("System Hosts"), this.getHostGroups()]).then((results) => {
                    console.log(_this.hostGroups);
                    _this.hostGroups[0].active = true;
                    _this.system.currentGroupName = _this.hostGroups[0].name;
                    _this.fullscreenLoading = false;
                    _this.systemHostsLoading = false;
                    console.log(results);
                });
            },
            changeGroup: function (groupName) {
                for (let i in this.hostGroups) {
                    let item = this.hostGroups[i];
                    item.active = groupName === item.name;
                }
            },
            getList: function (type) {
                let _this = this;
                return new Promise((resolve) => {
                    server.sendMessage("list", {type: type}, (message) => {
                        _this.systemHosts = message.payload;
                        _this.systemHostsLoading = false;
                        resolve(true)
                    });
                })
            },
            getHostGroups: function () {
                let _this = this;
                return new Promise((resolve) => {
                    server.sendMessage('groups', {}, (message) => {
                        _this.hostGroups = message.payload;
                        _this.hostGroups.unshift({
                            name: "System Hosts",
                            active: false,
                        });
                        resolve(true)
                        setTimeout(() => {
                            _this.fullscreenLoading = false;
                        }, 1000)
                    })
                })
            }
        }
    });
    app.initSystemHosts();
})();