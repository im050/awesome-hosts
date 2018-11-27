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
            systemHosts: []
        },
        methods: {
            initSystemHosts: function() {
                console.log("no problem");
                let _this = this;
                server.sendMessage("list", {}, function(message) {
                    console.log("收到消息了");
                    _this.systemHosts =  message.payload;
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