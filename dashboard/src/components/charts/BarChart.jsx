import React from 'react';
import { Bar } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

export const BarChart = ({ restValue, grpcValue, measure, unit = 'KB' }) => {
  const chartData = {
    // Le etichette sull'asse X: i due protocolli a confronto
    labels: ['REST (JSON)', 'gRPC (Protobuf)'],
    datasets: [
      {
        label: `${measure} (${unit})`,
        // Colore differenziato per le due barre
        backgroundColor: ['#6d28d9', '#ea580c'], 
        borderColor: ['#5b21b6', '#c2410c'],
        borderWidth: 1,
        // I due valori da confrontare
        data: [restValue || 0, grpcValue || 0],
        // Opzionale: restringe un po' le barre per renderle più eleganti
        barThickness: 60, 
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false }, // Nascondiamo la legenda perché i nomi sono già sulle X
      tooltip: {
        callbacks: {
          label: (context) => ` ${context.raw} ${unit}`
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        title: { 
          display: true, 
          text: `${measure} (${unit})`,
          font: { weight: 'bold' }
        },
      },
      x: {
        grid: { display: false } // Pulizia visiva
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={chartData} options={chartOptions} />
    </div>
  );
};