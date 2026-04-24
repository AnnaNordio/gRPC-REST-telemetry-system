import { useState, useEffect, useRef, useCallback } from 'react';
import * as protos from 'telemetry-proto-bundle';
import { useTelemetryData } from './useTelemetryData';
import { useConnection } from './useConnection';

export const useTelemetry = () => {
  const [history, setHistory] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  const [activeFilter, setActiveFilter] = useState('both');
  const [payloadSize, setPayloadSize] = useState('small');
  const [sensorNumber, setSensorNumber] = useState(1);

  const grpcClient = useRef(new protos.TelemetryServiceClient('http://localhost:8081'));
  const isConnected = useConnection();

  const updateHistory = useCallback((newData) => {
    const pointsToAdd = Array.isArray(newData) ? newData : [newData];
    setHistory(prev => {
      const map = new Map(prev.map(p => [`${p.timestamp}-${p.protocol}`, p]));
      pointsToAdd.forEach(p => p.timestamp && map.set(`${p.timestamp}-${p.protocol}`, p));
      return Array.from(map.values())
        .sort((a, b) => a.timestamp.localeCompare(b.timestamp))
        .slice(-600);
    });
  }, []);

  const { 
    restData, setRestData, grpcData, setGrpcData, 
    handleRestFetch, handleGrpcUnaryFetch 
  } = useTelemetryData(updateHistory);

  // --- LOGICA UNIFICATA POLLING (REST + gRPC Unary) ---
  useEffect(() => {
    const interval = setInterval(() => {
      if (activeFilter === 'both' || activeFilter === 'rest') {
        handleRestFetch();
      }
      
      if (!isStreaming && (activeFilter === 'both' || activeFilter === 'grpc')) {
        handleGrpcUnaryFetch(grpcClient.current, protos);
      }
    }, 500);

    return () => clearInterval(interval);
  }, [activeFilter, isStreaming, handleRestFetch, handleGrpcUnaryFetch]);

  // --- LOGICA STREAMING gRPC ---
  useEffect(() => {
    let stream = null;

    if (isStreaming && (activeFilter === 'both' || activeFilter === 'grpc')) {
      stream = grpcClient.current.getGrpcStream(new protos.Empty(), {});
      
      stream.on('data', (response) => {
        const g = response.toObject();
        setGrpcData({ 
          avg: g.avgLatency, 
          p99: g.p99Latency, 
          payloadSize: g.payloadSize, 
          overheadSize: g.overhead, 
          throughput: g.throughput, 
          marshalTime: g.marshaltime 
        });
        if (g.historyList) updateHistory(g.historyList);
      });

      stream.on('error', (err) => {
        console.error("gRPC Stream Error:", err);
      });
    }

    return () => {
      if (stream) stream.cancel();
    };
  }, [isStreaming, activeFilter, setGrpcData, updateHistory]);

  // --- RESET PROTOCOLLI NON ATTIVI ---
  useEffect(() => {
    if (activeFilter === 'grpc') {
      setRestData({ avg: 0, p99: 0, payloadSize: 0, overheadSize: 0, throughput: 0, marshalTime: 0 });
    } else if (activeFilter === 'rest') {
      setGrpcData({ avg: 0, p99: 0, payloadSize: 0, overheadSize: 0, throughput: 0, marshalTime: 0 });
    }
  }, [activeFilter, setRestData, setGrpcData]);

  // Reset iniziale al mount
  useEffect(() => {
    fetch('http://localhost:8080/reset', { method: 'POST' })
      .catch(err => console.error("Reset error:", err));
  }, []);

  return { 
    restData, grpcData, history, isStreaming, setIsStreaming, 
    payloadSize, setPayloadSize, sensorNumber, setSensorNumber, 
    setHistory, isConnected, activeFilter, setActiveFilter 
  };
};