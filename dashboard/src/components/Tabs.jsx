import React from 'react';

const Tabs = ({ activeTab, setActiveTab }) => {
  const tabs = [
    { id: 'latency', label: 'Latency Analysis' },
    { id: 'payload', label: 'Payload Size'},
  ];

  return (
    <div className="flex justify-center mb-8">
      <div className="inline-flex p-1 bg-slate-200/50 rounded-xl backdrop-blur-sm border border-slate-300/50">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`
              flex items-center px-6 py-2.5 text-sm font-semibold rounded-lg transition-all duration-200
              ${activeTab === tab.id 
                ? 'bg-white text-blue-600 shadow-md transform scale-105' 
                : 'text-slate-600 hover:text-slate-900 hover:bg-white/40'
              }
            `}
          >
            <span className="mr-2">{tab.icon}</span>
            {tab.label}
          </button>
        ))}
      </div>
    </div>
  );
};

export default Tabs;