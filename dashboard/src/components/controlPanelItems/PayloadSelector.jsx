export const PayloadSelector = ({ currentSize, onSizeChange }) => {
  const sizes = ['small', 'medium', 'large', 'nested'];
  return (
    <div className="grid grid-cols-2 gap-2">
      {sizes.map((size) => (
        <button
          key={size}
          onClick={() => onSizeChange(size)}
          className={`py-3 rounded-xl text-[10px] font-bold uppercase tracking-wider transition-all border ${
            currentSize === size 
              ? 'bg-slate-900 text-white border-slate-900 shadow-lg' 
              : 'bg-slate-50 text-slate-400 border-slate-100 hover:border-slate-300'
          }`}
        >
          {size}
        </button>
      ))}
    </div>
  );
};