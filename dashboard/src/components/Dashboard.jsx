import React, { useState, useEffect, useRef } from 'react';
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
} from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

const Dashboard = () => {
  const [data, setData] = useState({ history: [], avg_rest: 0, avg_grpc: 0 });
  const [isStreaming, setIsStreaming] = useState(false);
  const timerRef = useRef(null);

  const fetchData = async () => {
    try {
      const response = await fetch('/results');
      const rootData = await response.json();
      setData(rootData);
    } catch (err) {
      console.error("Errore nel recupero dati:", err);
    }
  };

  useEffect(() => {
    fetchData();
    timerRef.current = setInterval(fetchData, 1000);
    return () => clearInterval(timerRef.current);
  }, []);

  const toggleMode = async () => {
    const newMode = !isStreaming;
    const modeString = newMode ? "streaming" : "polling";
    
    try {
      await fetch(`/set-mode?mode=${modeString}`, { method: 'POST' });
      setIsStreaming(newMode);
      setData({ history: [], avg_rest: 0, avg_grpc: 0 });
    } catch (err) {
      console.error("Errore nel cambio modalità:", err);
    }
  };

  const limit = 30;
  const restPoints = data.history?.filter(d => d.protocol === 'REST').slice(-limit) || [];
  const grpcPoints = data.history?.filter(d => d.protocol === 'gRPC').slice(-limit) || [];

  const chartData = {
    labels: grpcPoints.map(d => d.timestamp),
    datasets: [
      {
        label: 'REST (µs)',
        borderColor: '#1e40af', // Blue 800
        backgroundColor: '#1e40af',
        data: restPoints.map(d => d.latency_ms),
        borderWidth: 2,
        pointRadius: 1,
        tension: 0.3,
      },
      {
        label: 'gRPC (µs)',
        borderColor: '#ea580c', // Orange 600
        backgroundColor: '#ea580c',
        data: grpcPoints.map(d => d.latency_ms),
        borderWidth: 2,
        pointRadius: 1,
        tension: 0.3,
      }
    ]
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    animation: false,
    plugins: {
      legend: { position: 'top' }
    },
    scales: {
      x: { grid: { display: false } },
      y: { beginAtZero: true }
    }
  };

  return (
    <div className="min-h-screen w-full bg-gray-50 text-gray-900 font-sans">
      <div className="max-w-6xl mx-auto p-4 md:p-8">
        
        {/* Header */}
        <div className="text-center mb-10">
          <h1 className="text-4xl font-extrabold tracking-tight text-gray-800">
            📡 IoT Telemetry Dashboard
          </h1>
          <p className="text-gray-500 mt-2 italic">Monitoraggio latenza gRPC vs REST</p>
        </div>

        {/* Modalità Toggle Card */}
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-8 flex flex-col md:flex-row items-center justify-center gap-6">
          <div className="flex items-center gap-3">
            <span className={`text-sm font-bold uppercase tracking-wider ${isStreaming ? 'text-gray-400' : 'text-blue-700'}`}>
              Polling
            </span>
            
            {/* Switch stilizzato */}
            <label className="relative inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                checked={isStreaming} 
                onChange={toggleMode} 
                className="sr-only peer" 
              />
              <div className="w-14 h-7 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:bg-orange-600"></div>
            </label>

            <span className={`text-sm font-bold uppercase tracking-wider ${isStreaming ? 'text-orange-600' : 'text-gray-400'}`}>
              Streaming
            </span>
          </div>
        </div>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white p-6 rounded-2xl shadow-sm border-l-8 border-blue-700">
            <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest">Avg Latency REST</h3>
            <div className="text-4xl font-black text-blue-800 mt-2">
              {data.avg_rest ? `${data.avg_rest.toFixed(2)} µs` : '--'}
            </div>
          </div>
          
          <div className="bg-white p-6 rounded-2xl shadow-sm border-l-8 border-orange-600">
            <h3 className="text-sm font-semibold text-gray-400 uppercase tracking-widest">Avg Latency gRPC</h3>
            <div className="text-4xl font-black text-orange-600 mt-2">
              {data.avg_grpc ? `${data.avg_grpc.toFixed(2)} µs` : '--'}
            </div>
          </div>
        </div>

        {/* Chart Card */}
        <div className="bg-white p-6 rounded-2xl shadow-md border border-gray-50 h-[450px]">
          <Line data={chartData} options={chartOptions} />
        </div>

      </div>
    </div>
  );
};

export default Dashboard;