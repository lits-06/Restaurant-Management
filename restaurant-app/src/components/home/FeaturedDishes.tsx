const DISHES = [
  {
    id: 1,
    name: "Atlantic Glazed Salmon",
    description: "Wild-caught, infused with ginger-soy and served over a bed of truffled leak purée.",
    price: "$42.00",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuAzUKUVu5-66jLwcjGRFDJODRbS_e5CzDtSpXVqfFQiHDGdTrTGj6D-9EZkyFcqrwsGMb-RFExHcLJgD8CHE-tCIiOr879Y8t63OlUpOSpwOSsVe7_TzKYC0UVF1-ABzAqRVhomlypLfnXwLn4ulqcxstDb8rQqAIFDV78V42hoSroocRonloRReULbslfGcsp1yFGkIwEbBkaVEosFNopQ7lIf6Mv6XYQ77BmBhrEfDVtO5UDDKUHpjD01E6VSvHvx-af5ZQCiuho",
    hasBorder: true
  },
  {
    id: 2,
    name: "Heritage Wagyu A5",
    description: "Slow-roasted beef with rosemary jus and garlic-confit fingerling potatoes.",
    price: "$125.00",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuDaJEMOH2Nyddn77lvIxW2kaZHyUNixvD6_8foRYvRcw9mzFUTvjDMqx9aV0IkhYL6l4jemZ-6ntfi_8b-ivP3FHpHcapGx3MU90qfYSahaT3xaRys2VPqykI2XxO3YyLYGMkJKj9Nugsih45CrNqQVQDZCbeanNFiYItNoA_PORIfTISGPAszrvuhpZrQwF6eWcpSgXzDlQQu1ONuFV9V4GLc8DPN7RhjYAHsXue6-PrwxYxIuwmpmOHra8ikX4wAJ-5wuxVnRhJE"
  },
  {
    id: 3,
    name: "Truffle Agnolotti",
    description: "Handmade pasta filled with ricotta and sage, finished with fresh Umbrian truffles.",
    price: "$38.00",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuCaa-E2LXx_pFsG6RUNGOiH5MssoakQs7eNaV6Ir1TXdstodVrtZ_I86k5vDX8EBkhAkr_SlQQu8NnBtMSGqLwVrSytTx3dsSnLLlpGZgiXSN3bcKA-TGooFho4gri0uJqLppr00_x3ToSui_o3dausazOmuw_iXYx5pNR38r9dLSE_3uLOwUXt4UvaaKaVe7xmGjRwhOXWdsCBz8ST-szoXQkIJl0kgeTqD_IW_rwO4E253DXKRwYnq_m7VTlCWsmnMt9OzM4Ntf4"
  }
];

export default function FeaturedDishes() {
  return (
    <section className="bg-surface-container-low py-24">
      <div className="px-margin-desktop w-full max-w-container-max mx-auto">
        <div className="flex justify-between items-end mb-12">
          <div className="space-y-2">
            <h2 className="text-3xl font-bold">Chef's Signature</h2>
            <p className="text-on-surface-variant">The essence of our seasonal tasting menu.</p>
          </div>
          <button className="text-primary text-xs font-semibold flex items-center gap-2 hover:underline decoration-primary">
            View Full Menu
            <span className="material-symbols-outlined text-sm">open_in_new</span>
          </button>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {DISHES.map((dish) => (
            <div key={dish.id} className="group bg-surface-container-lowest p-4 rounded-xl shadow-sm hover:shadow-md transition-all">
              <div className="aspect-[4/3] overflow-hidden rounded-lg mb-6">
                <img 
                  className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" 
                  alt={dish.name} 
                  src={dish.image}
                />
              </div>
              {dish.hasBorder && <div className="border-t-4 border-primary-container -mt-4 mb-4 w-12"></div>}
              <h3 className="text-xl font-bold mb-2">{dish.name}</h3>
              <p className="text-on-surface-variant text-base mb-4">{dish.description}</p>
              <span className="text-primary font-bold">{dish.price}</span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}