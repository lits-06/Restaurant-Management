import React from 'react';

interface DishData {
  name: string;
  unitsSold: number;
  revenue: string;
}

interface PerformanceTableProps {
  dishes?: DishData[];
}

const defaultDishes: DishData[] = [
    { name: 'Wagyu Ribeye Steak', unitsSold: 420, revenue: '$37,380.00' },
    { name: 'Truffle Cacio e Pepe', unitsSold: 385, revenue: '$12,320.00' },
    { name: 'Pan Seared Scallops', unitsSold: 310, revenue: '$10,540.00' },
    { name: 'Maine Lobster Risotto', unitsSold: 275, revenue: '$13,200.00' },
    { name: 'Crispy Duck Confit', unitsSold: 190, revenue: '$7,410.00' },
  ];

const PerformanceTable: React.FC<PerformanceTableProps> = ({ dishes = defaultDishes }) => {
  return (
    <div className="bg-white p-8 rounded-xl border border-[#d0c5af]/30 col-span-12">
      <div className="flex justify-between items-center mb-8">
        <h3 className="font-serif text-2xl font-semibold text-[#191c1d]">Signature Performance</h3>
        <button className="text-[#735c00] text-xs font-semibold hover:underline">
          Download Detailed CSV
        </button>
      </div>
      <table className="w-full text-left">
        <thead>
          <tr className="border-b border-[#d0c5af]/30 text-xs text-[#4d4635] uppercase tracking-widest">
            <th className="pb-4 font-semibold">Dish Name</th>
            <th className="pb-4 font-semibold text-right">Units Sold</th>
            <th className="pb-4 font-semibold text-right">Revenue</th>
          </tr>
        </thead>
        <tbody className="text-base text-[#191c1d]">
          {dishes.length === 0 ? (
            <tr>
              <td className="py-8 text-center text-sm text-[#4d4635]" colSpan={3}>
                No order item data available.
              </td>
            </tr>
          ) : dishes.map((dish, index) => (
            <tr 
              key={index} 
              className={`border-b border-[#d0c5af]/10 hover:bg-white transition-colors ${
                index === dishes.length - 1 ? 'border-none' : ''
              }`}
            >
              <td className="py-5 font-semibold">{dish.name}</td>
              <td className="py-5 text-right">{dish.unitsSold}</td>
              <td className="py-5 text-right font-semibold">{dish.revenue}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default PerformanceTable;
