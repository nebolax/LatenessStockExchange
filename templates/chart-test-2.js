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