import React, { useState, useCallback } from 'react';
import { useTelemetry } from './hooks/useTelemetry';
import { ControlPanel } from './components/ControlPanel';
import Tabs from './components/Tabs';

// Import dei nuovi sottocomponenti
import { LatencyView } from './components/views/LatencyView';
import { PayloadView } from './components/views/PayloadView';
import { ThroughputView } from './components/views/ThroughputView';
import { MarshalView } from './components/views/MarshalView';
import { ConnectionCard } from './components/ConnectionCard';

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
    sensorNumber,
    setSensorNumber,
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

  const handleSensorChange = useCallback(async (count) => {
    try {
      const resp = await fetch(`http://localhost:8080/set-sensors?count=${count}`, { 
        method: 'POST' 
      });
      if (resp.ok) {
        setSensorNumber(count);
        setHistory([]); 
      }
    } catch (err) {
      console.error("Errore aggiornamento sensori:", err);
    }
  }, [setSensorNumber, setHistory]);

  const renderActiveView = () => {
    switch (activeTab) {
      case 'latency':
        return <LatencyView restData={restData} grpcData={grpcData} history={history} />;
      case 'payload':
        return <PayloadView restData={restData} grpcData={grpcData} />;
      case 'scalability':
        return <ThroughputView restData={restData} grpcData={grpcData} />;
      case 'marshalling':
        return <MarshalView restData={restData} grpcData={grpcData} history={history} />;
      default:
        return null;
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
                onSensorChange={handleSensorChange}
                isConnected={isConnected}
              />
              
              <ConnectionCard isConnected={isConnected} />
            </div>
          </aside>

          <main className="flex-1">
            <Tabs activeTab={activeTab} setActiveTab={setActiveTab} />

            <div className="mt-8 transition-all duration-300">
              {renderActiveView()}
            </div>
          </main>

        </div>
      </div>
    </div>
  );
};

export default Dashboard;