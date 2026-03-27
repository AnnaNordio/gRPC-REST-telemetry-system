import { StatCard } from '../StatCard';
import { GroupedBarChart } from '../charts/GroupedBarChart';

export const PayloadView = ({ restData, grpcData, sizeComp }) => (
  <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <StatCard title="REST Total Size" value={restData.size + (restData.overhead / 1024)} subtitle="JSON Overhead" subValue={`${restData.overhead} B`} unit="KB" borderClass="border-violet-600" textColor="text-violet-700" />
      <StatCard title="gRPC Total Size" value={grpcData.size + (grpcData.overhead / 1024)} subtitle="Protobuf Overhead" subValue={`${grpcData.overhead} B`} unit="KB" borderClass="border-orange-500" textColor="text-orange-600" />
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
          restSize={restData.size} 
          restOverhead={restData.overhead / 1024} 
          grpcSize={grpcData.size} 
          grpcOverhead={grpcData.overhead / 1024} 
          unit="KB"
        />
      </div>
    </div>
  </div>
);