import React from 'react';
import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
} from 'chart.js';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

export const LineChart = ({ history, unit }) => {

  const getSampledData = (data) => {
    if (!data || data.length === 0) return [];

    const SAMPLE_INTERVAL_MS = 100;
    // Usiamo un oggetto per raggruppare i dati per "slot"
    const slots = {};

    data.forEach(d => {
      const [hours, minutes, secondsWithMs] = d.timestamp.split(':');
      const totalMs = 
        (parseInt(hours) * 3600000) + 
        (parseInt(minutes) * 60000) + 
        (parseFloat(secondsWithMs) * 1000);

      // Calcoliamo lo slot univoco
      const currentSlot = Math.floor(totalMs / SAMPLE_INTERVAL_MS);

      // Se lo slot non esiste ancora, lo creiamo
      if (!slots[currentSlot]) {
        slots[currentSlot] = {
          timestamp: d.timestamp, // Usiamo il timestamp del primo che arriva
          rest: null,
          grpc: null
        };
      }

      // Inseriamo i dati nel protocollo corretto all'interno dello stesso slot
      if (d.protocol.includes('REST')) {
        slots[currentSlot].rest = d;
      } else if (d.protocol.includes('gRPC')) {
        slots[currentSlot].grpc = d;
      }
    });

    // Convertiamo l'oggetto in un array ordinato e prendiamo gli ultimi 50
    return Object.values(slots).slice(-50);
  };
  const recentHistory = getSampledData(history);
  const chartLabels = recentHistory.map(d => d.timestamp);
  console.log("LineChart History:", recentHistory); // Debug: Verifica i dati in ingresso
  const COLORS = {
    rest: '#6d28d9',
    grpc: '#ea580c',
    grid: '#f1f5f9'
  };

  const chartData = {
    labels: chartLabels,
    datasets: [
      // --- REST DATA ---
      {
        label: `REST Avg`,
        borderColor: COLORS.rest,
        backgroundColor: COLORS.rest,
        data: recentHistory.map(d => d.rest ? d.rest.latency_ms : null),        
        borderWidth: 2.5,
        pointRadius: 0,
        tension: 0.3,
        spanGaps: true,
      },
      {
        label: `REST P99 (Tail)`,
        borderColor: COLORS.rest,
        borderDash: [5, 5], // Linea tratteggiata
        data: recentHistory.map(d => d.rest ? d.rest.p99_ms : null),        
        borderWidth: 1,
        pointRadius: 0,
        fill: false,
        spanGaps: true,
      },
      // --- gRPC DATA ---
      {
        label: `gRPC Avg`,
        borderColor: COLORS.grpc,
        backgroundColor: COLORS.grpc,
        data: recentHistory.map(d => d.grpc ? d.grpc.latencyms : null),        
        borderWidth: 2.5,
        pointRadius: 0,
        tension: 0.3,
        spanGaps: true,
      },
      {
        label: `gRPC P99 (Tail)`,
        borderColor: COLORS.grpc,
        borderDash: [5, 5],
        data: recentHistory.map(d => d.grpc ? d.grpc.p99 : null),        
        borderWidth: 1,
        pointRadius: 0,
        fill: false,
        spanGaps: true,
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    animation: false, // Disabilitato per performance con molti sensori
    interaction: {
      mode: 'index',
      intersect: false,
    },
    plugins: { 
      legend: { 
        position: 'top',
        labels: { boxWidth: 15, usePointStyle: true, font: { size: 12, weight: 'bold' } }
      },
      tooltip: {
        enabled: true,
        padding: 12,
        backgroundColor: 'rgba(255, 255, 255, 0.95)',
        titleColor: '#1e293b',
        bodyColor: '#1e293b',
        borderColor: '#e2e8f0',
        borderWidth: 1,
      }
    },
    scales: {
      x: { 
        grid: { display: false }, 
      },
      y: { 
        beginAtZero: false,
        grid: { color: COLORS.grid }, 
        title: { display: true, text: `Latency (${unit})`, font: { weight: 'bold' } } 
      }
    }
  };

  return <Line data={chartData} options={chartOptions} />;
};