import React, { useEffect, useState, useMemo } from 'react';
import { menuApi, type MenuItemDto, type CategoryDto } from '../api/gateway.api';

const resolveItemImage = (item: MenuItemDto): string | null => {
  const img = item.image_url ?? item.imageUrl;
  if (!img) return null;
  if (img.startsWith('http')) return img;
  if (img.startsWith('assets/')) return img.replace('assets/', '/');
  return img.startsWith('/') ? img : `/${img}`;
};

export default function MenuPage() {
  const [items, setItems] = useState<MenuItemDto[]>([]);
  const [categories, setCategories] = useState<CategoryDto[]>([]);
  const [activeCategoryId, setActiveCategoryId] = useState<string>('all');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    Promise.all([
      menuApi.listCategories(),
      menuApi.listItems({ page_size: 100 }),
    ])
      .then(([catRes, itemRes]) => {
        setCategories(catRes.categories ?? []);
        setItems(itemRes.items ?? []);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  // Backend puts category_id UUID into the "category" field of MenuItem.
  // Build a lookup map to resolve UUID → display name.
  const categoryMap = useMemo(() => {
    const m = new Map<string, string>();
    categories.forEach((cat) => {
      const id = cat.category_id ?? cat.categoryId ?? '';
      if (id) m.set(id, cat.name ?? '');
    });
    return m;
  }, [categories]);

  const getItemCategoryId = (it: MenuItemDto) =>
    it.category_id ?? it.categoryId ?? it.category ?? '';

  const filteredItems = useMemo(() => {
    if (activeCategoryId === 'all') return items;
    return items.filter((it) => getItemCategoryId(it) === activeCategoryId);
  }, [items, activeCategoryId]);

  return (
    <div className="bg-background text-on-surface selection:bg-primary-fixed-dim selection:text-on-primary-fixed min-h-screen flex flex-col font-sans">
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
            A symphony of seasonal flavors, meticulously curated for the discerning palate.
            Our ingredients are sourced daily from local artisanal producers.
          </p>
        </header>

        {/* Category Tabs */}
        <div className="flex flex-wrap justify-center gap-3 mb-20">
          <button
            onClick={() => setActiveCategoryId('all')}
            className={`px-5 py-2.5 rounded-full text-sm font-medium tracking-wide transition-all duration-200 ${
              activeCategoryId === 'all'
                ? 'bg-primary text-on-primary shadow-md scale-105'
                : 'bg-surface-container-high text-on-surface-variant hover:bg-outline-variant/50'
            }`}
          >
            All Selection
          </button>
          {categories.map((cat) => {
            const id = cat.category_id ?? cat.categoryId ?? '';
            return (
              <button
                key={id}
                onClick={() => setActiveCategoryId(id)}
                className={`px-5 py-2.5 rounded-full text-sm font-medium tracking-wide transition-all duration-200 ${
                  activeCategoryId === id
                    ? 'bg-primary text-on-primary shadow-md scale-105'
                    : 'bg-surface-container-high text-on-surface-variant hover:bg-outline-variant/50'
                }`}
              >
                {cat.name}
              </button>
            );
          })}
        </div>

        {/* Menu Items */}
        {loading ? (
          <div className="flex items-center justify-center py-24">
            <span className="material-symbols-outlined animate-spin text-4xl text-primary">
              progress_activity
            </span>
          </div>
        ) : (
          <div className="space-y-24">
            {filteredItems.map((item, index) => {
              const isEven = index % 2 === 0;
              return (
                <div
                  key={item.item_id ?? item.itemId ?? index}
                  className="flex flex-col md:flex-row gap-12 items-center transition-all duration-500 animate-in fade-in slide-in-from-bottom-6"
                >
                  <div
                    className={`w-full md:w-1/2 overflow-hidden rounded-xl shadow-lg transition-transform duration-300 hover:scale-[1.02] ${
                      isEven ? 'order-1 md:order-1' : 'order-1 md:order-2'
                    }`}
                  >
                    {resolveItemImage(item) ? (
                      <img
                        alt={item.name}
                        className="w-full aspect-[4/3] object-cover"
                        src={resolveItemImage(item)!}
                      />
                    ) : (
                      <div className="w-full aspect-[4/3] bg-surface-container-high flex flex-col items-center justify-center gap-3 text-on-surface-variant">
                        <span className="material-symbols-outlined text-5xl opacity-40">restaurant</span>
                        <span className="font-body-md text-sm opacity-60">{item.name}</span>
                      </div>
                    )}
                  </div>

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
                        {item.price?.toLocaleString('vi-VN')} ₫
                      </span>
                    </div>
                    {item.description && (
                      <p className="font-body-lg text-on-surface-variant italic leading-relaxed">
                        {item.description}
                      </p>
                    )}
                    {getItemCategoryId(item) && categoryMap.has(getItemCategoryId(item)) && (
                      <span className="inline-block text-xs font-semibold text-primary/70 uppercase tracking-widest">
                        {categoryMap.get(getItemCategoryId(item))}
                      </span>
                    )}
                  </div>
                </div>
              );
            })}

            {!loading && filteredItems.length === 0 && (
              <div className="text-center py-12 text-on-surface-variant italic">
                No menu items found in this section.
              </div>
            )}
          </div>
        )}

        <div className="mt-32 text-center">
          <p className="font-label-sm text-label-sm uppercase tracking-[0.2em] text-outline">
            An 18% service charge will be added to all tables.
          </p>
        </div>
      </main>
    </div>
  );
}
