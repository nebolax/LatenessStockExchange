var barCount = 60;
var initialDateStr = '01 Apr 2017 00:00 Z';

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
            data: []
        }]
    },
    options: {
        legend: {
            display: true,
            labels: {
                fontColor: 'rgb(255, 99, 132)'
            }
        }
    }
});

start();

function start() {
    console.log("started");
    var data = chart.config.data;
    var y = 100;
    for (var i = 0; i < 100; i++) {
        y += randVal(-5, 5)
        data.labels.push(i);
        data.datasets[0].data.push(y);
    }
    chart.update();
}

function addValue() {
    var data = chart.config.data;
    data.datasets[0].data.shift();
    data.datasets[0].data.push(randVal(-5, 30));
    chart.update();
}

function randVal(min, max) {
    return Math.random() * (max - min) + min;
}