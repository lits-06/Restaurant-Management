import React from 'react';

export interface KPIItem {
  title: string;
  value: string;
  trend: string;
  isUp: boolean | null;
  isPrimary: boolean;
}

interface KPIGridProps {
  kpis?: KPIItem[];
}

const defaultKpis: KPIItem[] = [
    {
      title: 'Total Revenue',
      value: '$142,580.00',
      trend: '12% from last month',
      isUp: true,
      isPrimary: true,
    },
    {
      title: 'Total Orders',
      value: '1,842',
      trend: '8.5% from last month',
      isUp: true,
      isPrimary: false,
    },
    {
      title: 'Avg. Order Value',
      value: '$77.40',
      trend: 'Stable performance',
      isUp: null,
      isPrimary: false,
    },
    {
      title: 'Total Covers',
      value: '3,844',
      trend: '3.2% from last month',
      isUp: false,
      isPrimary: false,
    },
  ];

const KPIGrid: React.FC<KPIGridProps> = ({ kpis = defaultKpis }) => {
  return (
    <div className="grid grid-cols-4 gap-6 mb-12">
      {kpis.map((kpi, index) => (
        <div key={index} className="bg-white p-8 rounded-xl border border-[#d0c5af]/30 flex flex-col gap-2 col-span-1">
          <h3 className="text-xs font-semibold text-[#4d4635] uppercase tracking-widest">{kpi.title}</h3>
          <p className={`font-serif text-[40px] leading-tight ${kpi.isPrimary ? 'text-[#735c00]' : 'text-[#191c1d]'}`}>
            {kpi.value}
          </p>
          {kpi.isUp !== null ? (
            <p className={`text-xs font-semibold flex items-center gap-1 ${kpi.isUp ? 'text-[#254188]' : 'text-[#ba1a1a]'}`}>
              <span className="material-symbols-outlined text-[16px]">
                {kpi.isUp ? 'trending_up' : 'trending_down'}
              </span>
              <span>{kpi.trend}</span>
            </p>
          ) : (
            <p className="text-xs text-[#4d4635] opacity-60">{kpi.trend}</p>
          )}
        </div>
      ))}
    </div>
  );
};

export default KPIGrid;
