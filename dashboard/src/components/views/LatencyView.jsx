import React from 'react';
import { StatCard } from '../StatCard';
import { LineChart } from '../charts/LineChart';
import { ComparisonBadge } from '../ComparisonBadge';
import { getLatencyComparison } from '../../utils/benchmarkUtils';

export const LatencyView = ({ restData, grpcData, history }) => {
  const comparison = getLatencyComparison(restData.avg, grpcData.avg);
  console.log("LatencyView - Rest Data:", restData);
  console.log("LatencyView - gRPC Data:", grpcData);
  console.log("LatencyView - History Data:", history);
  const formatVal = (val) => (val > 1000 ? (val / 1000).toFixed(2) : val.toFixed(2));
  const getUnit = (val) => (val > 1000 ? "ms" : "μs");

  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* REST CARD */}
        <StatCard 
          title="REST Latency (Avg)" 
          value={formatVal(restData.avg)} 
          unit={getUnit(restData.avg)}
          subtitle="99th Percentile (Tail)" 
          subValue={formatVal(restData.p99)} 
          subunit={getUnit(restData.p99)} 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />

        {/* GRPC CARD */}
        <StatCard 
          title="gRPC Latency (Avg)" 
          value={formatVal(grpcData.avg)} 
          unit={getUnit(grpcData.avg)}
          subtitle="99th Percentile (Tail)" 
          subValue={formatVal(grpcData.p99)} 
          subunit={getUnit(grpcData.p99)} 
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
            Latency & Reliability Timeline (Avg vs P99)
          </h3>
        </div>
        
        <div className="h-[400px]">
          <LineChart history={history} unit={'μs'} />
        </div>
      </div>
    </div>
  );
};