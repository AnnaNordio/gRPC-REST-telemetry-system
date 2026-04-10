import React from 'react';

const Tabs = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'latency', label: 'Latency & Reliability' },
    { id: 'marshalling', label: 'Computational Cost' },
    { id: 'payload', label: 'Payload & Efficiency' },
    { id: 'scalability', label: 'Scalability & Stress Test' },
  ];

  return (
    <div className="w-full mb-4"> 
      <div className="flex w-full bg-slate-200/50 rounded-xl backdrop-blur-sm shadow-sm">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`
              flex-1 flex items-center justify-center px-6 py-2.5 text-sm font-semibold rounded-lg transition-all duration-200
              ${activeTab === tab.id 
                ? 'bg-white text-blue-600 shadow-md transform scale-[1.02]'
                : 'text-slate-600 hover:text-slate-900 hover:bg-white/40'
              }
            `}
          >
            {tab.icon && <span className="mr-2">{tab.icon}</span>}
            {tab.label}
          </button>
        ))}
      </div>
    </div>
  );
};

export default Tabs;