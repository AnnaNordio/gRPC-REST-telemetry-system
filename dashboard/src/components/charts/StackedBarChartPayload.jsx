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

export const StackedBarChart = ({ restData, grpcData }) => {
  const data = {
    labels: ['REST', 'gRPC'],
    datasets: [
      {
        label: 'Payload Size',
        data: [restData.payloadSize, grpcData.payloadSize],
        // Viola per REST, Arancio per gRPC (Pieni)
        backgroundColor: ['#7c3aed', '#ea580c'], 
        borderRadius: 4,
        stack: 'Stack 0',
      },
      {
        label: 'Overhead Size',
        data: [
          restData.overheadSize, grpcData.overheadSize
        ],
        // Viola sbiadito per REST, Arancio sbiadito per gRPC
        backgroundColor: ['#ddd6fe', '#ffedd5'], 
        borderRadius: 4,
        stack: 'Stack 0',
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    
    plugins: {
      legend: { 
        position: 'top',
        labels: { boxWidth: 15, usePointStyle: true, font: { size: 12, weight: 'bold' }, generateLabels: (chart) => [
            { text: 'Payload Size (Solid)', fillStyle: '#4b5563', fontColor: '#666666' },
            { text: 'Overhead Size (Light)', fillStyle: '#d1d5db', fontColor: '#666666'},
          ]}
      },
      tooltip: {
        callbacks: {
          label: (context) => `${context.dataset.label}: ${context.raw}`
        }
      }
    },
    scales: {
      x: { stacked: true, grid: { display: false } },
      y: { 
        stacked: true, 
        beginAtZero: true,
        title: { display: true, text: 'Size (KB)' },
        grid: { color: '#f1f5f9' }
      }
    }
  };

  return (
    <div className="h-[400px] w-full">
      <Bar data={data} options={options} />
    </div>
  );
};