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
            systemHosts: [{ip: 'test', domain: 'a', enabled: true},{ip: 'test', domain: 'a', enabled: true}],
            systemHostsLoading: false,
            checked: true,
            checked2: false,
            input5: '',
            ops: {
                vuescroll: {
                    mode: 'slide',
                    sizeStrategy: 'number',
                    /** Whether to detect dom resize or not */
                    detectResize: false
                },
                bar: {
                    background: '#ffffff',
                    opacity: 0.5,
                }
            }
        },
        methods: {
            changeTest: function(value) {
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
                this.$notify({
                    title: '提示',
                    message: '左侧分组可通过点击复选框快速启用及关闭哦~',
                    position: 'top-left'
                });
                console.log("no problem");
                let _this = this;
                server.sendMessage("list", {}, function(message) {
                    _this.systemHosts =  message.payload;
                    let t = _this;
                    setTimeout(() => {
                        t.systemHostsLoading = false;
                    }, 500);

                    //
                    // let html = "";
                    // for (let i in data) {
                    //     let item = data[i];
                    //     html += "<tr><th scope=\"row\">"+(parseInt(i)+1)+"</th><td>"+item.ip+"</td><td>"+item.domain+"</td><td>"+item.enabled+"</td></tr>";
                    // }
                    // document
                })
            }
        }
    });
    app.initSystemHosts();
})();