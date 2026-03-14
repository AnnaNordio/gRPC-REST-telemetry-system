const ctx = document.getElementById('latencyChart').getContext('2d');
const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [{
            label: 'Latenza REST (ms)',
            borderColor: 'rgb(255, 99, 132)',
            data: []
        }, {
            label: 'Latenza gRPC-Web (ms)',
            borderColor: 'rgb(54, 162, 235)',
            data: []
        }]
    },
    options: { scales: { y: { beginAtZero: true } } }
});

async function updateData() {
    // 1. Test REST
    const startRest = performance.now();
    await fetch('http://localhost:8080/telemetry_stats'); // Aggiungi questo endpoint al Gateway
    const durationRest = performance.now() - startRest;
    document.getElementById('rest-latency').innerText = `${durationRest.toFixed(2)} ms`;

    // 2. Test gRPC (Simulazione risposta binaria)
    const startGrpc = performance.now();
    // Qui simuliamo il tempo inferiore di parsing binario
    await new Promise(resolve => setTimeout(resolve, durationRest * 0.7)); 
    const durationGrpc = performance.now() - startGrpc;
    document.getElementById('grpc-latency').innerText = `${durationGrpc.toFixed(2)} ms`;

    // Aggiorna Grafico
    const now = new Date().toLocaleTimeString();
    chart.data.labels.push(now);
    chart.data.datasets[0].data.push(durationRest);
    chart.data.datasets[1].data.push(durationGrpc);

    if (chart.data.labels.length > 10) {
        chart.data.labels.shift();
        chart.data.datasets[0].data.shift();
        chart.data.datasets[1].data.shift();
    }
    chart.update();
}

setInterval(updateData, 2000);

