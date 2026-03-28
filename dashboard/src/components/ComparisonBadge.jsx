export const ComparisonBadge = ({ data }) => {
  if (!data.winner) return (
    <div className="flex items-center gap-2 px-3 py-1.5 rounded-xl border-2 shadow-sm transition-all duration-500 border-slate-100 bg-slate-50">
      <div className="w-2 h-2 rounded-full animate-pulse bg-slate-300"></div>
        <div className="text-right">
          <span className="text-xs font-black text-slate-900">
            Collecting Data...
          </span>
      </div>
    </div>
  );

  return (
    <div className={`flex items-center gap-2 px-3 py-1.5 rounded-xl border-2 shadow-sm transition-all duration-500 ${data.border} ${data.bg}`}>
      <div className={`w-2 h-2 rounded-full animate-pulse ${data.winner === 'gRPC' ? 'bg-orange-500' : 'bg-violet-500'}`}></div>
      <div className="text-right">
        <span className={`text-xs font-black ${data.color}`}>
          {data.text}
        </span>
      </div>
    </div>
  );
};