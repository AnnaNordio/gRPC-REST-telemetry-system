import React from 'react';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement, // Questo è l'elemento "bar" che mancava!
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

// Registriamo i moduli necessari per far funzionare il grafico
ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

export const GroupedBarChart = ({ restSize, restOverhead, grpcSize, grpcOverhead, unit = 'KB' }) => {
  const COLORS = {
    rest: '#6d28d9',
    grpc: '#ea580c'
  };

  const chartData = {
    labels: ['REST (JSON)', 'gRPC (Proto)'],
    datasets: [
      {
        label: `Payload Size (${unit})`,
        backgroundColor: [COLORS.rest + '88', COLORS.grpc + '88'],
        borderColor: [COLORS.rest, COLORS.grpc],
        borderWidth: 1,
        data: [restSize || 0, grpcSize || 0],
      },
      {
        label: `Overhead (${unit})`,
        backgroundColor: [COLORS.rest, COLORS.grpc],
        borderColor: [COLORS.rest, COLORS.grpc],
        borderWidth: 1,
        data: [restOverhead || 0, grpcOverhead || 0],
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { position: 'top' },
      tooltip: {
        callbacks: {
          title: (items) => items[0].label,
          label: (context) => `${context.dataset.label}: ${context.raw} ${unit}`
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: { display: true, text: unit },
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={chartData} options={chartOptions} />
    </div>
  );
};