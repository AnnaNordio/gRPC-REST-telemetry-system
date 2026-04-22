import { useState, useCallback } from 'react';

export const useTelemetryData = (updateHistory) => {
  const [restData, setRestData] = useState({ avg: 0, p99: 0 });
  const [grpcData, setGrpcData] = useState({ avg: 0, p99: 0 });

  const handleRestFetch = useCallback(async () => {
    try {
      const response = await fetch('http://localhost:8080/results');
      const data = await response.json();
      setRestData({ 
        avg: data.avg_rest, p99: data.p99_rest, 
        payloadSize: data.payload_size, overheadSize: data.overhead_size, 
        throughput: data.throughput_rest, marshalTime: data.marshal_time_rest 
      }); 
      if (data.history) updateHistory(data.history);
    } catch (err) { console.error("REST Error:", err); }
  }, [updateHistory]);

  const handleGrpcUnaryFetch = useCallback((grpcClient, protos) => {
    grpcClient.getStats(new protos.Empty(), {}, (err, response) => {
      if (!err && response) {
        const g = response.toObject();
        setGrpcData({ 
          avg: g.avgLatency, p99: g.p99Latency, 
          payloadSize: g.payloadSize, overheadSize: g.overhead, 
          throughput: g.throughput, marshalTime: g.marshaltime 
        });
        if (g.historyList) updateHistory(g.historyList);
      }
    });
  }, [updateHistory]);

  return { restData, setRestData, grpcData, setGrpcData, handleRestFetch, handleGrpcUnaryFetch };
};