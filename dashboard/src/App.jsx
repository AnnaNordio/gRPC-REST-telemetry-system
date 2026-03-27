import { useTelemetry } from './hooks/useTelemetry';
import { ControlPanel } from './components/ControlPanel';
import { StatCard } from './components/StatCard';
import { LatencyChart } from './components/LatencyChart';

const Dashboard = () => {
  const { 
    restData, grpcData, history, isStreaming, 
    setIsStreaming, payloadSize, setPayloadSize, setHistory 
  } = useTelemetry();

  const handleModeToggle = async () => {
    const newMode = !isStreaming;
    const resp = await fetch(`http://localhost:8080/set-mode?mode=${newMode ? "streaming" : "polling"}`, { method: 'POST' });
    if (resp.ok) { setIsStreaming(newMode); setHistory([]); }
  };

  const handleSizeChange = async (size) => {
    const resp = await fetch(`http://localhost:8080/set-size?size=${size}`, { method: 'POST' });
    if (resp.ok) { setPayloadSize(size); setHistory([]); }
  };

  return (
    <div className="min-h-screen bg-slate-50 p-4 md:p-8 font-sans">
      <div className="max-w-6xl mx-auto">
        <header className="mb-10 text-center">
          <h1 className="text-4xl font-extrabold text-slate-800 tracking-tight">📡 IoT Telemetry Benchmark</h1>
        </header>

        <ControlPanel 
          payloadSize={payloadSize} onSizeChange={handleSizeChange} 
          isStreaming={isStreaming} onModeToggle={handleModeToggle} 
        />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-10">
          <StatCard title="REST Latency (Avg)" value={restData.avg} subtitle={"99th Percentile"} subValue={restData.p99} unit="µs" borderClass="border-violet-600" textColor="text-violet-700" />
          <StatCard title="gRPC Latency (Avg)" value={grpcData.avg} subtitle={"99th Percentile"} subValue={grpcData.p99} unit="µs" borderClass="border-orange-500" textColor="text-orange-600" />
        </div>

        <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
          <h3 className="text-lg font-semibold mb-4 text-slate-700">Real-time Comparison (Envoy proxy)</h3>
          <LatencyChart history={history} />
        </div>
      </div>
    </div>
  );
};

export default Dashboard;