import { StatCard } from '../StatCard';
import { LineChart } from '../charts/LineChart';
import { ComparisonBadge } from '../ComparisonBadge';

export const LatencyView = ({ restData, grpcData, history, comparison }) => (
  <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
      <StatCard title="REST Avg Latency" value={restData.avg} subtitle="99th Percentile" subValue={restData.p99} unit="μs" borderClass="border-violet-600" textColor="text-violet-700" />
      <StatCard title="gRPC Avg Latency" value={grpcData.avg} subtitle="99th Percentile" subValue={grpcData.p99} unit="μs" borderClass="border-orange-500" textColor="text-orange-600" />
    </div>
    
    <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">
      <div className="absolute top-6 right-6 z-10">
        <ComparisonBadge data={comparison} />
      </div>
      
      <div className="flex items-center gap-2 mb-6">
        <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
        <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">Latency Timeline</h3>
      </div>
      
      <div className="h-[400px]">
        <LineChart history={history} measure="Microseconds" unit="μs"/>
      </div>
    </div>
  </div>
);