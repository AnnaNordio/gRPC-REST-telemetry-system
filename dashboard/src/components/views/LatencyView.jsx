import React from 'react';
import { StatCard } from '../StatCard';
import { LineChart } from '../charts/LineChart';
import { ComparisonBadge } from '../ComparisonBadge';
import { getLatencyComparison } from '../../utils/benchmarkUtils';

export const LatencyView = ({ restData, grpcData, history }) => {
  const comparison = getLatencyComparison(restData.avg, grpcData.avg);
  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard 
          title="REST Average Aggregate Latency" 
          value={restData.avg > 1000  ?(restData.avg / 1000).toFixed(2) :restData.avg.toFixed(2)} 
          subtitle="99th Percentile (Tail Latency)" 
          subValue={restData.p99 > 1000  ?(restData.p99 / 1000).toFixed(2) :restData.p99.toFixed(2)} 
          unit={restData.avg > 1000  ? "ms" : "μs"} 
          subunit={restData.p99 > 1000  ? "ms" : "μs"} 
          addedTitle="Jitter"
          addedValue={restData.jitter > 1000  ?(restData.jitter / 1000) :restData.jitter} 
          addedUnit={restData.jitter > 1000  ? "ms" : "μs"} 

          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />
        <StatCard 
          title="gRPC Average Aggregate Latency" 
          value={grpcData.avg > 1000  ?(grpcData.avg / 1000).toFixed(2) :grpcData.avg.toFixed(2)} 
          subtitle="99th Percentile (Tail Latency)" 
          subValue={grpcData.p99 > 1000  ?(grpcData.p99 / 1000).toFixed(2) :grpcData.p99.toFixed(2)} 
          unit={grpcData.avg > 1000  ? "ms" : "μs"} 
          subunit={grpcData.p99 > 1000  ? "ms" : "μs"} 
          addedTitle="Jitter"
          addedValue={grpcData.jitter > 1000  ?(grpcData.jitter / 1000):grpcData.jitter} 
          addedUnit={grpcData.jitter > 1000  ? "ms" : "μs"} 
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>
      
      <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">
        <div className="absolute top-6 right-6 z-10">
          <ComparisonBadge data={comparison} />
        </div>
        
        <div className="flex items-center gap-2 mb-6">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">
            Instantaneous Latency Timeline (All Sensors)
          </h3>
        </div>
        
        <div>
          <LineChart history={history} measure="Microseconds" unit="μs" />
        </div>
      </div>
    </div>
  );
};