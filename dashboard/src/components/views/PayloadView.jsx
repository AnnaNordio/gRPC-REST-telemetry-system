import { getComparison } from '../../utils/benchmarkUtils'; 
import { StatCard } from '../StatCard';
import { GroupedBarChart } from '../charts/GroupedBarChart';
import { ComparisonBadge } from '../ComparisonBadge';

export const PayloadView = ({ restData, grpcData }) => {
  const totalRestKB = (restData.totalSize + restData.totalOverhead) / 1024;
  const totalGrpcKB = (grpcData.totalSize + grpcData.totalOverhead) / 1024;
  const sizeComp = getComparison(totalRestKB, totalGrpcKB, 'size');

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        
        {/* Card per REST */}
        <StatCard 
          title="REST Total Traffic" 
          value={restData.payloadSize.toFixed(2)} 
          unit="KB"
          subtitle="Total JSON Overhead" 
          subValue={`${(restData.overhead / 1024).toFixed(2)} KB`} 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />

        {/* Card per gRPC */}
        <StatCard 
          title="gRPC Total Traffic" 
          value={grpcData.payloadSize.toFixed(2)} 
          unit="KB"
          subtitle="Total Proto Overhead" 
          subValue={(grpcData.overhead / 1024).toFixed(2)}
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>

      {/* Il grafico mostrerà ora le "montagne" di dati accumulati */}
      <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">
        <div className="absolute top-8 right-8 z-10">
          <ComparisonBadge data={sizeComp} />
        </div>
        <div className="flex items-center gap-2 mb-6">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">
            Data Volume Comparison
          </h3>
        </div>
        <div>
          <GroupedBarChart 
            restSize={restData.payloadSize} 
            restOverhead={restData.overhead} 
            grpcSize={grpcData.payloadSize} 
            grpcOverhead={grpcData.overhead} 
            unit="KB"
          />
        </div>
      </div>
    </div>
  );
};