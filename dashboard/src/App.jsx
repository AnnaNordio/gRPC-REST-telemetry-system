import { useState, useEffect, useRef } from 'react';
import { ControlPanel } from './components/ControlPanel';
import { StatCard } from './components/StatCard';
import { LatencyChart } from './components/LatencyChart';

const Dashboard = () => {
  const [data, setData] = useState({ history: [], avg_rest: 0, avg_grpc: 0, p99_rest: 0, p99_grpc: 0 });
  const [isStreaming, setIsStreaming] = useState(false);
  const [payloadSize, setPayloadSize] = useState('small');
  const timerRef = useRef(null);

  const fetchData = async () => {
    try {
      const response = await fetch('/results');
      const rootData = await response.json();
      setData(rootData);
    } catch (err) { console.error("Fetch error:", err); }
  };

  useEffect(() => {
    fetchData();
    timerRef.current = setInterval(fetchData, 1000);
    return () => clearInterval(timerRef.current);
  }, []);

  const handleModeToggle = async () => {
    const newMode = !isStreaming;
    await fetch(`/set-mode?mode=${newMode ? "streaming" : "polling"}`, { method: 'POST' });
    setIsStreaming(newMode);
    setData(prev => ({ ...prev, history: [], avg_rest: 0, avg_grpc: 0 }));
  };

  const handleSizeChange = async (size) => {
    await fetch(`/set-size?size=${size}`, { method: 'POST' });
    setPayloadSize(size);
    setData(prev => ({ ...prev, history: [], avg_rest: 0, avg_grpc: 0 }));
  };

  return (
    <div className="min-h-screen w-full bg-gray-50 text-gray-900 font-sans p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        <header className="text-center mb-10">
          <h1 className="text-4xl font-extrabold tracking-tight text-gray-800">📡 IoT Telemetry Benchmark</h1>
          <p className="text-slate-500 mt-2">Real-time performance analysis: REST vs gRPC</p>
        </header>

        <ControlPanel 
          payloadSize={payloadSize} onSizeChange={handleSizeChange}
          isStreaming={isStreaming} onModeToggle={handleModeToggle}
        />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <StatCard 
            title="REST Average Latency" 
            subtitle="P99" 
            value={data.avg_rest} 
            subValue={data.p99_rest} 
            unit="µs"
            borderClass="border-violet-800" 
            textColor="text-violet-800" 
          />
          <StatCard 
            title="gRPC Average Latency" 
            subtitle="P99" 
            value={data.avg_grpc} 
            subValue={data.p99_grpc} 
            unit="µs"
            borderClass="border-orange-600" 
            textColor="text-orange-600" 
          />
        </div>

        <LatencyChart history={data.history} />
      </div>
    </div>
  );
};

export default Dashboard;