import React, { useState } from 'react';
import { useTelemetry } from './hooks/useTelemetry';
import { ControlPanel } from './components/ControlPanel';
import Tabs from './components/Tabs';

// Import dei nuovi sottocomponenti
import { LatencyView } from './components/views/LatencyView';
import { PayloadView } from './components/views/PayloadView';

const Dashboard = () => {
  const [activeTab, setActiveTab] = useState('latency');

  const { 
    restData, 
    grpcData, 
    history, 
    isStreaming, 
    setIsStreaming, 
    payloadSize, 
    setPayloadSize, 
    setHistory, 
    isConnected 
  } = useTelemetry();
  

  // --- HANDLERS ---
  const handleModeToggle = async () => {
    const newMode = !isStreaming;
    const modeStr = newMode ? "streaming" : "polling";
    try {
      const resp = await fetch(`http://localhost:8080/set-mode?mode=${modeStr}`, { method: 'POST' });
      if (resp.ok) { 
        setIsStreaming(newMode); 
        setHistory([]); // Reset grafico al cambio modalità
      }
    } catch (err) {
      console.error("Errore cambio modalità:", err);
    }
  };

  const handleSizeChange = async (size) => {
    try {
      const resp = await fetch(`http://localhost:8080/set-size?size=${size}`, { method: 'POST' });
      if (resp.ok) { 
        setPayloadSize(size); 
        setHistory([]); // Reset grafico al cambio dimensione
      }
    } catch (err) {
      console.error("Errore cambio dimensione:", err);
    }
  };

  return (
    <div className="min-h-screen bg-slate-50 font-sans text-slate-900 selection:bg-blue-100">
      <div className="max-w-[1400px] mx-auto p-4 md:p-8">
        
        <header className="mb-10">
          <h1 className="text-3xl font-extrabold text-slate-800 tracking-tight">
            IoT Telemetry <span className="text-blue-600">Benchmark</span>
          </h1>
        </header>
        <div className="flex flex-col lg:flex-row gap-8 items-stretch">
          
          <aside className="lg:w-80 flex-shrink-0">
            <div className="sticky top-8">
              <ControlPanel 
                payloadSize={payloadSize} 
                onSizeChange={handleSizeChange} 
                isStreaming={isStreaming} 
                onModeToggle={handleModeToggle}
                isConnected={isConnected}
              />
              
              <div className="mt-6 p-4 bg-blue-600 rounded-2xl text-white shadow-lg shadow-blue-200">
                <h4 className="text-xs font-bold uppercase tracking-widest opacity-80">System Status</h4>
                <div className="flex items-center gap-2 mt-2">
                  <span 
                    className={`h-3 w-3 rounded-full flex-shrink-0 ${
                      isConnected 
                        ? "bg-green-400 shadow-[0_0_10px_rgba(74,222,128,0.9)]" 
                        : "bg-red-400 shadow-[0_0_10px_rgba(248,113,113,0.9)]"
                    }`}
                  ></span>
                  
                  <p className="text-sm font-medium leading-none">
                    {isConnected ? "Connected to Gateway" : "Gateway Offline"}
                  </p>
                </div>
              </div>
            </div>
          </aside>

          <main className="flex-1">
            <Tabs activeTab={activeTab} setActiveTab={setActiveTab} />

            <div className="mt-8 transition-all duration-300">
              {activeTab === 'latency' ? (
                <LatencyView 
                  restData={restData} 
                  grpcData={grpcData} 
                  history={history} 
                />
              ) : (
                <PayloadView 
                  restData={restData} 
                  grpcData={grpcData} 
                />
              )}
            </div>
          </main>

        </div>
      </div>
    </div>
  );
};

export default Dashboard;