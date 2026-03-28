import { StatCard } from '../StatCard';
import { GroupedBarChart } from '../charts/GroupedBarChart';

export const PayloadView = ({ restData, grpcData, sizeComp }) => {
  // 1. Calcoliamo tutto in KB una volta sola all'inizio
  const restKB = {
    payload: restData.payloadSize / 1024,
    overhead: restData.overheadSize / 1024,
  };

  const grpcKB = {
    payload: grpcData.payloadSize / 1024,
    overhead: grpcData.overheadSize / 1024,
  };

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Usiamo i valori già convertiti */}
        <StatCard 
          title="REST Payload Size" 
          value={restKB.payload} 
          subtitle="JSON Overhead" 
          subValue={restKB.overhead} 
          unit="KB" 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />
        <StatCard 
          title="gRPC Payload Size" 
          value={grpcKB.payload} 
          subtitle="Protobuf Overhead" 
          subValue={grpcKB.overhead} 
          unit="KB" 
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>

      <div className="bg-white p-8 rounded-3xl shadow-sm border border-slate-200">
        <div className="flex justify-between items-center mb-8">
          <div className="flex items-center gap-2">
            <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
            <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">Weight Analysis</h3>
          </div>
          <div className={`px-4 py-2 rounded-xl border-2 ${sizeComp.border} ${sizeComp.bg} text-right`}>
            <span className="text-[10px] font-black text-slate-400 uppercase tracking-widest block leading-none mb-1">Efficiency</span>
            <span className={`font-bold ${sizeComp.color}`}>{sizeComp.label}</span>
          </div>
        </div>
        
        <div className="h-[400px]">
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