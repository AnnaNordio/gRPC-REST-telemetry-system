import React from 'react';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement, // Registriamo l'elemento Bar
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

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
    // Sull'asse X ora definiamo cosa stiamo guardando
    labels: ['Payload Size', 'Overhead Size'],
    datasets: [
      {
        label: 'REST (JSON)',
        backgroundColor: COLORS.rest, // Colore fisso viola
        borderColor: COLORS.rest,
        borderWidth: 1,
        // Dati per REST: [payload, overhead]
        data: [restSize || 0, restOverhead || 0],
      },
      {
        label: 'gRPC (Proto)',
        backgroundColor: COLORS.grpc, // Colore fisso arancione
        borderColor: COLORS.grpc,
        borderWidth: 1,
        // Dati per gRPC: [payload, overhead]
        data: [grpcSize || 0, grpcOverhead || 0],
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
          label: (context) => `${context.dataset.label}: ${context.raw}`
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: { display: true, text: 'Kilobyte (KB)' },
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={chartData} options={chartOptions} />
    </div>
  );
};