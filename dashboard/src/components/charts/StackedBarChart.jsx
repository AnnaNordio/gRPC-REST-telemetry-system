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
        label: 'CPU Marshalling',
        data: [restData.marshalTime, grpcData.marshalTime],
        backgroundColor: ['#7c3aed', '#ea580c'], 
        borderRadius: 4,
        stack: 'Stack 0',
      },
      {
        label: 'Network & Logic',
        data: [
          restData.avg - restData.marshalTime,
          grpcData.avg - grpcData.marshalTime
        ],
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
            { text: 'CPU Marshalling (Solid)', fillStyle: '#4b5563', fontColor: '#666666' },
            { text: 'Network & Logic (Light)', fillStyle: '#d1d5db', fontColor: '#666666'},
          ]}
      },
      tooltip: {
        mode: 'index',
        intersect: false,
        callbacks: {
          label: (ctx) => ` ${ctx.dataset.label}: ${ctx.raw.toFixed(2)} μs`
        }
      }
    },
    scales: {
      x: { stacked: true, grid: { display: false } },
      y: { 
        stacked: true, 
        beginAtZero: true,
        title: { display: true, text: 'Time (μs)' },
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