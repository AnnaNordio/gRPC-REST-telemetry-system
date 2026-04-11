import React from 'react';
import { StatCard } from '../StatCard';
import { ComparisonBadge } from '../ComparisonBadge';
import { getMarshalComparison } from '../../utils/benchmarkUtils';
import { StackedBarChart } from '../charts/StackedBarChart';

export const MarshalView = ({ restData, grpcData, protocol }) => {
  // Calcolo della comparazione CPU
  const marshalComp = getMarshalComparison(restData.marshalTime, grpcData.marshalTime);

  // Calcolo impatto percentuale sulla latenza totale
  const calculateImpact = (marshal, totalMs) => {
    return totalMs > 0 ? ((marshal / totalMs) * 100).toFixed(1) : 0;
  };

  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard 
          title="REST Serialization" 
          value={restData.marshalTime.toFixed(1)} 
          unit="μs"
          subtitle="CPU Weight on Total" 
          subValue={calculateImpact(restData.marshalTime, restData.avg)} 
          subunit="%" 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />

        <StatCard 
          title="gRPC Serialization" 
          value={grpcData.marshalTime.toFixed(1)} 
          unit="μs"
          subtitle="CPU Weight on Total" 
          subValue={calculateImpact(grpcData.marshalTime, grpcData.avg)} 
          subunit="%" 
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>
      
      <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">
          
        <div className="absolute top-6 right-6 z-10">
          <ComparisonBadge data={marshalComp} protocol={protocol} />
        </div>
        <div className="flex items-center gap-2 mb-6">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <h3 className="text-sm font-black text-slate-700 uppercase tracking-wider">
            Computational Cost Analysis
          </h3>
        </div>
        
        <div className="h-[400px]">
          <StackedBarChart restData={restData} grpcData={grpcData} />
        </div>
      </div>
    </div>
  );
};