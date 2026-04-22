import React from 'react';

export const ProtocolFilter = ({ activeFilter, onFilterChange }) => {
  const options = [
    { id: 'rest', label: 'Rest Only'},
    { id: 'both', label: 'Dual Stack'},
    { id: 'grpc', label: 'gRPC Only'}
  ];

  return (
    <div className="flex p-1.5 bg-slate-100/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 shadow-inner">
      {options.map((option) => (
        <button
          key={option.id}
          onClick={() => onFilterChange(option.id)}
          className={`
            flex-1 flex items-center justify-center gap-1.5 
            py-2 px-3 rounded-xl text-[10px] font-black uppercase tracking-tighter 
            transition-all duration-200 ease-out
            ${activeFilter === option.id 
              ? 'bg-white text-blue-600 shadow-md ring-1 ring-slate-200/50 scale-[1.02]' 
              : 'text-slate-400 hover:text-slate-600 hover:bg-white/50'
            }
          `}
        >
          {option.label}
        </button>
      ))}
    </div>
  );
};