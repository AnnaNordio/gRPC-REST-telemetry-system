import { useEffect, useState } from 'react';

export const SensorScaleConfig = ({ onSensorChange }) => {
  const STEPS = [1, 10, 50, 100];
  
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const timer = setTimeout(() => onSensorChange(STEPS[index]), 300);
    return () => clearTimeout(timer);
  }, [index, onSensorChange]);

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-end">
        <span className="text-lg font-black text-slate-900">
          {STEPS[index]} <small className="text-[10px] text-slate-400 uppercase">Sensor{STEPS[index] !== 1 ? 's' : ''}</small>
        </span>
      </div>
      
      <input 
        type="range" 
        min="0" 
        max={STEPS.length - 1}
        step="1"
        value={index}
        onChange={(e) => setIndex(parseInt(e.target.value))}
        className="w-full h-2 bg-slate-100 rounded-lg appearance-none cursor-pointer accent-blue-600 transition-all"
      />
      
      <div className="flex justify-between text-[10px] font-bold text-slate-400 uppercase px-1">
        {STEPS.map((val) => (
          <span key={val} className={STEPS[index] === val ? "text-blue-600" : ""}>
            {val}
          </span>
        ))}
      </div>
    </div>
  );
};