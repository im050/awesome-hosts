let Server = function () {
};

Server.prototype.sendMessage = function (name, payload, callback) {
    // This will send a message to Go
    astilectron.sendMessage({name: name, payload: payload}, function (message) {
        console.log(message);
        callback(message.payload)
    });
};

(function () {
    let system = {
        isSystemHosts: true,
        currentGroupName: '',
        currentHosts: [],
        systemHosts: []
    };
    const SYSTEM_HOSTS_NAME = "System Hosts";
    let server = new Server();
    let app = new Vue({
        el: '#app',
        data: {
            currentPage: '',
            hostsLoading: false,
            fullscreenLoading:  false,
            addHostLoading: false,
            hostGroups: [],
            inputIp: '',
            inputHost: '',
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
                server.sendMessage("enableGroup", {groupName: this.system.currentGroupName, enabled: value}, (message) => {
                    if (message.code === 1) {
                        if (value) {
                            this.$message({
                                message: 'The group has been enabled.',
                                type: 'success'
                            });
                        } else {
                            this.$message({
                                message: 'The group has been disabled.',
                                type: 'info'
                            });
                        }
                    } else {
                        this.$message({
                            message: 'An error occured while your operating',
                            type: 'error'
                        });
                    }
                });

            },
            addHost: function () {
                if (this.system.currentGroupName === SYSTEM_HOSTS_NAME) {
                    return ;
                }
                let ip = this.inputIp;
                let domain = this.inputHost;
                if (ip === '' || domain === '') {
                    this.$message({
                        message: "IP or Domain was empty.",
                        type: "error"
                    });
                    return ;
                }

                let groupName = this.system.currentGroupName;
                this.addHostLoading = true;
                server.sendMessage("addHost", {groupName: groupName, ip: ip, domain: domain}, (message) => {
                    this.system.currentHosts.push({
                        ip: ip,
                        domain: domain,
                        enabled: true
                    });
                    //focus on ip input
                    this.$refs.ipInput.focus();
                    this.inputIp = '';
                    this.inputHost = '';
                    this.addHostLoading = false;
                    this.$message({
                        message: 'Added successfully',
                        type: 'success'
                    });
                });
            },
            init: function () {
                let _this = this;
                Promise.all([this.getList(SYSTEM_HOSTS_NAME), this.getHostGroups()]).then((results) => {
                    _this.hostGroups[0].active = true;
                    _this.system.currentGroupName = _this.hostGroups[0].name;
                    setTimeout(() => {
                        _this.fullscreenLoading = false;
                        _this.hostsLoading = false;
                    }, 1000);
                    console.log(results);
                });
            },
            changeHost: function (value, index) {
                let groupName = this.system.currentGroupName;
                let host = this.system.currentHosts[index];
                if (host === null) {
                    this.$message({
                        message: 'Badly Host',
                        type: 'error'
                    });
                }
                server.sendMessage(
                    "updateHost",
                    {groupName: groupName, ip: host.ip, domain: host.domain, enabled: value, index: index},
                    (message) => {
                        if (message.code !== 1) {
                            //it turns switch button to old status
                            this.system.currentHosts[index].enabled = !value;
                            this.$message({
                                message: 'An error occured while updating host',
                                type: 'error'
                            });
                        }
                });
            },
            changeGroup: function (groupName) {
                this.system.isSystemHosts = (groupName === SYSTEM_HOSTS_NAME);
                for (let i in this.hostGroups) {
                    let item = this.hostGroups[i];
                    if (groupName === item.name) {
                        item.active = true;
                        this.system.currentHosts = (groupName === SYSTEM_HOSTS_NAME)
                            ? this.system.systemHosts :
                            (item.hosts === null ? [] : item.hosts);
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
                        if (message.payload === null) {
                            message.payload = [];
                        }
                        _this.system.currentHosts = message.payload;
                        if (type === SYSTEM_HOSTS_NAME) {
                            _this.system.systemHosts = message.payload;
                        }
                        _this.hostsLoading = false;
                        resolve(true)
                    });
                })
            },
            getHostGroups: function () {
                let _this = this;
                return new Promise((resolve) => {
                    server.sendMessage('groups', {}, (message) => {
                        if (message.payload === null) {
                            message.payload = []
                        }
                        _this.hostGroups = message.payload;
                        _this.hostGroups.unshift({
                            name: SYSTEM_HOSTS_NAME,
                            active: false,
                        });
                        resolve(true);
                    })
                })
            },
            needPassword: function (payload) {
                this.$prompt('In order to sync hosts file you have to type in the administrator password.', 'Password', {
                    confirmButtonText: 'Confirm',
                    cancelButtonText: 'Cancel'
                }).then(({ value }) => {
                    server.sendMessage(payload, {password: value}, (message) => {
                        if (message.code == 1) {
                            this.$message({
                                type: 'success',
                                message: 'Synchronization success'
                            });
                        } else {
                            this.$message({
                                type: 'error',
                                message: 'An error occured while your operating'
                            });
                        }
                    });
                }).catch(() => {
                    this.$message({
                        type: 'error',
                        message: 'Synchronization failure'
                    });
                });
            }
        },
        mounted() {
            document.addEventListener('astilectron-ready',  () => {
                this.init();
                this.loadIpPrepareList();
                astilectron.onMessage((message) => {
                    switch (message.name) {
                        case 'needPassword':
                            this.needPassword(message.payload);
                            break;
                    }
                });
            })

        }
    });
})();