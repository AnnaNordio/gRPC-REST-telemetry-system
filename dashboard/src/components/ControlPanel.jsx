export const ControlPanel = ({ payloadSize, onSizeChange, isStreaming, onModeToggle, isConnected = true }) => {
  return (
    <div className="bg-white rounded-3xl shadow-sm border border-slate-200 p-6 flex flex-col gap-10 h-full">
      
      {/* Sezione Payload */}
      <div className="flex flex-col gap-4">
        <div className="flex items-center gap-2">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <span className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Payload Config</span>
        </div>
        
        <div className="grid grid-cols-2 gap-2">
          {['small', 'medium', 'large', 'nested'].map((size) => (
            <button
              key={size}
              onClick={() => onSizeChange(size)}
              className={`py-3 rounded-xl text-[10px] font-bold uppercase tracking-wider transition-all border ${
                payloadSize === size 
                  ? 'bg-slate-900 text-white border-slate-900 shadow-lg' 
                  : 'bg-slate-50 text-slate-400 border-slate-100 hover:border-slate-300'
              }`}
            >
              {size}
            </button>
          ))}
        </div>
      </div>

      <div className="border-t border-slate-100"></div>

      {/* Sezione Transmission - Ora con flex-1 per occupare lo spazio centrale */}
      <div className="flex flex-col gap-4 flex-1">
        <div className="flex items-center gap-2">
          <div className="w-1 h-4 bg-blue-600 rounded-full"></div>
          <span className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">Network Mode</span>
        </div>

        <div className="flex flex-col gap-3 bg-slate-50 p-4 rounded-3xl border border-slate-100">
          <div className="flex justify-between items-center px-2 mb-2">
            <span className={`text-[10px] font-black uppercase tracking-widest transition-colors ${!isStreaming ? 'text-blue-600' : 'text-slate-300'}`}>
              Polling
            </span>
            <span className={`text-[10px] font-black uppercase tracking-widest transition-colors ${isStreaming ? 'text-slate-900' : 'text-slate-300'}`}>
              Streaming
            </span>
          </div>
          
          <label className="relative w-full h-14 bg-slate-200 rounded-2xl cursor-pointer p-1 transition-colors duration-300 has-[:checked]:bg-slate-900">
            <input 
              type="checkbox" 
              checked={isStreaming} 
              onChange={onModeToggle} 
              className="sr-only peer" 
            />
            
            <div className={`
              absolute top-1 bottom-1 w-[calc(50%-4px)] bg-white rounded-xl shadow-md transition-all duration-500 ease-in-out
              ${isStreaming ? 'left-[calc(50%+2px)]' : 'left-1'}
            `}>
              <div className="flex items-center justify-center h-full">
                <div className={`w-1 h-4 rounded-full ${isStreaming ? 'bg-slate-900' : 'bg-blue-600'} opacity-20`}></div>
              </div>
            </div>
          </label>
        </div>
      </div>
    </div>
  );
};