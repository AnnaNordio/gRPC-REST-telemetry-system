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
    // Sull'asse X confrontiamo i due protocolli
    labels: ['REST (JSON)', 'gRPC (Proto)'],
    datasets: [
      {
        label: `Payload Netto (${unit})`,
        // Colore pieno ma leggermente opaco per il dato "puro"
        backgroundColor: [COLORS.rest + '88', COLORS.grpc + '88'],
        borderColor: [COLORS.rest, COLORS.grpc],
        borderWidth: 1,
        // Prendiamo i dati reali passati dalle props
        data: [restSize || 0, grpcSize || 0],
      },
      {
        label: `Overhead (${unit})`, // L'overhead spesso è in Byte, adatta se necessario
        // Colore solido e acceso per evidenziare lo "spreco"
        backgroundColor: [COLORS.rest, COLORS.grpc],
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
          label: (context) => `${context.dataset.label}: ${context.raw}`
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: { display: true, text: 'Size / Overhead' },
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={chartData} options={chartOptions} />
    </div>
  );
};