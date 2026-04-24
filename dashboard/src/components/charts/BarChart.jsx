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
    labels: ['REST (JSON)', 'gRPC (Protobuf)'],
    datasets: [
      {
        label: `${measure} (${unit})`,
        backgroundColor: ['#6d28d9', '#ea580c'], 
        borderColor: ['#5b21b6', '#c2410c'],
        borderWidth: 1,
        data: [restValue || 0, grpcValue || 0],
        barThickness: 60, 
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (context) => ` ${context.raw} ${unit}`
        }
      }
    },
    scales: {
      y: { 
        stacked: true, 
        beginAtZero: true,
        title: { display: true, text: `${measure} (${unit})` },
        grid: { color: '#f1f5f9' }
      },
      x: {
        grid: { display: false }
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={chartData} options={chartOptions} />
    </div>
  );
};