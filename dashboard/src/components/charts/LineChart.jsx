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
  const recentHistory = history?.slice(-50) || [];
  console.log("LineChart - Recent History:", recentHistory); // Debug log per verificare i dati in ingresso
  const chartLabels = recentHistory.map(d => d.timestamp);

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
        data: recentHistory.map(d => d.protocol.includes('REST') ? d.latency_ms : null),
        borderWidth: 2.5,
        pointRadius: 0,
        tension: 0.3,
        spanGaps: true,
      },
      {
        label: `REST P99 (Tail)`,
        borderColor: COLORS.rest,
        borderDash: [5, 5], // Linea tratteggiata
        data: recentHistory.map(d => d.protocol.includes('REST') ? d.p99_ms : null),
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
        data: recentHistory.map(d => d.protocol.includes('gRPC') ? d.latency_ms : null),
        borderWidth: 2.5,
        pointRadius: 0,
        tension: 0.3,
        spanGaps: true,
      },
      {
        label: `gRPC P99 (Tail)`,
        borderColor: COLORS.grpc,
        borderDash: [5, 5],
        data: recentHistory.map(d => d.protocol.includes('gRPC') ? d.p99 : null),
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
        ticks: { display: false } // Nascondiamo i timestamp fitti per pulizia
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