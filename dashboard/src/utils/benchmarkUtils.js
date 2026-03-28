export const getPayloadComparison = (restTotalKB, grpcTotalKB) => {
  // Controlliamo che i dati esistano e non siano identici
  if (!restTotalKB || !grpcTotalKB || restTotalKB === grpcTotalKB) {
    return { 
      text: "Calculating savings...", 
      winner: null, 
      color: "text-slate-400", 
      bg: "bg-slate-50", 
      border: "border-slate-100" 
    };
  }

  const isGrpcLighter = grpcTotalKB < restTotalKB;
  const savedKB = Math.abs(restTotalKB - grpcTotalKB);
  
  // Se il risparmio supera i 1024 KB, mostriamo MB, altrimenti KB
  const savedDisplay = savedKB > 1024 
    ? `${(savedKB / 1024).toFixed(2)} MB` 
    : `${savedKB.toFixed(2)} KB`;

  return {
    winner: isGrpcLighter ? 'gRPC' : 'REST',
    text: `${isGrpcLighter ? 'gRPC' : 'REST'} saved ${savedDisplay}`,
    color: isGrpcLighter ? "text-orange-600" : "text-violet-600",
    bg: isGrpcLighter ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcLighter ? "border-orange-200" : "border-violet-200"
  };
};

export const getLatencyComparison = (restAvg, grpcAvg) => {
  if (!restAvg || !grpcAvg || restAvg === grpcAvg) {
    return { text: "Analyzing latency...", winner: null, color: "text-slate-400", bg: "bg-slate-50", border: "border-slate-100" };
  }

  const isGrpcFaster = grpcAvg < restAvg;
  const diff = Math.abs(restAvg - grpcAvg);
  const maxVal = Math.max(restAvg, grpcAvg);
  const percentage = ((diff / maxVal) * 100).toFixed(2);

  return {
    winner: isGrpcFaster ? 'gRPC' : 'REST',
    text: `${isGrpcFaster ? 'gRPC' : 'REST'} is ${percentage}% faster`,
    color: isGrpcFaster ? "text-orange-600" : "text-violet-600",
    bg: isGrpcFaster ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcFaster ? "border-orange-200" : "border-violet-200"
  };
};