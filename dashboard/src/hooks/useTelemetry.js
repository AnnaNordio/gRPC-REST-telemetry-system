import { useState, useEffect, useRef, useCallback } from 'react';
import * as protos from 'telemetry-proto-bundle';
import { formatTimestamp } from '../utils/formatters';

export const useTelemetry = () => {
  const [restData, setRestData] = useState({ avg: 0, p99: 0 });
  const [grpcData, setGrpcData] = useState({ avg: 0, p99: 0 });
  const [history, setHistory] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [isConnected, setIsConnected] = useState(false);
  const [payloadSize, setPayloadSize] = useState('small');
  const [sensorNumber, setSensorNumber] = useState(1);
  const [activeFilter, setActiveFilter] = useState('both');

  const grpcClient = useRef(new protos.TelemetryServiceClient('http://localhost:8081'));

  const updateHistory = useCallback((newData) => {
    const pointsToAdd = Array.isArray(newData) ? newData : [newData];
    setHistory(prev => {
      const map = new Map();
      prev.forEach(p => map.set(`${p.timestamp}-${p.protocol}`, p));
      pointsToAdd.forEach(p => {
        if (p.timestamp) map.set(`${p.timestamp}-${p.protocol}`, p);
      });
      return Array.from(map.values())
        .sort((a, b) => a.timestamp.localeCompare(b.timestamp))
        .slice(-600);
    });
  }, []);

  const fetchData = useCallback(async () => {
    // REST
    try {
      const response = await fetch('http://localhost:8080/results');
      const data = await response.json();
      setRestData({ avg: data.avg_rest, p99: data.p99_rest, payloadSize: data.payload_size, overheadSize: data.overhead_size, throughput: data.throughput_rest, marshalTime: data.marshal_time_rest }); 
      if (data.history) updateHistory(data.history);
    } catch (err) { console.error("REST Error:", err); }

    // gRPC Unary
    if (!isStreaming) {
      grpcClient.current.getStats(new protos.Empty(), {}, (err, response) => {
        if (!err && response) {
          const g = response.toObject();
          const ts = g.timestamp || g.Timestamp || 0;
          const syncTime = formatTimestamp(ts);
          setGrpcData({ avg: g.avgLatency, p99: g.p99Latency, payloadSize: g.payloadSize, overheadSize: g.overhead, throughput: g.throughput, marshalTime: g.marshaltime });
          if (syncTime) updateHistory({ timestamp: syncTime, protocol: 'gRPC', latency_ms: g.avgLatency, p99: g.p99Latency, marshal: g.marshaltime});
        }
      });
    }
  }, [isStreaming, updateHistory]);

  useEffect(() => {
    let stream = null;
    if (isStreaming) {
      stream = grpcClient.current.getGrpcStream(new protos.Empty(), {});
      stream.on('data', (response) => {
        const g = response.toObject();
        const syncTime = formatTimestamp(g.timestamp);
        setGrpcData({ avg: g.avgLatency, p99: g.p99Latency, payloadSize: g.payloadSize, overheadSize: g.overhead, throughput: g.throughput, marshalTime: g.marshaltime });
        if (syncTime) updateHistory({ timestamp: syncTime, protocol: 'gRPC', latency_ms: g.avgLatency, p99: g.p99Latency, marshal: g.marshaltime});
      });
    }
    const interval = setInterval(fetchData, 500);
    return () => { if (stream) stream.cancel(); clearInterval(interval); };
  }, [isStreaming, fetchData, updateHistory]);

  useEffect(() => {
    let socket;
    let reconnectTimeout;
    let isMounted = true; // Per evitare aggiornamenti di stato su componenti smontati

    // Cerca la funzione connect() dentro il tuo useEffect in useTelemetry.js
    const connect = () => {
      // Invece di usare window.location.host (che sarebbe localhost:3000)
      // puntiamo direttamente al backend sulla 8080
      const wsUrl = `ws://${window.location.hostname}:8080/ws`;
      
      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        if (isMounted) setIsConnected(true);
      };

      socket.onclose = () => {
        if (isMounted) {
          setIsConnected(false);
          reconnectTimeout = setTimeout(connect, 5000);
        }
      };

      socket.onerror = () => {
        socket.close();
      };
    };

    const initialTimeout = setTimeout(connect, 500);

    return () => {
      isMounted = false;
      clearTimeout(initialTimeout);
      clearTimeout(reconnectTimeout);
      if (socket) socket.close();
    };
  }, []);

  useEffect(() => {
    const resetBackendAndLocal = async () => {
      try {
        // Avvisa il backend di resettare le statistiche globali
        await fetch('http://localhost:8080/reset', { method: 'POST' });
        
        // Resetta lo stato locale per partire da un grafico pulito
        setHistory([]);
        setRestData({ avg: 0, p99: 0 });
        setGrpcData({ avg: 0, p99: 0 });
        setSensorNumber(1);
        setActiveFilter('both');
      } catch (err) {
        console.error("Errore durante il reset iniziale:", err);
      }
    };
    resetBackendAndLocal();
  }, []);

  return { restData, grpcData, history, isStreaming, setIsStreaming, payloadSize, setPayloadSize, sensorNumber, setSensorNumber, setHistory, isConnected, activeFilter, setActiveFilter };
};