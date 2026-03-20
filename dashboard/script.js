const ctx = document.getElementById('latencyChart').getContext('2d');

const chart = new Chart(ctx, {
    type: 'line',
    data: {
        labels: [],
        datasets: [
            {
                label: 'REST (µs)',
                borderColor: '#18069b',
                data: [],
                borderWidth: 2,
                pointRadius: 2,
                fill: false,
            },
            {
                label: 'gRPC (µs)',
                borderColor: '#ff5e00',
                data: [],
                borderWidth: 2,
                pointRadius: 2,
                fill: false,
            }
        ]
    },
    options: {
        responsive: true,
        animation: false,
        scales: {
            x: { title: { display: true, text: 'Tempo' } },
            y: { 
                beginAtZero: true, 
                title: { display: true, text: 'Latenza (µs)' } 
            }
        }
    }
});

async function fetchData() {
    try {
        const response = await fetch('/results');
        const rootData = await response.json(); 
        
        const history = rootData.history;
        if (!history || history.length === 0) return;

        const limit = 30;
        const restPoints = history.filter(d => d.protocol === 'REST').slice(-limit);
        const grpcPoints = history.filter(d => d.protocol === 'gRPC').slice(-limit);

        chart.data.labels = grpcPoints.map(d => d.timestamp);
        chart.data.datasets[0].data = restPoints.map(d => d.latency_ms);
        chart.data.datasets[1].data = grpcPoints.map(d => d.latency_ms);
        chart.update();

        document.getElementById('rest-lat').innerText = rootData.avg_rest.toFixed(2) + " µs";
        document.getElementById('grpc-lat').innerText = rootData.avg_grpc.toFixed(2) + " µs";

    } catch (err) {
        console.error("Errore nel recupero dati:", err);
    }
}


setInterval(fetchData, 1000);

async function toggleMode() {
    const isStreaming = document.getElementById('modeToggle').checked;
    const mode = isStreaming ? "streaming" : "polling";
    
    await fetch(`/set-mode?mode=${mode}`, { method: 'POST' });

    chart.data.labels = [];
    chart.data.datasets[0].data = []; 
    chart.data.datasets[1].data = []; 
    chart.update();

    document.getElementById('rest-lat').innerText = "-- µs";
    document.getElementById('grpc-lat').innerText = "-- µs";
}