export const ControlSection = ({ title, iconColor = "bg-blue-600", children, showSeparator = true }) => (
  <>
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <div className={`w-1 h-4 ${iconColor} rounded-full`}></div>
        <span className="text-[10px] font-black uppercase tracking-[0.2em] text-slate-400">
          {title}
        </span>
      </div>
      {children}
    </div>
    {showSeparator && <div className="border-t border-slate-100"></div>}
  </>
);