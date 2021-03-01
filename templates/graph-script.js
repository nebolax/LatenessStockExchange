var vm = new Vue({
    el: "#main",
    data: {
        curOffers: 0
    }
})

var ctx = document.getElementById('chart').getContext('2d');
ctx.canvas.width = 1000;
ctx.canvas.height = 1000;
var chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            borderColor: 'rgb(107, 87, 250)',
            pointRadius: 0,
            borderWidth: 1.5,
            fill: false,
            data: []
        }]
    },
    options: {
        legend: {
            display: false
        }
    }
});

ws = new WebSocket('ws://' + window.location.host + "/ws" + window.location.pathname);
ws.onmessage = function(e) {
    var message = JSON.parse(e.data);
    switch (message.type) {
        case "setup":
            var conf = chart.config.data;
            conf.labels = [];
            conf.datasets[0].data = []
            for (let i = 0; i < message.history.length; i++) {
                conf.labels.push(0);
                conf.datasets[0].data.push(message.history[i])
            }
            chart.update();
            vm.curOffers = message.offers;
            break;
        case "gpoint":
            var conf = chart.config.data;
            conf.labels.push(0);
            conf.datasets[0].data.push(message.stockPrice);
            if (conf.labels.length > 20) {
                conf.labels = conf.labels.slice(conf.labels.length - 20)
                conf.datasets[0].data = conf.datasets[0].data.slice(conf.datasets[0].data.length - 20)
            }
            chart.update();
            break;
        case "offers":
            vm.curOffers = message.offersCount;
            break;
        default:
            console.error("Unknow message type of mesage:");
            console.error(message);
    }
}

$("#buyOfferBtn").click(function() {
    ws.send(JSON.stringify({
        offerType: "buy"
    }))
});

$("#sellOfferBtn").click(function() {
    ws.send(JSON.stringify({
        offerType: "sell"
    }))
});