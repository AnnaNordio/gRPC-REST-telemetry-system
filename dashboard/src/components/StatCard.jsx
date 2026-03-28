export const StatCard = ({ 
  title, 
  subtitle,   
  value,       
  subValue,   
  percentageTitle,
  percentageValue,
  unit = "µs", 
  borderClass, 
  textColor 
}) => {
  // Funzione di formattazione più robusta
  const formatNum = (num) => {
    const n = parseFloat(num);
    return !isNaN(n) ? n.toFixed(2) : '0.00';
  };

  // Verifica se il valore esiste ed è un numero valido (anche se è 0)
  const hasValue = (val) => val !== undefined && val !== null && !isNaN(parseFloat(val));

  return (
    <div className={`bg-white p-6 rounded-2xl shadow-sm border-l-8 ${borderClass} transition-all hover:shadow-md`}>
      <h3 className="text-sm font-semibold text-gray-400 tracking-widest uppercase">
        {title}
      </h3>
      
      <div className={`text-4xl font-black ${textColor} mt-2`}>
        {/* Cambiato il controllo: se il valore esiste, mostralo (anche se è 0) */}
        {hasValue(value) ? formatNum(value) : '--'}
        <span className="text-lg ml-1 font-medium opacity-60">{unit}</span>
      </div>
      
      {subtitle && (
        <div className={`text-sm opacity-70 ${textColor} mt-1`}>
          <span className="font-bold">{subtitle}</span>
          {hasValue(subValue) ? `: ${formatNum(subValue)} ${unit}` : ' --'}
        </div>
      )}

      {percentageTitle && (
        <div className={`text-sm opacity-70 ${textColor} mt-1`}>
          <span className="font-bold">{percentageTitle}</span>
          {hasValue(percentageValue) ? `: ${formatNum(percentageValue)}%` : ' --'}
        </div>
      )}
    </div>
  );
};