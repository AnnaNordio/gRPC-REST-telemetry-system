import React from 'react';
import { StatCard } from '../StatCard';
import { BarChart } from '../charts/BarChart';
import { ComparisonBadge } from '../ComparisonBadge';
import { getThroughputComparison } from '../../utils/benchmarkUtils';

export const ThroughputView = ({ restData, grpcData, protocol, sensorNumber }) => {
  const comparison = getThroughputComparison(restData.throughput, grpcData.throughput);
  
  const n = sensorNumber || 1;
  const restEff = (restData.throughput / n).toFixed(1);
  const grpcEff = (grpcData.throughput / n).toFixed(1);

  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard 
          title="REST Avg Efficiency" 
          value={restEff} 
          unit="msg/node" 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />
        <StatCard 
          title="gRPC Avg Efficiency" 
          value={grpcEff} 
          unit="msg/node" 
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>
      
      <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">

          <div className="absolute top-6 right-6 z-10">
            <ComparisonBadge data={comparison} protocol={protocol} />
          </div>
          
          <div className="flex items-center gap-2 mb-6">
            <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
            <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">
              Total System Throughput (Stress Test)
            </h3>
          </div>
          
          <div>
            <BarChart 
              restValue={restData.throughput} 
              grpcValue={grpcData.throughput} 
              measure="Messages per Second" 
              unit="msg/s"
            />
          </div>
      </div>
    </div>
  );
};