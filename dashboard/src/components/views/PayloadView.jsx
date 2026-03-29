import { StatCard } from '../StatCard';
import { GroupedBarChart } from '../charts/GroupedBarChart';
import { ComparisonBadge } from '../ComparisonBadge';
import { getPayloadComparison } from '../../utils/benchmarkUtils';

export const PayloadView = ({ restData, grpcData }) => {
  // 1. Calcoliamo tutto in KB una volta sola all'inizio
  const restKB = {
    payload: restData.payloadSize / 1024,
    overhead: restData.overheadSize / 1024,
    total: (restData.payloadSize + restData.overheadSize) / 1024,
  };

  const grpcKB = {
    payload: grpcData.payloadSize / 1024,
    overhead: grpcData.overheadSize / 1024,
    total: (grpcData.payloadSize + grpcData.overheadSize) / 1024,
  };

  // Calcolo Efficienza REST
  const restEfficiency = restKB.payload > 0 
    ? ((restKB.payload / (restKB.payload + restKB.overhead)) * 100).toFixed(2) 
    : 0;

  // Calcolo Efficienza gRPC
  const grpcEfficiency = grpcKB.payload > 0 
    ? ((grpcKB.payload / (grpcKB.payload + grpcKB.overhead)) * 100).toFixed(2) 
    : 0;

const comparison = getPayloadComparison(restKB.total, grpcKB.total);  return (
    <div className="space-y-4 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard 
          title="REST Payload Size" 
          value={restKB.payload > 1024  ?(restKB.payload / 1024).toFixed(2) :restKB.payload.toFixed(2)} 
          subtitle="JSON Overhead" 
          subValue={restKB.overhead > 1024  ?(restKB.overhead / 1024).toFixed(2) :restKB.overhead.toFixed(2)} 
          addedTitle={"Efficiency"}
          addedValue={restEfficiency}
          addedUnit={"%"}
          unit={restKB.payload > 1024  ?"MB" :"KB"}  
          subunit={restKB.overhead > 1024  ?"MB" :"KB"}  
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />
        <StatCard 
          title="gRPC Payload Size" 
          value={grpcKB.payload > 1024  ?(grpcKB.payload / 1024).toFixed(2) :grpcKB.payload.toFixed(2)} 
          subtitle="Protobuf Overhead" 
          subValue={grpcKB.overhead > 1024  ?(grpcKB.overhead / 1024).toFixed(2) :grpcKB.overhead.toFixed(2)} 
          addedTitle={"Efficiency"}
          addedValue={grpcEfficiency}
          addedUnit={"%"}
          unit={grpcKB.payload > 1024  ?"MB" :"KB"}  
          subunit={grpcKB.overhead > 1024  ?"MB" :"KB"}  
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
          <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">Payload Size Analysis</h3>
        </div>
        
         <div>
          <GroupedBarChart 
            restSize={restKB.payload} 
            restOverhead={restKB.overhead} 
            grpcSize={grpcKB.payload} 
            grpcOverhead={grpcKB.overhead} 
            unit="KB"
          />
        </div>
      </div>       
    </div>
  );
};