export const StatCard = ({ 
  title, 
  subtitle,   
  value,       
  subValue,   
  percentageTitle,
  percentageValue,
  unit,
  borderClass, 
  textColor 
}) => {
  const formatNum = (num) => {
    const n = parseFloat(num);
    return !isNaN(n) ? n.toFixed(1) : '0.0';
  };

  const hasValue = (val) => val !== undefined && val !== null && !isNaN(parseFloat(val));
  console.log(unit)
  return (
    <div className={`relative bg-white p-6 rounded-2xl shadow-sm border-l-[10px] ${borderClass} transition-all hover:shadow-md overflow-hidden`}>
      
      {/* Badge Percentuale - Più visibile e spaziato */}
      {percentageTitle && hasValue(percentageValue) && (
        <div className={`absolute top-5 right-5 px-3 py-2 rounded-xl bg-slate-50 border border-slate-100 flex flex-col items-end shadow-sm`}>
          <span className="text-[12px] uppercase font-bold text-slate-400 leading-none mb-1.5 tracking-wider">
            {percentageTitle}
          </span>
          <span className={`text-base font-black ${textColor}`}>
            {formatNum(percentageValue)}%
          </span>
        </div>
      )}

      {/* Titolo principale: passato a text-sm (14px) */}
      <h3 className={`text-sm font-bold tracking-widest uppercase opacity-80 ${textColor}`}>
        {title}
      </h3>
      
      {/* Valore principale: text-5xl per un impatto massiccio */}
      <div className={`text-4xl font-black ${textColor} mt-2 flex items-baseline tracking-tighter`}>
        {hasValue(value) ? formatNum(value) : '--'}
        <span className="text-xl ml-2 font-bold opacity-30">{unit}</span>
      </div>
      
      {subtitle && (
        <div className={`text-sm mt-4 flex items-center gap-2 font-bold opacity-90 ${textColor}`}>
          <span className="uppercase tracking-wide">{subtitle}:</span>
          <span className="px-1 py-0.5">
            {hasValue(subValue) ? `${formatNum(subValue)} ${unit}` : '--'}
          </span>
        </div>
      )}
    </div>
  );
};