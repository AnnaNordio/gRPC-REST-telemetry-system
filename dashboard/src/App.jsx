import { useState, useEffect, useRef } from 'react';
import { ControlPanel } from './components/ControlPanel';
import { StatCard } from './components/StatCard';
import { LatencyChart } from './components/LatencyChart';

// 1. Import dal bundle generato da Webpack (dist/bundle.js)
// Questo risolve il problema del "require is not defined"
import * as protos from 'my-grpc-protos';

const Dashboard = () => {
  const [restData, setRestData] = useState({ avg: 0, p99: 0 });
  const [grpcData, setGrpcData] = useState({ avg: 0, p99: 0 });
  const [history, setHistory] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [payloadSize, setPayloadSize] = useState('small');

  // Inizializzazione del client gRPC puntando a Envoy (porta 8081)
  const grpcClient = useRef(new protos.TelemetryServiceClient('http://localhost:8081'));
  // --- LOGICA REST (Polling) ---
  const fetchRestData = async () => {
    try {
      const response = await fetch('/results');
      const data = await response.json();
      
      setRestData({ avg: data.avg_rest, p99: data.p99_rest });

      if (isStreaming) {
        setHistory(prev => {
          const newRestPoints = data.history.filter(h => h.protocol === 'REST');
          const combined = [...prev, ...newRestPoints];
          const unique = Array.from(new Map(combined.map(item => [item.timestamp + item.protocol, item])).values());
          return unique.slice(-50); 
        });
      } else {
        setHistory(data.history);
      }
    } catch (err) { console.error(err); }
  };

  const fetchGrpcUpdate = () => {
    if (isStreaming) return;

    grpcClient.current.getStats(new protos.Empty(), {}, (err, response) => {
      if (err) {
        console.error("gRPC Polling error:", err);
        return;
      }
      const g = response.toObject();
      setGrpcData({ avg: g.avgLatency || 0, p99: g.p99Latency || 0 });
    });
  };


  useEffect(() => {
    const interval = setInterval(() => {
      fetchRestData();    // Polling REST per box viola
      fetchGrpcUpdate();  // Polling gRPC per box arancione
    }, 1000);
    return () => clearInterval(interval);
  }, [isStreaming]);

  // --- LOGICA gRPC (Streaming) ---
  useEffect(() => {
    if (!isStreaming) return;

    try {
      // Usiamo protos.Empty direttamente dal bundle
      const stream = grpcClient.current.getGrpcStream(new protos.Empty(), {});
      stream.on('data', (response) => {
        const g = response.toObject();
        const avg = g.avgLatency || 0;
        
        setGrpcData({ avg, p99: g.p99Latency || 0 });
        setHistory(prev => [
          ...prev.slice(-49), 
          { 
            timestamp: new Date().toLocaleTimeString(), 
            protocol: 'gRPC', 
            latency_ms: avg 
          }
        ]);
      });

      stream.on('error', (err) => {
        console.error("gRPC Stream error:", err);
        setIsStreaming(false); // Opzionale: ferma lo streaming in caso di errore
      });
      
      return () => {
        if (stream) stream.cancel();
      };
    } catch (e) {
      console.error("Failed to initialize gRPC stream:", e);
    }
  }, [isStreaming]);

  // --- GESTORI EVENTI ---
 const handleModeToggle = async () => {
    const newMode = !isStreaming;
    const modeString = newMode ? "streaming" : "polling";
    
    try {
      await fetch(`/set-mode?mode=${modeString}`, { method: 'POST' });
      
      setIsStreaming(newMode);
      
      setRestData({ avg: 0, p99: 0 });
      setGrpcData({ avg: 0, p99: 0 });
      
    } catch (e) { 
      console.error("Error toggling mode:", e); 
    }
};

  const handleSizeChange = async (size) => {
    try {
      await fetch(`/set-size?size=${size}`, { method: 'POST' });
      setPayloadSize(size);
      setRestData({ avg: 0, p99: 0 });
      setGrpcData({ avg: 0, p99: 0 });
      setHistory([]);
    } catch (e) { console.error(e); }
  };
  console.log("Rendering Dashboard with:", { isStreaming, restData, grpcData });
  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-6xl mx-auto">
        <header className="text-center mb-10">
          <h1 className="text-4xl font-bold">📡 IoT Telemetry Benchmark</h1>
          <p className="text-gray-500 mt-2">REST (Polling) vs gRPC-Web (Streaming)</p>
        </header>
        
        <ControlPanel 
          payloadSize={payloadSize} 
          onSizeChange={handleSizeChange} 
          isStreaming={isStreaming} 
          onModeToggle={handleModeToggle} 
        />
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <StatCard 
            title="REST Avg Latency" 
            value={restData.avg} 
            subtitle={"P99: "}
            subValue={restData.p99} 
            unit="µs" 
            borderClass="border-violet-800" 
            textColor="text-violet-800" 
          />
          <StatCard 
            title="gRPC Avg Latency" 
            value={grpcData.avg} 
            subtitle={"P99: "}
            subValue={grpcData.p99} 
            unit="µs" 
            borderClass="border-orange-600" 
            textColor="text-orange-600" 
          />
        </div>
        
        <LatencyChart history={history} />
      </div>
    </div>
  );
};

export default Dashboard;