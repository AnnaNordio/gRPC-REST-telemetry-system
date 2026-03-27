import React, { useState } from 'react';
import { useTelemetry } from './hooks/useTelemetry';
import { ControlPanel } from './components/ControlPanel';
import Tabs from './components/Tabs';

// Import dei nuovi sottocomponenti
import { DashboardHeader } from './components/DashboardHeader';
import { LatencyView } from './components/views/LatencyView';
import { PayloadView } from './components/views/PayloadView';

import { getComparison } from './utils/benchmarkUtils'; 

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

  // Calcolo delle performance per i due tab
  const latencyComp = getComparison(restData.avg, grpcData.avg, 'latency');
  
  const restTotalSize = restData.size + (restData.overhead / 1024);
  const grpcTotalSize = grpcData.size + (grpcData.overhead / 1024);
  const sizeComp = getComparison(restTotalSize, grpcTotalSize, 'size');

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
        
        {/* Intestazione */}
        <DashboardHeader payloadSize={payloadSize} networkMode={isStreaming ? "Streaming" : "Polling"} />

        <div className="flex flex-col lg:flex-row gap-8 items-stretch">
          
          {/* Barra Laterale dei Controlli */}
          <aside className="lg:w-80 flex-shrink-0">
            <div className="sticky top-8">
              <ControlPanel 
                payloadSize={payloadSize} 
                onSizeChange={handleSizeChange} 
                isStreaming={isStreaming} 
                onModeToggle={handleModeToggle}
                isConnected={isConnected}
              />
              
              {/* Info Box Aggiuntiva opzionale */}
              <div className="mt-6 p-4 bg-blue-600 rounded-2xl text-white shadow-lg shadow-blue-200">
                <h4 className="text-xs font-bold uppercase tracking-widest opacity-80">System Status</h4>
                <p className="text-sm font-medium mt-1">
                  {isConnected ? "✅ Connected to Gateway" : "❌ Gateway Offline"}
                </p>
              </div>
            </div>
          </aside>

          {/* Area Principale dei Contenuti */}
          <main className="flex-1">
            <Tabs activeTab={activeTab} setActiveTab={setActiveTab} />

            <div className="mt-8 transition-all duration-300">
              {activeTab === 'latency' ? (
                <LatencyView 
                  restData={restData} 
                  grpcData={grpcData} 
                  history={history} 
                  comparison={latencyComp} 
                />
              ) : (
                <PayloadView 
                  restData={restData} 
                  grpcData={grpcData} 
                  sizeComp={sizeComp} 
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