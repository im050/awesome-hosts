let Server = function () {
};

Server.prototype.sendMessage = function (name, payload, callback) {
    document.addEventListener('astilectron-ready', function () {
        // This will send a message to Go
        astilectron.sendMessage({name: name, payload: payload}, function (message) {
            console.log(message)
            callback(message)
        });
    });
};

(function () {
    let system = {
        currentGroupName: '',
        currentHosts: [],
        systemHosts: []
    };
    let server = new Server();
    let app = new Vue({
        el: '#app',
        data: {
            hostsLoading: true,
            fullscreenLoading: true,
            hostGroups: [],
            addIp: '',
            addHost: '',
            system: system,
            ipPrepareList: []
        },
        methods: {
            querySearch(queryString, show) {
                let ipPrepareList = this.ipPrepareList;
                let results = queryString ? ipPrepareList.filter(this.createFilter(queryString)) : ipPrepareList;
                // 调用 callback 返回建议列表的数据
                show(results);
            },
            createFilter(queryString) {
                return (ipPrepareList) => {
                    return (ipPrepareList.value.indexOf(queryString) === 0);
                };
            },
            handleSelect(item) {
                console.log(item);
            },
            loadIpPrepareList() {
                server.sendMessage("intranet", {}, (message) => {
                    var data = [{value: "127.0.0.1"}];
                    let ip = message.payload;
                    if (ip !== "") {
                        data.push({value: ip})
                    }
                    this.ipPrepareList = data;
                })
            },
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
            init: function () {
                let _this = this;
                Promise.all([this.getList("System Hosts"), this.getHostGroups()]).then((results) => {
                    _this.hostGroups[0].active = true;
                    _this.system.currentGroupName = _this.hostGroups[0].name;
                    setTimeout(() => function () {
                        _this.fullscreenLoading = false;
                        _this.hostsLoading = false;
                    }, 1000);
                    console.log(results);
                });
            },
            changeGroup: function (groupName) {
                for (let i in this.hostGroups) {
                    let item = this.hostGroups[i];
                    if (groupName === item.name) {
                        item.active = true;
                        this.system.currentHosts = (groupName === "System Hosts") ? this.system.systemHosts : item.hosts;
                        this.system.currentGroupName = item.name;
                    } else {
                        item.active = false;
                    }
                }
            },
            getList: function (type) {
                let _this = this;
                return new Promise((resolve) => {
                    server.sendMessage("list", {type: type}, (message) => {
                        _this.system.currentHosts = message.payload;
                        _this.system.systemHosts = message.payload;
                        _this.hostsLoading = false;
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
                        resolve(true);
                    })
                })
            }
        },
        mounted() {
            this.init();
            this.loadIpPrepareList();
        }
    });

})();