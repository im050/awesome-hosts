let Server = function(){};

Server.prototype.sendMessage = function(name, payload, callback) {
    document.addEventListener('astilectron-ready', function() {
        // This will send a message to Go
        astilectron.sendMessage({name: name, payload: payload}, function(message) {
            console.log(message);
            callback(message)
        });
    });
};

(function () {
    let server = new Server();
    let app = new Vue({
        el: '#app',
        data: {
            systemHosts: [],//[{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true}],
            systemHostsLoading: false,
            fullscreenLoading: false,
            hostGroups: []
        },
        methods: {
            groupSwitch: function(value) {
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
            initSystemHosts: function() {
                let _this = this;
                setTimeout(function() {
                    console.log("yes")
                    _this.systemHosts = [{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true}]

                }, 2000);return
                this.$notify({
                    title: '提示',
                    message: '左侧分组可通过点击复选框快速启用及关闭哦~',
                    position: 'top-left'
                });

                server.sendMessage("event.name", "world", function(message) {
                    console.log(message)
                });
                server.sendMessage("list", {type: "SystemHosts"}, function(message) {
                    _this.systemHosts =  message.payload;
                    let t = _this;
                    setTimeout(() => {
                        t.systemHostsLoading = false;
                    }, 500);
                });
                this.getHostGroups();
            },
            getHostGroups: function() {
                let _this = this;
                server.sendMessage('groups', {}, function(message) {
                    _this.hostGroups = message.payload
                    let t = _this;
                    setTimeout(() => {
                            t.fullscreenLoading = false;
                        }, 1000)
                })
            }
        }
    });
    app.initSystemHosts();
})();