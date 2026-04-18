// Costante per il pareggio
export const COMPARISON_WINNER = {
  REST: 'REST',
  GRPC: 'gRPC',
  DRAW: 'DRAW'
};

const DRAW_STYLE = {
  winner: COMPARISON_WINNER.DRAW,
  color: "text-slate-500",
  bg: "bg-slate-50",
  border: "border-slate-200"
};

export const getPayloadComparison = (restTotalKB, grpcTotalKB) => {
  if (!restTotalKB || !grpcTotalKB) return { ...DRAW_STYLE, text: "Calculating..." };
  
  // Caso Pareggio
  if (restTotalKB === grpcTotalKB) {
    return { ...DRAW_STYLE, text: "Equal Payload Size" };
  }

  const isGrpcLighter = grpcTotalKB < restTotalKB;
  const savedKB = Math.abs(restTotalKB - grpcTotalKB);
  const savedDisplay = savedKB > 1024 ? `${(savedKB / 1024).toFixed(2)} MB` : `${savedKB.toFixed(2)} KB`;

  return {
    winner: isGrpcLighter ? COMPARISON_WINNER.GRPC : COMPARISON_WINNER.REST,
    text: `${isGrpcLighter ? 'gRPC' : 'REST'} saved ${savedDisplay}`,
    color: isGrpcLighter ? "text-orange-600" : "text-violet-600",
    bg: isGrpcLighter ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcLighter ? "border-orange-200" : "border-violet-200"
  };
};

export const getLatencyComparison = (restAvg, grpcAvg) => {
  if (!restAvg || !grpcAvg) return { ...DRAW_STYLE, text: "Analyzing..." };

  if (restAvg === grpcAvg) {
    return { ...DRAW_STYLE, text: "Identical Latency" };
  }

  const isGrpcFaster = grpcAvg < restAvg;
  const diff = Math.abs(restAvg - grpcAvg);
  const percentage = ((diff / Math.max(restAvg, grpcAvg)) * 100).toFixed(2);

  return {
    winner: isGrpcFaster ? COMPARISON_WINNER.GRPC : COMPARISON_WINNER.REST,
    text: `Best Avg. Latency: ${isGrpcFaster ? 'gRPC' : 'REST'} (-${percentage}%)`,
    color: isGrpcFaster ? "text-orange-600" : "text-violet-600",
    bg: isGrpcFaster ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcFaster ? "border-orange-200" : "border-violet-200"
  };
};

export const getThroughputComparison = (restThroughput, grpcThroughput) => {
  if (!restThroughput || !grpcThroughput) return { text: "Benchmarking...", color: "text-slate-400" };

  const isGrpcHigher = grpcThroughput > restThroughput;

  const factor = isGrpcHigher 
    ? (grpcThroughput / restThroughput).toFixed(2) 
    : (restThroughput / grpcThroughput).toFixed(2);

  return {
    winner: isGrpcHigher ? 'gRPC' : 'REST',
    // Il testo ora è molto più "da competizione"
    text: isGrpcHigher 
      ? `gRPC is ${factor}x more scalable` 
      : `REST is ${factor}x more scalable`,
    color: isGrpcHigher ? "text-orange-600" : "text-violet-600",
    bg: isGrpcHigher ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcHigher ? "border-orange-200" : "border-violet-200",
  };
};

export const getMarshalComparison = (restMarshalUs, grpcMarshalUs) => {
  if (!restMarshalUs || !grpcMarshalUs) return { ...DRAW_STYLE, text: "Measuring CPU..." };

  if (Math.abs(restMarshalUs - grpcMarshalUs) < 0.1) {
    return { ...DRAW_STYLE, text: "Similar CPU Usage" };
  }

  const isGrpcFaster = grpcMarshalUs < restMarshalUs;
  const ratio = (Math.max(restMarshalUs, grpcMarshalUs) / Math.min(restMarshalUs, grpcMarshalUs)).toFixed(1);

  return {
    winner: isGrpcFaster ? COMPARISON_WINNER.GRPC : COMPARISON_WINNER.REST,
    text: `CPU Efficiency: ${isGrpcFaster ? 'gRPC' : 'REST'} is ${ratio}x faster`,
    color: isGrpcFaster ? "text-orange-600" : "text-violet-600",
    bg: isGrpcFaster ? "bg-orange-50" : "bg-violet-50",
    border: isGrpcFaster ? "border-orange-200" : "border-violet-200"
  };
};