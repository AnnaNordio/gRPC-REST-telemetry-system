import { StatCard } from '../StatCard';
import { LineChart } from '../charts/LineChart';
import { ComparisonBadge } from '../ComparisonBadge';
import { getComparison } from '../../utils/benchmarkUtils'; // Assicurati che il percorso sia corretto

export const LatencyView = ({ restData, grpcData, history }) => {
  // Calcolo della comparazione direttamente qui dentro
  const latencyComp = getComparison(restData.avg, grpcData.avg, 'latency');

  return (
    <div className="space-y-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      
      {/* Grid delle Statistiche */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <StatCard 
          title="REST Avg Latency" 
          value={restData.avg.toFixed(2)} 
          subtitle="99th Percentile" 
          subValue={`${restData.p99.toFixed(0)} μs`} 
          unit="μs" 
          borderClass="border-violet-600" 
          textColor="text-violet-700" 
        />
        <StatCard 
          title="gRPC Avg Latency" 
          value={grpcData.avg.toFixed(2)} 
          subtitle="99th Percentile" 
          subValue={`${grpcData.p99.toFixed(0)} μs`} 
          unit="μs" 
          borderClass="border-orange-500" 
          textColor="text-orange-600" 
        />
      </div>
      
      {/* Box del Grafico */}
      <div className="bg-white p-6 rounded-3xl shadow-sm border border-slate-200 relative">
        <div className="absolute top-8 right-8 z-10">
          {/* Usiamo la variabile calcolata sopra */}
          <ComparisonBadge data={latencyComp} />
        </div>
        
        <div className="flex items-center gap-2 mb-6">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <h3 className="text-sm font-bold text-slate-700 uppercase tracking-wider">
            Latency Timeline
          </h3>
        </div>
        
        <div>
          <LineChart history={history} measure="Microseconds" unit="μs"/>
        </div>
      </div>
    </div>
  );
};