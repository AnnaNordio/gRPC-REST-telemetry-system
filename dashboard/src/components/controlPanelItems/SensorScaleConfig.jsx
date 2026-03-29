import { useEffect, useState } from 'react';
export const SensorScaleConfig = ({ onSensorChange }) => {
  const [count, setCount] = useState(1);

  useEffect(() => {
    const timer = setTimeout(() => onSensorChange(count), 300);
    return () => clearTimeout(timer);
  }, [count, onSensorChange]);

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-end">
        <span className="text-lg font-black text-slate-900">
          {count} <small className="text-[10px] text-slate-400 uppercase">Sensors</small>
        </span>
      </div>
      <input 
        type="range" min="1" max="100" value={count}
        onChange={(e) => setCount(parseInt(e.target.value))}
        className="w-full h-2 bg-slate-100 rounded-lg appearance-none cursor-pointer accent-blue-600 transition-all"
      />
      <div className="flex justify-between text-[9px] font-bold text-slate-300 uppercase">
        <span>1 Unit</span>
        <span>100 Units</span>
      </div>
    </div>
  );
};