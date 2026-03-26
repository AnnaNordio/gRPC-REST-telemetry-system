import React from 'react';

export const StatCard = ({ 
  title, 
  subtitle,   
  value,       
  subValue,   
  unit = "µs", 
  borderClass, 
  textColor 
}) => {
  const formatNum = (num) => (typeof num === 'number' ? num.toFixed(2) : '0.00');

  return (
    <div className={`bg-white p-6 rounded-2xl shadow-sm border-l-8 ${borderClass} transition-all hover:shadow-md`}>
      <h3 className="text-sm font-semibold text-gray-400 tracking-widest uppercase">
        {title}
      </h3>
      
      <div className={`text-4xl font-black ${textColor} mt-2`}>
        {formatNum(value)}
        <span className="text-lg ml-1 font-medium opacity-60">{unit}</span>
      </div>
      
      {subtitle && (
        <div className={`text-sm opacity-70 ${textColor} mt-1`}>
          <span className="font-bold">{subtitle}</span>
          {subValue !== undefined && `: ${formatNum(subValue)} ${unit}`}
        </div>
      )}
    </div>
  );
};