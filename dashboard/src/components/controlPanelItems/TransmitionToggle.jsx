export const TransmissionToggle = ({ isStreaming, onModeToggle }) => (
  <div className="flex flex-col gap-3 bg-slate-50 p-4 rounded-3xl border border-slate-100">
    <div className="flex justify-between items-center px-2 mb-2">
      <span className={`text-[10px] font-black uppercase tracking-widest transition-colors ${!isStreaming ? 'text-blue-600' : 'text-slate-300'}`}>
        Unary (Req-Res)
      </span>
      <span className={`text-[10px] font-black uppercase tracking-widest transition-colors ${isStreaming ? 'text-slate-900' : 'text-slate-300'}`}>
        Stream
      </span>
    </div>
    <label className="relative w-full h-14 bg-blue-600 rounded-2xl cursor-pointer p-1 transition-colors duration-300 has-[:checked]:bg-slate-900">
      <input type="checkbox" checked={isStreaming} onChange={onModeToggle} className="sr-only peer" />
      <div className={`absolute top-1 bottom-1 w-[calc(50%-4px)] bg-white rounded-xl shadow-md transition-all duration-500 ease-in-out ${isStreaming ? 'left-[calc(50%+2px)]' : 'left-1'}`}>
        <div className="flex items-center justify-center h-full">
          <div className={`w-1 h-4 rounded-full ${isStreaming ? 'bg-slate-900' : 'bg-blue-600'} opacity-20`}></div>
        </div>
      </div>
    </label>
  </div>
);