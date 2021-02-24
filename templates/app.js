new Vue({
    el: "#app",

    data: {
        ws: null,
        curOffers: 0,
        chart: Chart
    },

    created: function() {
        var ctx = document.getElementById('chart').getContext('2d');
        console.log("ctx: " + ctx)
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

        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.onmessage = function(e) {
            var message = JSON.parse(e.data);
            switch (message.type) {
                case "gpoint":
                    var conf = chart.config.data;
                    conf.labels.push(0);
                    conf.datasets[0].data.push(message.stockPrice);
                    chart.update();
                    console.log(chart);
                    break;
                case "offers":
                    self.curOffers = message.offersCount;
                    break;
                default:
                    console.error("Unknow message type of mesage:");
                    console.error(message);
            }
        }
    }
})