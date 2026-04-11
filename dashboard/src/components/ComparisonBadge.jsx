export const ComparisonBadge = ({ data, protocol }) => {
  if (!data.winner) return (
    <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl border-2 shadow-sm transition-all duration-500 border-slate-100 bg-slate-50">
      <div className="w-2 h-2 rounded-full animate-pulse bg-slate-500"></div>
        <div className="text-right">
          <span className="text-xs font-black text-slate-900">
            Initializing Benchmark...
          </span>
      </div>
    </div>
  );

  if (protocol !== 'both') return null;

  const dotColor = data.winner === 'DRAW' 
    ? 'bg-slate-500'
    : (data.winner === 'gRPC' ? 'bg-orange-500' : 'bg-violet-500');

  return (
    <div className={`flex items-center gap-2 px-3 py-1.5 rounded-xl border-2 shadow-sm transition-all duration-500 ${data.border} ${data.bg}`}>
      <div className={`w-2 h-2 rounded-full animate-pulse ${dotColor}`}></div>
      <div className="text-right">
        <span className={`text-xs font-black ${data.color}`}>
          {data.text}
        </span>
      </div>
    </div>
  );
};