import { useState } from 'react';
import { useTelemetry } from './hooks/useTelemetry';
import { ControlPanel } from './components/ControlPanel';
import { StatCard } from './components/StatCard';
import { Chart } from './components/Chart';
import Tabs from './components/Tabs'; // Assicurati che il file esista in components/

const Dashboard = () => {
  const [activeTab, setActiveTab] = useState('latency'); // Stato per gestire il menu a schede

  const { 
    restData, grpcData, history, isStreaming, 
    setIsStreaming, payloadSize, setPayloadSize, setHistory 
  } = useTelemetry();

  const handleModeToggle = async () => {
    const newMode = !isStreaming;
    const modeStr = newMode ? "streaming" : "polling";
    const resp = await fetch(`http://localhost:8080/set-mode?mode=${modeStr}`, { method: 'POST' });
    if (resp.ok) { 
      setIsStreaming(newMode); 
      setHistory([]); 
    }
  };

  const handleSizeChange = async (size) => {
    const resp = await fetch(`http://localhost:8080/set-size?size=${size}`, { method: 'POST' });
    if (resp.ok) { 
      setPayloadSize(size); 
      setHistory([]); 
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 font-sans">
      <div className="max-w-[1400px] mx-auto p-4 md:p-8">
        
        <header className="mb-10">
          <h1 className="text-3xl font-extrabold text-slate-800 tracking-tight">
            📡 IoT Telemetry <span className="text-blue-600">Benchmark</span>
          </h1>
        </header>

        {/* --- LAYOUT PRINCIPALE --- */}
        <div className="flex flex-col lg:flex-row gap-8">
          
          {/* SIDEBAR SINISTRA (Control Panel) */}
          <aside className="lg:w-80 flex-shrink-0">
            <ControlPanel 
              payloadSize={payloadSize} 
              onSizeChange={handleSizeChange} 
              isStreaming={isStreaming} 
              onModeToggle={handleModeToggle} 
            />
          </aside>

          {/* CONTENUTO PRINCIPALE (Tabs + Cards + Charts) */}
          <main className="flex-1">
            <Tabs activeTab={activeTab} setActiveTab={setActiveTab} />

            <div className="mt-8">
              {activeTab === 'latency' && (
                <div className="space-y-8 animate-in fade-in slide-in-from-left-4 duration-500">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <StatCard title="REST Avg" value={restData.avg} subtitle="99th Percentile" subValue={restData.p99} unit="μs" borderClass="border-violet-600" textColor="text-violet-700" />
                    <StatCard title="gRPC Avg" value={grpcData.avg} subtitle="99th Percentile" subValue={grpcData.p99} unit="μs" borderClass="border-orange-500" textColor="text-orange-600" />
                  </div>
                  <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <Chart history={history} measure="Microseconds" unit="μs"/>
                  </div>
                </div>
              )}

              {activeTab === 'payload' && (
                <div className="space-y-8 animate-in fade-in slide-in-from-right-4 duration-500">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <StatCard title="REST Payload Size" value={restData.size || 0} subtitle="Json overhead" subValue={restData.overhead} unit="KB" borderClass="border-blue-600" textColor="text-blue-700" />
                    <StatCard title="gRPC Payload Size" value={grpcData.size || 0} subtitle="Protobuf overhead" subValue={grpcData.overhead} unit="KB" borderClass="border-orange-500" textColor="text-orange-600" />
                  </div>
                  <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-200">
                    <Chart history={history} measure="Size" unit="KB"/>
                  </div>
                </div>
              )}
            </div>
          </main>
        </div>
      </div>
    </div>
  );
};

export default Dashboard;