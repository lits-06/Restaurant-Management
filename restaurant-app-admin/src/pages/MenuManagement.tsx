import React, { useEffect, useMemo, useState } from 'react';
import { menuApi, type MenuItemDto } from '../services/api';

// ── types ──────────────────────────────────────────────────────────────────────

interface Dish {
  id: string;
  name: string;
  price: number;
  description: string;
  image: string;
  category: string;   // display name
  categoryId: string; // for API calls
  isVip: boolean;
}

interface DishForm {
  name: string;
  price: string;
  description: string;
  image: string;
  categoryId: string;
}

const EMPTY_FORM: DishForm = { name: '', price: '', description: '', image: '', categoryId: '' };

// ── helpers ────────────────────────────────────────────────────────────────────

const getItemId = (item: MenuItemDto) => item.item_id ?? item.itemId ?? '';

const getImageUrl = (item: MenuItemDto) => {
  const image = item.image_url ?? item.imageUrl;
  if (!image) return '/images/default-food.jpg';
  if (image.startsWith('assets/')) return image.replace('assets/', '/');
  return image.startsWith('/') ? image : `/${image}`;
};

const fmtVnd = (n: number) => `${Math.round(n).toLocaleString('vi-VN')}đ`;

// ── component ──────────────────────────────────────────────────────────────────

const MenuManagement: React.FC = () => {
  const [dishes, setDishes]         = useState<Dish[]>([]);
  const [categories, setCategories] = useState<{ id: string; name: string }[]>([]);
  const [activeCategory, setActiveCategory] = useState('All');
  const [searchQuery, setSearchQuery]       = useState('');
  const [isLoading, setIsLoading]           = useState(true);
  const [error, setError]                   = useState('');

  // Modal
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalMode, setModalMode]     = useState<'add' | 'edit'>('add');
  const [editingDish, setEditingDish] = useState<Dish | null>(null);
  const [form, setForm]               = useState<DishForm>(EMPTY_FORM);
  const [isSaving, setIsSaving]       = useState(false);

  // ── helpers that close over categories ────────────────────────────────────────

  const mapDto = (item: MenuItemDto): Dish => {
    const catId   = item.category_id ?? item.categoryId ?? '';
    const catName = categories.find(c => c.id === catId)?.name ?? item.category ?? 'Other';
    return {
      id:          getItemId(item),
      name:        item.name ?? '',
      price:       item.price ?? 0,
      description: item.description ?? '',
      image:       getImageUrl(item),
      category:    catName,
      categoryId:  catId,
      isVip:       (item.price ?? 0) >= 1000000,
    };
  };

  // ── load ───────────────────────────────────────────────────────────────────────

  const loadCategories = async () => {
    try {
      const res = await menuApi.listCategories();
      return (res.categories ?? []).map(cat => ({ id: cat.category_id ?? cat.categoryId ?? '', name: cat.name ?? 'Other' }));
    } catch {
      return [];
    }
  };

  const loadMenu = async (cats: { id: string; name: string }[]) => {
    setIsLoading(true);
    setError('');
    try {
      const res = await menuApi.listItems({ page: 1, page_size: 100 });
      const mapped = (res.items ?? []).map(item => {
        const catId   = item.category_id ?? item.categoryId ?? '';
        const catName = cats.find(c => c.id === catId)?.name ?? item.category ?? 'Other';
        return {
          id:          getItemId(item),
          name:        item.name ?? '',
          price:       item.price ?? 0,
          description: item.description ?? '',
          image:       getImageUrl(item),
          category:    catName,
          categoryId:  catId,
          isVip:       (item.price ?? 0) >= 1000000,
        } as Dish;
      }).filter(d => d.id);
      setDishes(mapped);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load menu.');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadCategories().then(cats => {
      setCategories(cats);
      loadMenu(cats);
    });
  }, []);

  // ── filtering ─────────────────────────────────────────────────────────────────

  const filteredDishes = useMemo(() =>
    dishes.filter(d => {
      const matchCat  = activeCategory === 'All' || d.category === activeCategory;
      const matchName = d.name.toLowerCase().includes(searchQuery.toLowerCase());
      return matchCat && matchName;
    }),
    [dishes, activeCategory, searchQuery]
  );

  // ── modal ─────────────────────────────────────────────────────────────────────

  const openModal = (mode: 'add' | 'edit', dish?: Dish) => {
    setModalMode(mode);
    setEditingDish(dish ?? null);
    setForm(dish
      ? { name: dish.name, price: String(dish.price), description: dish.description, image: dish.image, categoryId: dish.categoryId }
      : EMPTY_FORM
    );
    setError('');
    setIsModalOpen(true);
  };

  const handleDelete = async (dish: Dish) => {
    if (!window.confirm(`Delete "${dish.name}"?`)) return;
    try {
      await menuApi.deleteItem(dish.id);
      setDishes(prev => prev.filter(d => d.id !== dish.id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete item.');
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    setError('');

    const payload = {
      name:        form.name.trim(),
      description: form.description.trim(),
      price:       Number(form.price) || 0,
      category_id: form.categoryId,
      image_url:   form.image.trim(),
    };

    try {
      if (modalMode === 'edit' && editingDish) {
        const res = await menuApi.updateItem(editingDish.id, payload);
        const updated: Dish = res.item
          ? mapDto(res.item)
          : { ...editingDish, ...payload, category: categories.find(c => c.id === form.categoryId)?.name ?? editingDish.category, image: payload.image_url, isVip: payload.price >= 1000000 };
        setDishes(prev => prev.map(d => d.id === editingDish.id ? updated : d));
      } else {
        const res = await menuApi.createItem(payload);
        if (res.item) {
          setDishes(prev => [mapDto(res.item!), ...prev]);
        } else {
          await loadMenu(categories);
        }
      }
      setIsModalOpen(false);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to save item.');
    } finally {
      setIsSaving(false);
    }
  };

  // ── render ────────────────────────────────────────────────────────────────────

  const categoryTabs = ['All', ...categories.map(c => c.name)];

  return (
    <div className="flex flex-col gap-6 animate-fadeIn">
      {/* Header */}
      <div className="flex justify-between items-end border-b border-outline-variant/20 pb-6">
        <div>
          <nav className="flex gap-2 text-on-surface-variant text-[10px] mb-2 uppercase tracking-widest">
            <span>Admin</span><span>/</span><span className="text-primary font-bold">Menu</span>
          </nav>
          <h2 className="font-serif text-5xl font-bold text-on-surface">Menu Management</h2>
        </div>
        <button
          onClick={() => openModal('add')}
          className="bg-[#735c00] text-white text-xs font-semibold px-6 py-3 rounded-lg shadow-sm hover:shadow-md transition-all active:scale-95 flex items-center gap-2"
        >
          <span className="material-symbols-outlined text-[18px]">add</span>
          Add Item
        </button>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap justify-between items-center gap-6 my-2">
        {/* Category tabs */}
        <div className="flex gap-6 border-b border-outline-variant/30 pb-px overflow-x-auto">
          {categoryTabs.map(name => (
            <button
              key={name}
              onClick={() => setActiveCategory(name)}
              className={`text-xs font-semibold py-2 relative whitespace-nowrap transition-colors ${activeCategory === name ? 'text-[#735c00]' : 'text-on-surface-variant hover:text-[#735c00]'}`}
            >
              {name}
              {activeCategory === name && <span className="absolute bottom-0 left-0 w-full h-[2px] bg-[#735c00]" />}
            </button>
          ))}
        </div>
        {/* Search */}
        <div className="relative">
          <span className="material-symbols-outlined absolute left-3 top-1/2 -translate-y-1/2 text-on-surface-variant text-[20px]">search</span>
          <input
            className="pl-10 pr-4 py-2 bg-[#f3f4f5] border border-outline-variant rounded-full text-sm w-64 focus:outline-none focus:border-[#735c00] transition-all"
            placeholder="Search dishes..."
            type="text"
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      {error   && <p className="text-sm text-red-600">{error}</p>}
      {isLoading && <p className="text-sm text-on-surface-variant">Loading menu...</p>}

      {/* Dish grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
        {filteredDishes.map(dish => (
          <div key={dish.id} className="group relative bg-white rounded-xl shadow-sm border border-outline-variant/10 overflow-hidden hover:shadow-lg transition-all duration-300 flex flex-col">
            <div className="relative h-48 overflow-hidden bg-gray-100">
              {dish.isVip && (
                <div className="absolute top-0 left-0 bg-[#735c00] px-3 py-1 text-white text-[10px] z-10 font-bold uppercase tracking-widest rounded-br-lg">
                  VIP Selection
                </div>
              )}
              <img alt={dish.name} className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110" src={dish.image} />
              <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 flex items-center justify-center gap-4 transition-opacity duration-300">
                <button
                  onClick={() => openModal('edit', dish)}
                  className="w-12 h-12 bg-white text-on-surface rounded-full flex items-center justify-center hover:bg-[#735c00] hover:text-white transition-all shadow-lg"
                >
                  <span className="material-symbols-outlined">edit</span>
                </button>
                <button
                  onClick={() => handleDelete(dish)}
                  className="w-12 h-12 bg-white text-red-600 rounded-full flex items-center justify-center hover:bg-red-600 hover:text-white transition-all shadow-lg"
                >
                  <span className="material-symbols-outlined">delete</span>
                </button>
              </div>
            </div>
            <div className={`p-6 flex-grow ${dish.isVip ? 'border-t-4 border-[#735c00]' : ''}`}>
              <div className="flex justify-between items-start mb-1">
                <h3 className="font-serif text-xl font-bold text-on-surface">{dish.name}</h3>
                <span className="text-sm font-bold text-[#735c00] ml-2 whitespace-nowrap">{fmtVnd(dish.price)}</span>
              </div>
              <p className="text-[11px] text-on-surface-variant mb-1">{dish.category}</p>
              <p className="text-on-surface-variant text-sm line-clamp-2">{dish.description}</p>
            </div>
          </div>
        ))}

        {/* Add placeholder card */}
        <div
          onClick={() => openModal('add')}
          className="border-2 border-dashed border-outline-variant/50 rounded-xl flex flex-col items-center justify-center p-12 group cursor-pointer hover:border-[#735c00]/50 hover:bg-[#edeeef] transition-all"
        >
          <div className="w-16 h-16 rounded-full bg-[#e1e3e4] flex items-center justify-center group-hover:bg-[#735c00] group-hover:text-white transition-all mb-4 text-on-surface-variant">
            <span className="material-symbols-outlined text-[32px]">add_circle</span>
          </div>
          <p className="text-xs font-semibold text-on-surface-variant group-hover:text-[#735c00] transition-colors">Add new dish</p>
        </div>
      </div>

      {/* ── MODAL: Add / Edit ──────────────────────────────────────────────────── */}
      {isModalOpen && (
        <div className="fixed inset-0 z-[100] flex items-center justify-center bg-black/60 backdrop-blur-sm">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-2xl overflow-hidden">
            <div className="px-8 py-6 border-b border-outline-variant/20 flex justify-between items-center">
              <h3 className="font-serif text-xl font-bold text-on-surface">
                {modalMode === 'add' ? 'Add New Item' : 'Update Item'}
              </h3>
              <button className="text-on-surface-variant hover:text-[#735c00] transition-colors" onClick={() => setIsModalOpen(false)}>
                <span className="material-symbols-outlined">close</span>
              </button>
            </div>

            <form className="p-8 space-y-5" onSubmit={handleSubmit}>
              {error && <p className="text-sm text-red-600">{error}</p>}

              <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Dish Name</label>
                  <input
                    className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface"
                    placeholder="e.g. Wagyu Beef A5..."
                    type="text"
                    value={form.name}
                    onChange={e => setForm({ ...form, name: e.target.value })}
                    required
                  />
                </div>
                <div className="space-y-1">
                  <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Price (VNĐ)</label>
                  <input
                    className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface"
                    placeholder="0"
                    type="number"
                    min={0}
                    value={form.price}
                    onChange={e => setForm({ ...form, price: e.target.value })}
                    required
                  />
                </div>
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Category</label>
                {categories.length > 0 ? (
                  <select
                    className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface"
                    value={form.categoryId}
                    onChange={e => setForm({ ...form, categoryId: e.target.value })}
                    required
                  >
                    <option value="">-- Select category --</option>
                    {categories.map(cat => (
                      <option key={cat.id} value={cat.id}>{cat.name}</option>
                    ))}
                  </select>
                ) : (
                  <input
                    className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface"
                    placeholder="e.g. Main Course"
                    type="text"
                    value={form.categoryId}
                    onChange={e => setForm({ ...form, categoryId: e.target.value })}
                  />
                )}
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Description</label>
                <textarea
                  className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface resize-none"
                  placeholder="Ingredients, preparation..."
                  rows={3}
                  value={form.description}
                  onChange={e => setForm({ ...form, description: e.target.value })}
                />
              </div>

              <div className="space-y-1">
                <label className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider block">Image URL</label>
                <input
                  className="w-full px-4 py-3 bg-[#f3f4f5] border border-outline-variant rounded-lg focus:ring-1 focus:ring-[#735c00] focus:border-[#735c00] outline-none text-on-surface"
                  placeholder="/images/mon-an.jpg"
                  type="text"
                  value={form.image}
                  onChange={e => setForm({ ...form, image: e.target.value })}
                />
              </div>

              <div className="flex justify-end items-center gap-4 pt-2">
                <button type="button" className="text-xs font-semibold text-on-surface-variant hover:text-[#735c00] transition-colors px-4 py-2" onClick={() => setIsModalOpen(false)}>
                  Cancel
                </button>
                <button
                  type="submit"
                  disabled={isSaving}
                  className="px-10 py-3 bg-[#735c00] text-white text-xs font-bold rounded-lg shadow-md hover:bg-[#735c00]/90 transition-all active:scale-95 uppercase tracking-widest disabled:opacity-60"
                >
                  {isSaving ? 'Saving...' : modalMode === 'add' ? 'Create Item' : 'Save Changes'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default MenuManagement;
