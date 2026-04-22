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

export const LineChart = ({ history }) => {

  const getSampledData = (data) => {
    if (!data || data.length === 0) return [];

    const SAMPLE_INTERVAL_MS = 500; // Raggruppiamo i 100 sensori ogni 0.5s
    const slots = {};

    data.forEach(d => {
      const [hours, minutes, secondsWithMs] = d.timestamp.split(':');
      const totalMs = (parseInt(hours) * 3600000) + (parseInt(minutes) * 60000) + (parseFloat(secondsWithMs) * 1000);
      const currentSlot = Math.floor(totalMs / SAMPLE_INTERVAL_MS);

      const protoKey = d.protocol.includes('REST') ? 'rest' : 'grpc';
      if (!slots[currentSlot]) slots[currentSlot] = { timestamp: d.timestamp, rest: [], grpc: [] };
      
      // Accumuliamo le latenze nello slot
      slots[currentSlot][protoKey].push(d.latency_ms || d.latencyMs || d.latencyms);
    });

    return Object.values(slots).map(slot => {
      const calc = (lats) => {
        if (!lats.length) return null;
        return {
          avg: lats.reduce((a, b) => a + b, 0) / lats.length,
          min: Math.min(...lats),
          max: Math.max(...lats)
        };
      };
      return { timestamp: slot.timestamp, rest: calc(slot.rest), grpc: calc(slot.grpc) };
    }).slice(-40); 
  };

  const recentHistory = getSampledData(history);
  const chartLabels = recentHistory.map(d => d.timestamp);

  const COLORS = {
    rest: { main: '#6d28d9', area: 'rgba(109, 40, 217, 0.15)' },
    grpc: { main: '#ea580c', area: 'rgba(234, 88, 12, 0.15)' }
  };

  const chartData = {
    labels: chartLabels,
    datasets: [
      // --- DATASET REST ---
      {
        label: 'REST Max',
        data: recentHistory.map(d => d.rest?.max),
        borderColor: 'transparent',
        pointRadius: 0,
        fill: false,
      },
      {
        label: 'REST Range',
        data: recentHistory.map(d => d.rest?.min),
        borderColor: 'transparent',
        backgroundColor: COLORS.rest.area,
        fill: '-1', // Riempie l'area tra Max e Min
        pointRadius: 0,
      },
      {
        label: 'REST Avg',
        data: recentHistory.map(d => d.rest?.avg),
        borderColor: COLORS.rest.main,
        borderWidth: 3,
        pointRadius: 0,
        tension: 0.3,
      },
      // --- DATASET gRPC ---
      {
        label: 'gRPC Max',
        data: recentHistory.map(d => d.grpc?.max),
        borderColor: 'transparent',
        pointRadius: 0,
        fill: false,
      },
      {
        label: 'gRPC Range',
        data: recentHistory.map(d => d.grpc?.min),
        borderColor: 'transparent',
        backgroundColor: COLORS.grpc.area,
        fill: '-1',
        pointRadius: 0,
      },
      {
        label: 'gRPC Avg',
        data: recentHistory.map(d => d.grpc?.avg),
        borderColor: COLORS.grpc.main,
        borderWidth: 3,
        pointRadius: 0,
        tension: 0.3,
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
        labels: { boxWidth: 15, usePointStyle: true, font: { size: 12, weight: 'bold' }, 
          filter: (item) => item.text.includes('Avg') // Mostra solo la media in legenda
        }
      },
      tooltip: {
        filter: (tooltipItem) => {
          return tooltipItem.dataset.label.includes('Avg');
        },
        callbacks: {
          label: (ctx) => {
            const entry = recentHistory[ctx.dataIndex];
            const p = ctx.dataset.label.includes('REST') ? entry.rest : entry.grpc;
            return p ? [`Avg: ${p.avg.toFixed(2)}μs`, `Range: [${p.min.toFixed(1)}-${p.max.toFixed(1)}]μs`] : '';
          }
        }
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Line data={chartData} options={chartOptions} />
    </div>
  );
};