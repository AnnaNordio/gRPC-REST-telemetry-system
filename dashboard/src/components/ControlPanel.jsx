import React from 'react';

export const ControlPanel = ({ payloadSize, onSizeChange, isStreaming, onModeToggle }) => {
  const toggleActiveColor = "text-slate-900";
  const toggleInactiveColor = "text-slate-300";

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-slate-100 p-8 mb-8 flex flex-col items-center gap-8">
      {/* Payload Size Selector */}
      <div className="flex flex-col items-center gap-3">
        <span className="text-xs font-black uppercase tracking-widest text-slate-600">Payload Size</span>
        <div className="flex flex-wrap justify-center gap-3">
          {['small', 'medium', 'large', 'nested'].map((size) => (
            <button
              key={size}
              onClick={() => onSizeChange(size)}
              className={`px-6 py-2 rounded-full text-xs font-bold uppercase tracking-wider transition-all border ${
                payloadSize === size 
                  ? 'bg-slate-800 text-white border-slate-800 shadow-md' 
                  : 'bg-white text-slate-400 border-slate-200 hover:border-slate-400'
              }`}
            >
              {size}
            </button>
          ))}
        </div>
      </div>

      <div className="w-full max-w-xs border-t border-slate-100"></div>

      {/* Transmission Mode Toggle */}
      <div className="flex flex-col items-center gap-3">
        <span className="text-xs font-black uppercase tracking-widest text-slate-600">Transmission Mode</span>
        <div className="flex items-center gap-4 bg-slate-50 px-6 py-3 rounded-full border border-slate-200">
          <span className={`text-xs font-black uppercase tracking-widest ${!isStreaming ? toggleActiveColor : toggleInactiveColor}`}>Polling</span>
          <label className="relative inline-flex items-center cursor-pointer">
            <input type="checkbox" checked={isStreaming} onChange={onModeToggle} className="sr-only peer" />
            <div className="w-14 h-7 bg-slate-200 rounded-full peer peer-checked:bg-slate-800 after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-6 after:w-6 after:transition-all peer-checked:after:translate-x-full"></div>
          </label>
          <span className={`text-xs font-black uppercase tracking-widest ${isStreaming ? toggleActiveColor : toggleInactiveColor}`}>Streaming</span>
        </div>
      </div>
    </div>
  );
};
