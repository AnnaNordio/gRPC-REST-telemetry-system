import React from 'react';

export const ConnectionCard = ({ isConnected }) => {
  return (
    <div className="p-4 bg-blue-600 rounded-2xl text-white shadow-lg shadow-blue-200 transition-all duration-300">
      <h4 className="text-xs font-bold uppercase tracking-widest opacity-80">
        System Status
      </h4>
      <div className="flex items-center gap-2 mt-1">
        <span 
          className={`h-3 w-3 rounded-full flex-shrink-0 transition-all duration-500 ${
            isConnected 
              ? "bg-green-400 shadow-[0_0_10px_rgba(74,222,128,0.9)] animate-pulse" 
              : "bg-red-400 shadow-[0_0_10px_rgba(248,113,113,0.9)]"
          }`}
        ></span>
        
        <p className="text-sm font-medium leading-none tracking-tight">
          {isConnected ? "Gateway Online" : "Gateway Offline"}
        </p>
      </div>
    </div>
  );
};