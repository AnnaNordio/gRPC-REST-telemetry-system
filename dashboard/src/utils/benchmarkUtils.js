/**
 * Calcola il confronto prestazionale tra REST e gRPC.
 * @param {number} restValue - Value of REST
 * @param {number} grpcValue - Value of gRPC
 * @param {string} type - benchmark type
 * @returns {object} ui-friendly comparison result with text, label, color, and winner indication
 */
export const getComparison = (restValue, grpcValue, type = 'latency') => {
  if (!restValue || !grpcValue) {
    return { 
      text: "Waiting for data...", 
      label: "N/D",
      color: "text-slate-400", 
      bg: "bg-slate-100", 
      border: "border-slate-200", 
      winner: null 
    };
  }

  const diff = restValue - grpcValue;
  const isGrpcBetter = diff > 0;
  const absDiff = Math.abs(diff);
  
  // Percentuale rispetto al valore più alto (il più lento/pesante)
  const percentage = ((absDiff / (isGrpcBetter ? restValue : grpcValue)) * 100).toFixed(1);
  const unitText = type === 'latency' ? 'faster' : 'lighter';

  return isGrpcBetter ? {
    text: `gRPC is ${percentage}% ${unitText} than REST`,
    label: `gRPC +${percentage}%`,
    color: "text-orange-600",
    bg: "bg-orange-50",
    border: "border-orange-200",
    winner: 'gRPC'
  } : {
    text: `REST is ${percentage}% ${unitText} than gRPC`,
    label: `REST +${percentage}%`,
    color: "text-violet-600",
    bg: "bg-violet-50",
    border: "border-violet-200",
    winner: 'REST'
  };
};