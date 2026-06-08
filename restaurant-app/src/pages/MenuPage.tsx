import React, { useState } from 'react';

// 1. Khởi tạo mảng dữ liệu các món ăn và đồ uống từ HTML gốc
const MENU_DATA = [
  {
    id: 1,
    name: "Heirloom Tomato & Burrata",
    price: 24,
    category: "starters",
    description: "Creamy Puglia burrata, balsamic reduction, basil pearls, and hand-selected heirloom tomatoes finished with cold-pressed olive oil.",
    image: "https://lh3.googleusercontent.com/aida/AP1WRLtLp3LoicM9_KKBJhYybJVpwyiTTvC0seuEE3IhWuOmBePPA927R6RcxxQ4YFJ9rh5V-i2ohbmwEPMF5f5JQ-9GQzx7lzpoINgwCdciqVFRkkJVcaiy9fpbtS25gcCa_zvZbzwz_xYLrvs3Gbu_pvtGTu-CeURAGLsVyIE5aMINuvK8jv_xfE1EjPfNuIrNyPAxifaV2TIbvvNjXcnZBo6ywDzdC9bm-QaOuXtM5vinxjnQFkrMg181B1M"
  },
  {
    id: 2,
    name: "Black Truffle Carpaccio",
    price: 28,
    category: "starters",
    description: "Prime Wagyu beef, shaved Perigord truffles, 24-month aged parmesan, and wild arugula.",
    image: "https://lh3.googleusercontent.com/aida/AP1WRLvXxQ9aRNBW__ZBQti_ebm3nDdh5nCxS6605mvtCUQ1lPD9cEp1G9B-I5iHMr3wB7ptYG0FrXPXBOnhLIruYStD_BYvORib9L3wPrLAte6rrdstSxIWw8GlV1UGPVlhCKXd3sF5V_pv0pHsDFZU8ATXiiJ2fEvglVcDkVM3-8co1IgP5PCLu6RRl9TDQhKUZapFGAeZmn_8nqx0-4Mc_4DhC77YfYDC85vLX8mfub7uT3UMiNGtFI-1Jq0"
  },
  {
    id: 3,
    name: "Seared Hokkaido Scallops",
    price: 32,
    category: "starters",
    description: "Pan-seared scallops, cauliflower silk, crispy pancetta, and citrus-infused foam.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuCpY1JgsXFLfcq4XfHMI4oNg9lCjYTSMphO2pWCW3kBVXJmv_DPz2nERfrmkMXkEBr3C_vC_8ZT0YwgKJdWy8-V0ck-jUIyiFNiun0-QXio7_DQAMR-hdTGrbsiAnIwIwSaXafVHvwAHWWrbIxcsqxlbB0ras8ZUxR1oRfvd5xGSBBgcNHHg82s-WZs1KccZfCR429vGaPW7UWAR4mnvZh583kJd2ffDMY29ZywOgSDWNqzABbOHGu0ZpY-tguY8OlPb-10kF7Vgxc"
  },
  {
    id: 4,
    name: "Herb-Crusted Rack of Lamb",
    price: 54,
    category: "mains",
    description: "Provencal herb crust, mint pea purée, glazed heritage carrots, and a vintage Port reduction.",
    image: "https://lh3.googleusercontent.com/aida/AP1WRLvZUphISYjfUN6erOIloFHHAk8z-QMVzO_tSY9SnV7BifaGMMOHyNd8pIiLOomzTbkv8ZTNOAr29UQ51oDlH2SUXLcR80xNbXXMMW8TjpWFLekxsA5Lu7WFjpx6iNFXxT-tVlKKLEcWwFHbR3T8O1A1uAnmfmX7eipHc8T2ZsMmx3MYWetmfre-iW-L7TfxjZRPJnW5Fulb-_MwXyOxELnGOsjZl85RIglfH0ugZxJod2jEEHvCju5kH7I"
  },
  {
    id: 5,
    name: "Miso Glazed Black Cod",
    price: 48,
    category: "mains",
    description: "Sustainably sourced cod, ginger-infused dashi, baby bok choy, and toasted sesame oil drizzle.",
    image: "https://lh3.googleusercontent.com/aida/AP1WRLuspDtxoqCr8LXxVumV2y1X5zq0t-NCvC86Z6PlPcZCxMv-ZqOtlJZ0yW2ZR4wxneoHm6xvNQZzJZs9aXAInYeqCqQ8n-0gm12bFTPt0mej9s1QjJY3GEF_GRP0ab2_GzYDDx7XeiuX7iwToEDTZi0YM8h-BwqwRLt-uL3QRsoqxKLFJ3uIlPTRHbaC2xyyFZpfDJUk8u3zMyASwWfHhlocub8GxW2EKoqalM2xS5u6TAD7WzMGmSNMBMU"
  },
  {
    id: 6,
    name: "Gold Leaf Chocolate Fondant",
    price: 18,
    category: "desserts",
    description: "70% Dark Valrhona chocolate, 24k edible gold, vanilla bean gelato.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuAeeQ5vO1tBSjK9lbwE21GUeiei19a10_4WJ8aQeQ28SgQPJuyDcr7lrOBwxQascCmeurqoxtiE1CI9LpAIz8ER-Ukpk5pHb6ErLmFGNngCwqV2XpGWQXGnqj1U_18ONQw2MerSiFcWBLLiXbIXwpfG-bDuyHyoKkdFBF3068IZe3nbMKpyDt4UpAPdr5Hhok-WCp_LkR-Qf6U0DwQzvseimUW0rdgmWsuA-iIWO7enZDrN3Dm6msqFYzc4gLu68Cuyc1Et69t4uo8"
  },
  {
    id: 7,
    name: "Deconstructed Lemon Tart",
    price: 16,
    category: "desserts",
    description: "Sicilian lemon curd, torched meringue, shortbread soil, raspberry coulis.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuCYBVbJ6kvRAzeouOa1GfwbImqngt3J-WFnIQ2nPwQ-whvc6WlfjT65d9o6ot-2ivwkRZkORxj4lW_f-RxferTJenvlJar9EcWtbXryhJB6xE8g5FMZpFXxt2ln44BrYV14jK_S4iTRd4sB6Qqn07Y3DUqoc3W2Te0VBHe-ny7C9iyv8HDzqxi3tDyVLrB2XAUJstld91H5P8WDAdWrgXqECpjhWaACYKuAmhsG8se5fN2NHW8dG9eY6AfgZk4qXHUBV32cFbJNpos"
  },
  {
    id: 8,
    name: "Artisan Cheese Selection",
    price: 22,
    category: "desserts",
    description: "Chef's selection of three European cheeses, honeycomb, sourdough crisps.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuBbu2bLbcLyAkpKOmYm0KvB-EqnqxIOOhfj5tmU4gvDIGp0p1snTbue2rAunJo5ZhKPfw0DPWF6zua-5wYMvsy_h__CZuS92gwEApQOdO5KUJs5YxV2sK83IbV4qv0bKYlv6OQtatqEvW_eNviSVJ4OoUZlktxtFJK5Hrgfc4FaIhfWMVN-dTRgC07x2gw6RQhxF2YJhO-T1REPGMSfdPNc8Tpx0-X8o_VOVz4F0cpyB60LaXRotRvFCDbj00fSZEGqWdzR5UNMHCc"
  },
  {
    id: 9,
    name: "Chateau Margaux, Grand Cru Classe",
    price: 750,
    category: "wines",
    description: "Bordeaux, France | 2015. A vintage of exceptional depth and refined tannins.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuDcry2vwndhOcLvpkoEeZoW8DmFgozCaGb6mzw-ZaR5rHPBfyhSzIGvfusxUQpr30Le4rYH9tOhXrJjOD-GFxOLoPbVKX3k0LWd_bV5WkrvsenwZaBwwPNJPSiJLv-Kw1vuYT_SON9oUegb_YGZ8cNxEL6XojtK-GtYBlH88n-tfDu_cui1iJdNox9iEg62TxL1umV8Y4AvVRu2SEZkUvRQQ_P5-Kno9APmLuPHcoEux13P2vMRkL9xU6FWuoE9-RebRkxOzHV4c1c"
  },
  {
    id: 10,
    name: "Dom Perignon, Brut Champagne",
    price: 420,
    category: "wines",
    description: "Epernay, France | 2012. Luminous and airy, with notes of stone fruit and toast.",
    image: "https://lh3.googleusercontent.com/aida-public/AB6AXuB0sL2FtG5W9ePeQA0b_640332591eqAi7FeeNC9tY7BhHjP-37vH_WrOa9BbTOyKm-743B7GttsOI8aXCp0RJP46NmXUAp-G5P6zWGsouzeeia_xtUqf7uabLczAOMvrtQGUUpnGyU0rtcWqDpqM16fuoWXAtzxC5IvoqpPAD90krHQ92ZioIh4HZtY3RwDTJ_JXeZ7AmSSEnexA0osOPSPT_DrYenJa18sIOXeTxXNM0QwYcLF8AFjmdpRBpWCXppmDvQq7FLwwI"
  }
];

const CATEGORIES = [
  { id: 'all', label: 'All Selection' },
  { id: 'starters', label: 'Starters' },
  { id: 'mains', label: 'Mains' },
  { id: 'desserts', label: 'Desserts' },
  { id: 'wines', label: 'Fine Wines' }
];

export default function MenuPage() {
  const [activeCategory, setActiveCategory] = useState('all');

  // Lọc món ăn dựa trên danh mục đang được chọn
  const filteredItems = activeCategory === 'all' 
    ? MENU_DATA 
    : MENU_DATA.filter(item => item.category === activeCategory);

  return (
    <div className="bg-background text-on-surface selection:bg-primary-fixed-dim selection:text-on-primary-fixed min-h-screen flex flex-col font-sans">

      {/* Main Content */}
      <main className="max-w-container-max mx-auto px-margin-mobile md:px-margin-desktop py-20 flex-grow w-full">
        
        {/* Hero Header */}
        <header className="text-center mb-16 space-y-4">
          <span className="font-label-sm text-label-sm uppercase tracking-widest text-primary block">
            Informed Hospitality
          </span>
          <h1 className="font-headline-xl text-headline-xl text-on-surface">
            The Culinary Selection
          </h1>
          <p className="max-w-2xl mx-auto font-body-lg text-body-lg text-on-surface-variant">
            A symphony of seasonal flavors, meticulously curated for the discerning palate. Our ingredients are sourced daily from local artisanal producers.
          </p>
        </header>

        {/* TÍNH NĂNG THÊM: Bộ lọc Tabs chuyển đổi linh hoạt giữa các danh mục */}
        <div className="flex flex-wrap justify-center gap-3 mb-20">
          {CATEGORIES.map((category) => (
            <button
              key={category.id}
              onClick={() => setActiveCategory(category.id)}
              className={`px-5 py-2.5 rounded-full text-sm font-medium tracking-wide transition-all duration-200 ${
                activeCategory === category.id
                  ? 'bg-primary text-on-primary shadow-md scale-105'
                  : 'bg-surface-container-high text-on-surface-variant hover:bg-outline-variant/50'
              }`}
            >
              {category.label}
            </button>
          ))}
        </div>

        {/* Luxurious Menu List với hiệu ứng xen kẽ ảnh tự động */}
        <div className="space-y-24">
          {filteredItems.map((item, index) => {
            // Logic xen kẽ: Phần tử chẵn ảnh bên trái, phần tử lẻ ảnh bên phải trên màn hình Desktop
            const isEven = index % 2 === 0;

            return (
              <div 
                key={item.id} 
                className="menu-item flex flex-col md:flex-row gap-12 items-center transition-all duration-500 animate-in fade-in slide-in-from-bottom-6"
              >
                {/* Image Wrapper */}
                <div 
                  className={`w-full md:w-1/2 overflow-hidden rounded-xl shadow-lg transition-transform duration-300 hover:scale-[1.02] ${
                    isEven ? 'order-1 md:order-1' : 'order-1 md:order-2'
                  }`}
                >
                  <img 
                    alt={item.name} 
                    className="w-full aspect-[4/3] object-cover" 
                    src={item.image} 
                  />
                </div>

                {/* Content Wrapper */}
                <div 
                  className={`w-full md:w-1/2 space-y-4 ${
                    isEven ? 'order-2 md:order-2' : 'order-2 md:order-1'
                  }`}
                >
                  <div className="flex justify-between items-baseline border-b border-outline-variant pb-2">
                    <h3 className="font-headline-md text-headline-md text-on-surface">
                      {item.name}
                    </h3>
                    <span className="text-primary font-headline-md ml-4">
                      ${item.price}
                    </span>
                  </div>
                  <p className="font-body-lg text-on-surface-variant italic leading-relaxed">
                    {item.description}
                  </p>
                </div>
              </div>
            );
          })}

          {/* Trạng thái trống khi không tìm thấy món ăn phù hợp */}
          {filteredItems.length === 0 && (
            <div className="text-center py-12 text-on-surface-variant italic">
              No menu items found in this section.
            </div>
          )}
        </div>

        {/* Service Charge Footer Notes */}
        <div className="mt-32 text-center">
          <p className="font-label-sm text-label-sm uppercase tracking-[0.2em] text-outline">
            An 18% service charge will be added to all tables.
          </p>
        </div>
      </main>
    </div>
  );
}