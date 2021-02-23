new Vue({
    el: "#app",

    data: {
        ws: null,
        curVal: 10
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.onmessage = function(e) {
            self.curVal = 'lol ' + JSON.parse(e.data).value
        }
    }
})