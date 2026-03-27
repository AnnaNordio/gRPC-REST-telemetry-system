import React from 'react';
import { Line } from 'react-chartjs-2';
// 1. Importa i componenti necessari da chart.js
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

// 2. Registra i componenti
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

export const LineChart = ({ history, measure, unit }) => {
  const recentHistory = history?.slice(-40) || [];
  const chartLabels = recentHistory.map(d => d.timestamp);
  const COLORS = {
    rest: '#6d28d9',
    grpc: '#ea580c'
    };

  const chartData = {
    labels: chartLabels,
    datasets: [
      {
        label: `REST (${unit})`,
        borderColor: COLORS.rest,
        backgroundColor: COLORS.rest,
        data: recentHistory.map(d => d.protocol.includes('REST') ? d.latency_ms : null),
        borderWidth: 2, pointRadius: 2, tension: 0, spanGaps: true,
      },
      {
        label: `gRPC (${unit})`,
        borderColor: COLORS.grpc,
        backgroundColor: COLORS.grpc,
        data: recentHistory.map(d => d.protocol.includes('gRPC') ? d.latency_ms : null),
        borderWidth: 2, pointRadius: 2, tension: 0, spanGaps: true,
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    plugins: { legend: { position: 'top' } },
    scales: {
      x: { 
        type: 'category', // Specifica esplicitamente il tipo se l'errore persiste
        grid: { display: false }, 
        ticks: { maxRotation: 0, autoSkip: true, maxTicksLimit: 10 } 
      },
      y: { 
        beginAtZero: true, 
        title: { display: true, text: `${measure} (${unit})` } 
      }
    }
  };

  return (
    <div className="bg-white p-6 rounded-2xl shadow-md border border-gray-50 h-[450px]">
      <Line data={chartData} options={chartOptions} />
    </div>
  );
};