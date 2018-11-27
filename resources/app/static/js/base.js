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
    let app = new Vue({
        el: '#app',
        data: {
            sites: [
                { name: 'Runoob' },
                { name: 'Google' },
                { name: 'Taobao' }
            ]
        }
    })
    let server = new Server();
    server.sendMessage("list", {}, function(message) {
        let data = message.payload;
        for (let i in data) {
            let item = data[i];
            let html = "<tr><th scope=\"row\">"+(parseInt(i)+1)+"</th><td>"+item.ip+"</td><td>"+item.domain+"</td><td>"+item.enabled+"</td></tr>";
            $("#host-list").append(html)
        }
    })
})();